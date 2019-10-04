/*
 *
 * Copyright 2015 gRPC authors.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 *
 */

// Package main implements a client for Greeter service.
package main

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/coreos/go-systemd/sdjournal"
	grpc_retry "github.com/grpc-ecosystem/go-grpc-middleware/retry"
	pb "github.com/metal-pod/droptailer/dropsink"
	"google.golang.org/grpc"
)

const (
	address            = "localhost:50051"
	nftablesDropPrefix = "nftables-metal-dropped: "
)

type dropreader struct {
	jr *sdjournal.JournalReader
	dc pb.DropSinkClient
}

func main() {
	// Set up a connection to the server.
	opts := []grpc_retry.CallOption{
		grpc_retry.WithBackoff(grpc_retry.BackoffLinear(100 * time.Millisecond)),
	}
	conn, err := grpc.Dial(address, grpc.WithInsecure(),
		// grpc.WithStreamInterceptor(grpc_retry.StreamClientInterceptor(opts...)),
		grpc.WithUnaryInterceptor(grpc_retry.UnaryClientInterceptor(opts...)),
	)

	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	c := pb.NewDropSinkClient(conn)

	jr, err := sdjournal.NewJournalReader(
		sdjournal.JournalReaderConfig{
			NumFromTail: 100,
			// Matches on message only match the whole message not the start
			Matches: []sdjournal.Match{
				{
					Field: sdjournal.SD_JOURNAL_FIELD_SYSLOG_IDENTIFIER,
					Value: "kernel",
				},
			},
			Formatter: protoMessageFormatter,
		})
	if err != nil {
		log.Fatalf("Error opening journal: %s", err)
	}
	if jr == nil {
		log.Fatalf("Got a nil reader")
	}
	defer jr.Close()
	j := &dropreader{
		jr: jr,
		dc: c,
	}
	j.run()
}

func newCtx(timeout time.Duration) context.Context {
	ctx, _ := context.WithTimeout(context.TODO(), timeout)
	return ctx
}

func (d *dropreader) run() {
	pr, pw := io.Pipe()
	until := make(chan time.Time)
	go func() {
		if err := d.jr.Follow(until, pw); err != sdjournal.ErrExpired {
			log.Fatalf("Error during follow: %s", err)
		}
		pw.Close()
	}()
	d.writeTo(pr)
}

// writeTo
// actual message will be like
// nftables-metal-dropped: IN=vrf104009 OUT= MAC=12:99:fd:3b:ce:f8:1a:ae:e9:a7:95:50:08:00 SRC=222.73.197.30 DST=212.34.89.87 LEN=40 TOS=0x00 PREC=0x00 TTL=238 ID=46474 PROTO=TCP SPT=59265 DPT=445 WINDOW=1024 RES=0x00 SYN URGP=0
func (d *dropreader) writeTo(r io.ReadCloser) {
	br := bufio.NewReader(r)
	for {
		line, _, err := br.ReadLine()
		if err != nil {
			r.Close()
			break
		}

		parts := strings.Split(string(line), "@")
		if len(parts) < 2 {
			continue
		}
		ts, err := strconv.ParseInt(parts[0], 10, 64)
		if err != nil {
			log.Printf("unable to parse timestamp:%v", err)
			continue
		}
		msg := parts[1]
		if !strings.HasPrefix(msg, nftablesDropPrefix) {
			fmt.Printf("message:%s dropped\n", msg)
			continue
		}
		msg = strings.TrimPrefix(msg, nftablesDropPrefix)
		fields := parseFields(msg)
		de := &pb.DropEntry{
			Timestamp: ts,
			Fields:    fields,
		}
		_, err = d.dc.Push(
			newCtx(35*time.Second),
			de,
			grpc_retry.WithMax(30),
			grpc_retry.WithPerRetryTimeout(1*time.Second))
		if err != nil {
			log.Printf("unable to send dropentry:%v", err)
		}
	}
}

func parseFields(msg string) map[string]string {
	fields := make(map[string]string)
	parts := strings.Fields(msg)
	for _, part := range parts {
		fieldParts := strings.Split(part, "=")
		if len(fieldParts) == 0 {
			continue
		}
		key := fieldParts[0]
		fields[key] = ""
		if len(fieldParts) == 1 {
			continue
		}
		fields[key] = fieldParts[1]
	}
	return fields
}

func protoMessageFormatter(entry *sdjournal.JournalEntry) (string, error) {
	msg, ok := entry.Fields["MESSAGE"]
	if !ok {
		return "", fmt.Errorf("no MESSAGE field present in journal entry")
	}
	usec := entry.RealtimeTimestamp
	timestamp := time.Unix(0, int64(usec)*int64(time.Microsecond))
	return fmt.Sprintf("%d@%s\n", timestamp.Unix(), msg), nil
}

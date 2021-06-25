package forwarder

import (
	"fmt"
	"github.com/coreos/go-systemd/v22/sdjournal"
	"time"
)

func messageFormatter(entry *sdjournal.JournalEntry) (string, error) {
	msg, ok := entry.Fields[sdjournal.SD_JOURNAL_FIELD_MESSAGE]
	if !ok {
		return "", fmt.Errorf("no %s field present in journal entry", sdjournal.SD_JOURNAL_FIELD_MESSAGE)
	}
	usec := entry.RealtimeTimestamp
	timestamp := time.Unix(0, int64(usec)*int64(time.Microsecond))
	return fmt.Sprintf("%d@%s\n", timestamp.Unix(), msg), nil
}

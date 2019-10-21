package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"time"

	"github.com/metal-pod/droptailer/pkg/client"
	"github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	k8s "k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

const (
	moduleName      = "droptailer-client"
	serverAppName   = "droptailer"
	certificateBase = "/etc"
)

var (
	debug                  = false
	defaultPrefixesOfDrops = []string{"nftables-metal-dropped: ", "nftables-firewall-dropped: "}
	rootCmd                = &cobra.Command{
		Use:     moduleName,
		Short:   "a service that streams information of dropped packages to a sink in a k8s cluster",
		Version: "",
		Run: func(cmd *cobra.Command, args []string) {
			run()
		},
	}
)

func main() {
	if err := rootCmd.Execute(); err != nil {
		log.Fatalf("failed executing root command; err: %v", err)
	}
}

func init() {
	viper.SetEnvPrefix("droptailer")
	homedir, err := homedir.Dir()
	if err != nil {
		log.Fatal(err)
	}
	rootCmd.PersistentFlags().StringP("kubeconfig", "k", homedir+"/.kube/config", "kubeconfig path to the cluster")
	rootCmd.PersistentFlags().StringSlice("prefixesOfDrops", defaultPrefixesOfDrops, "comma seperated list of kernel journal messages that contain packet drop information")
	rootCmd.PersistentFlags().StringP("serverAddress", "s", "lookup", "directly specify the server address")
	viper.AutomaticEnv()
	err = viper.BindPFlags(rootCmd.PersistentFlags())
	if err != nil {
		log.Fatal(err)
	}
}

func run() {
	cs, err := loadClient(viper.GetString("kubeconfig"))
	if err != nil {
		log.Fatalf("unable to connect to k8s; err: %v", err)
	}

	address := viper.GetString("serverAddress")
	if address == "lookup" {
		address, err = fetchServerAddress(cs)
		if err != nil {
			log.Fatalf("could not fetch server address; err: %v", err)
		}
	}
	_, err = net.DialTimeout("tcp", address, time.Second*5)
	if err != nil {
		log.Fatalf("could not reach droptailer server within 5 seconds: %v", err)
	}

	certs, err := fetchCerts(cs)
	if err != nil {
		log.Fatalf("could not fetch certificates from k8s secrets; err: %v", err)
	}
	c := client.Client{
		ServerAddress:   address,
		PrefixesOfDrops: viper.GetStringSlice("prefixesOfDrops"),
		Certificates:    *certs,
	}
	err = c.Start()
	if err != nil {
		log.Fatalf("client could not start or died, %v", err)
	}
}

func fetchCerts(c k8s.Interface) (*client.Certificates, error) {
	s, err := c.CoreV1().Secrets("droptailer").Get("droptailer-client", metav1.GetOptions{})
	if err != nil {
		return nil, fmt.Errorf("could not find droptailer-client secret; err: %w", err)
	}
	keys := []string{"tls.key", "tls.crt", "ca.crt"}
	for _, k := range keys {
		v := s.Data[k]
		f := fmt.Sprintf("%s/%s", certificateBase, k)
		err = ioutil.WriteFile(f, v, 0640)
		if err != nil {
			return nil, fmt.Errorf("could not write secret to droptailer folder; err: %w", err)
		}
	}
	return &client.Certificates{
		Crt: fmt.Sprintf("%s/%s", certificateBase, "tls.crt"),
		Key: fmt.Sprintf("%s/%s", certificateBase, "tls.key"),
		CA:  fmt.Sprintf("%s/%s", certificateBase, "ca.crt"),
	}, nil
}

func fetchServerAddress(c k8s.Interface) (string, error) {
	labelMap := map[string]string{"app": serverAppName}
	opts := metav1.ListOptions{
		LabelSelector: labels.SelectorFromSet(labelMap).String(),
	}
	for {
		watcher, err := c.CoreV1().Pods(serverAppName).Watch(opts)
		if err != nil {
			return "", fmt.Errorf("could not watch for pods; err: %w", err)
		}
		for event := range watcher.ResultChan() {
			p, ok := event.Object.(*v1.Pod)
			if !ok || p.Status.PodIP == "" {
				continue
			}
			return p.Status.PodIP, nil
		}
	}
}

func loadClient(kubeconfigPath string) (*k8s.Clientset, error) {
	config, err := clientcmd.BuildConfigFromFlags("", kubeconfigPath)
	if err != nil {
		return nil, err
	}
	return k8s.NewForConfig(config)
}

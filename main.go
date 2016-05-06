package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/mlafeldt/havoc/chaosmonkey"
)

func main() {
	var (
		groupName string
		chaosType string
		endpoint  string
	)

	flag.StringVar(&groupName, "group-name", "", "Group name")
	flag.StringVar(&chaosType, "chaos-type", "ShutdownInstance", "Chaos type")
	flag.StringVar(&endpoint, "endpoint", "http://127.0.0.1:8080", "Endpoint")
	flag.Parse()

	if groupName == "" {
		abort("group name missing")
	}

	config := chaosmonkey.Config{
		Endpoint: endpoint,
	}
	client, _ := chaosmonkey.NewClient(&config)
	err := client.TriggerChaosEvent(groupName, chaosType)
	if err != nil {
		abort("%s", err)
	}
}

func abort(format string, a ...interface{}) {
	fmt.Fprintf(os.Stderr, "error: "+format+"\n", a...)
	os.Exit(1)
}

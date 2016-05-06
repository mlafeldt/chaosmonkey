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
		username  string
		password  string
	)

	flag.StringVar(&groupName, "group-name", "", "Group name")
	flag.StringVar(&chaosType, "chaos-type", "ShutdownInstance", "Chaos type")
	flag.StringVar(&endpoint, "endpoint", "http://127.0.0.1:8080", "HTTP endpoint")
	flag.StringVar(&username, "username", "", "HTTP username")
	flag.StringVar(&password, "password", "", "HTTP password")
	flag.Parse()

	config := chaosmonkey.Config{
		Endpoint: endpoint,
		Username: username,
		Password: password,
	}

	client, err := chaosmonkey.NewClient(&config)
	if err != nil {
		abort("%s", err)
	}

	if groupName != "" {
		event, err := client.TriggerEvent(groupName, chaosType)
		if err != nil {
			abort("%s", err)
		}
		fmt.Printf("%+v\n", event)
	} else {
		events, err := client.GetEvents()
		if err != nil {
			abort("%s", err)
		}
		for _, e := range events {
			fmt.Printf("%+v\n", e)
		}
	}
}

func abort(format string, a ...interface{}) {
	fmt.Fprintf(os.Stderr, "error: "+format+"\n", a...)
	os.Exit(1)
}

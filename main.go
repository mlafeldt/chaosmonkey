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

	if groupName == "" {
		abort("group name missing")
	}

	config := chaosmonkey.Config{
		Endpoint: endpoint,
		Username: username,
		Password: password,
	}

	client, err := chaosmonkey.NewClient(&config)
	if err != nil {
		abort("%s", err)
	}

	if err := client.TriggerEvent(groupName, chaosType); err != nil {
		abort("%s", err)
	}
}

func abort(format string, a ...interface{}) {
	fmt.Fprintf(os.Stderr, "error: "+format+"\n", a...)
	os.Exit(1)
}

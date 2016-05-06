package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/mlafeldt/havoc/chaosmonkey"
)

func main() {
	var (
		asgName  string
		strategy string
		endpoint string
		username string
		password string
	)

	flag.StringVar(&asgName, "asg", "", "Name of auto scaling group")
	flag.StringVar(&strategy, "strategy", "", "Chaos strategy to use")
	flag.StringVar(&endpoint, "endpoint", "http://127.0.0.1:8080", "HTTP endpoint")
	flag.StringVar(&username, "username", "", "HTTP username")
	flag.StringVar(&password, "password", "", "HTTP password")
	flag.Parse()

	config := chaosmonkey.Config{
		Endpoint:   endpoint,
		Username:   username,
		Password:   password,
		HTTPClient: &http.Client{Timeout: 10 * time.Second},
	}

	client, err := chaosmonkey.NewClient(&config)
	if err != nil {
		abort("%s", err)
	}

	if asgName != "" {
		event, err := client.TriggerEvent(asgName, strategy)
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

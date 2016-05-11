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
		group    string
		strategy string
		endpoint string
		username string
		password string
	)

	flag.StringVar(&group, "group", "", "Name of auto scaling group")
	flag.StringVar(&strategy, "strategy", "", "Chaos strategy to use")
	flag.StringVar(&endpoint, "endpoint", "", "HTTP endpoint")
	flag.StringVar(&username, "username", "", "HTTP username")
	flag.StringVar(&password, "password", "", "HTTP password")
	flag.Parse()

	client, err := chaosmonkey.NewClient(&chaosmonkey.Config{
		Endpoint:   endpoint,
		Username:   username,
		Password:   password,
		HTTPClient: &http.Client{Timeout: 10 * time.Second},
	})
	if err != nil {
		abort("%s", err)
	}

	if group != "" {
		event, err := client.TriggerEvent(group, strategy)
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

package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/mlafeldt/havoc/chaosmonkey"
	"github.com/ryanuber/columnize"
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
		printEvents(*event)
	} else {
		events, err := client.GetEvents()
		if err != nil {
			abort("%s", err)
		}
		printEvents(events...)
	}
}

func printEvents(event ...chaosmonkey.ChaosEvent) {
	lines := []string{"InstanceID|AutoScalingGroupName|Region|Strategy|TriggeredAt"}
	for _, e := range event {
		lines = append(lines, fmt.Sprintf("%s|%s|%s|%s|%s",
			e.InstanceID,
			e.AutoScalingGroupName,
			e.Region,
			e.Strategy,
			e.TriggeredAt.Format(time.RFC3339),
		))
	}
	fmt.Println(columnize.SimpleFormat(lines))
}

func abort(format string, a ...interface{}) {
	fmt.Fprintf(os.Stderr, "error: "+format+"\n", a...)
	os.Exit(1)
}

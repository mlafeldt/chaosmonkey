package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"
	"runtime"
	"time"

	"github.com/ryanuber/columnize"

	"github.com/mlafeldt/chaosmonkey/aws"
	chaosmonkey "github.com/mlafeldt/chaosmonkey/lib"
)

func main() {
	var (
		endpoint = flag.String("endpoint", "", "Address and port of Chaos Monkey API server")
		username = flag.String("username", "", "Username for HTTP basic authentication")
		password = flag.String("password", "", "Password for HTTP basic authentication")

		group    = flag.String("group", "", "Name of auto scaling group, see -list-groups")
		strategy = flag.String("strategy", "", "Chaos strategy to use, see -list-strategies")

		count    = flag.Int("count", 1, "Number of times to trigger chaos event")
		interval = flag.Duration("interval", 5*time.Second, "Time to wait between chaos events")

		listStrategies = flag.Bool("list-strategies", false, "List chaos strategies")
		listGroups     = flag.Bool("list-groups", false, "List auto scaling groups")
		wipeState      = flag.String("wipe-state", "", "Wipe state of Chaos Monkey by deleting given SimpleDB domain")
		showVersion    = flag.Bool("version", false, "Show program version")
	)
	flag.Parse()

	if flag.NArg() > 0 {
		abort("program expects no arguments, but %d given", flag.NArg())
	}

	if v := os.Getenv("CHAOSMONKEY_ENDPOINT"); v != "" && *endpoint == "" {
		*endpoint = v
	}
	if v := os.Getenv("CHAOSMONKEY_USERNAME"); v != "" && *username == "" {
		*username = v
	}
	if v := os.Getenv("CHAOSMONKEY_PASSWORD"); v != "" && *password == "" {
		*password = v
	}

	switch {
	case *listStrategies:
		for _, s := range chaosmonkey.Strategies {
			fmt.Println(s)
		}
		return
	case *listGroups:
		if err := listAutoScalingGroups(); err != nil {
			abort("failed to list auto scaling groups: %s", err)
		}
		return
	case *wipeState != "":
		if err := aws.DeleteSimpleDBDomain(*wipeState); err != nil {
			abort("failed to wipe state: %s", err)
		}
		return
	case *showVersion:
		fmt.Printf("chaosmonkey %s %s/%s %s\n", Version,
			runtime.GOOS, runtime.GOARCH, runtime.Version())
		return
	}

	client, err := chaosmonkey.NewClient(&chaosmonkey.Config{
		Endpoint:   *endpoint,
		Username:   *username,
		Password:   *password,
		UserAgent:  fmt.Sprintf("chaosmonkey Go client %s", Version),
		HTTPClient: &http.Client{Timeout: 10 * time.Second},
	})
	if err != nil {
		abort("%s", err)
	}

	if *group != "" {
		for i := 1; i <= *count; i++ {
			event, err := client.TriggerEvent(*group, chaosmonkey.Strategy(*strategy))
			if err != nil {
				abort("%s", err)
			}
			printEvents(*event)
			if i < *count {
				time.Sleep(*interval)
			}
		}
	} else {
		events, err := client.Events()
		if err != nil {
			abort("%s", err)
		}
		printEvents(events...)
	}
}

func listAutoScalingGroups() error {
	groups, err := aws.AutoScalingGroups()
	if err != nil {
		return err
	}
	lines := []string{"AutoScalingGroupName|Instances|Desired|Min|Max"}
	for _, g := range groups {
		lines = append(lines, fmt.Sprintf("%s|%d|%d|%d|%d",
			g.Name,
			g.InstancesInService,
			g.DesiredCapacity,
			g.MinSize,
			g.MaxSize,
		))
	}
	fmt.Println(columnize.SimpleFormat(lines))
	return nil
}

var addHeader = true

func printEvents(event ...chaosmonkey.Event) {
	var lines []string
	if addHeader {
		lines = append(lines, "InstanceID|AutoScalingGroupName|Region|Strategy|TriggeredAt")
		addHeader = false
	}
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

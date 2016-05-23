package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"
	"sort"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/autoscaling"

	"github.com/mlafeldt/chaosmonkey"
	"github.com/ryanuber/columnize"
)

func main() {
	var (
		group    string
		strategy string
		endpoint string
		username string
		password string

		listGroups     bool
		listStrategies bool
	)

	flag.StringVar(&group, "group", "", "Name of auto scaling group")
	flag.StringVar(&strategy, "strategy", "", "Chaos strategy to use")
	flag.StringVar(&endpoint, "endpoint", "", "HTTP endpoint")
	flag.StringVar(&username, "username", "", "HTTP username")
	flag.StringVar(&password, "password", "", "HTTP password")
	flag.BoolVar(&listGroups, "list-groups", false, "List auto scaling groups")
	flag.BoolVar(&listStrategies, "list-strategies", false, "List default chaos strategies")
	flag.Parse()

	if listGroups {
		groups, err := autoScalingGroups()
		if err != nil {
			abort("failed to get auto scaling groups: %s", err)
		}
		fmt.Println(strings.Join(groups, "\n"))
		return
	}

	if listStrategies {
		for _, s := range chaosmonkey.Strategies {
			fmt.Println(s)
		}
		return
	}

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
		event, err := client.TriggerEvent(group, chaosmonkey.Strategy(strategy))
		if err != nil {
			abort("%s", err)
		}
		printEvents(*event)
	} else {
		events, err := client.Events()
		if err != nil {
			abort("%s", err)
		}
		printEvents(events...)
	}
}

func autoScalingGroups() ([]string, error) {
	var groups []string
	svc := autoscaling.New(session.New())
	err := svc.DescribeAutoScalingGroupsPages(nil, func(out *autoscaling.DescribeAutoScalingGroupsOutput, last bool) bool {
		for _, g := range out.AutoScalingGroups {
			groups = append(groups, aws.StringValue(g.AutoScalingGroupName))
		}
		return !last
	})
	if err != nil {
		return nil, err
	}
	sort.Strings(groups)
	return groups, nil
}

func printEvents(event ...chaosmonkey.Event) {
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

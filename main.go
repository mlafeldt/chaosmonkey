package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"
	"runtime"
	"sort"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/autoscaling"
	"github.com/aws/aws-sdk-go/service/simpledb"
	"github.com/ryanuber/columnize"

	chaosmonkey "github.com/mlafeldt/chaosmonkey/lib"
)

// Version is the current version of the chaosmonkey tool. A ".dev" suffix
// denotes that the version is currently being developed.
const Version = "v0.3.0.dev"

func main() {
	var (
		group    string
		strategy string
		endpoint string
		username string
		password string

		listGroups     bool
		listStrategies bool
		wipeState      string
		showVersion    bool
	)

	flag.StringVar(&group, "group", "", "Name of auto scaling group")
	flag.StringVar(&strategy, "strategy", "", "Chaos strategy to use")
	flag.StringVar(&endpoint, "endpoint", "", "HTTP endpoint")
	flag.StringVar(&username, "username", "", "HTTP username")
	flag.StringVar(&password, "password", "", "HTTP password")
	flag.BoolVar(&listGroups, "list-groups", false, "List auto scaling groups")
	flag.BoolVar(&listStrategies, "list-strategies", false, "List default chaos strategies")
	flag.StringVar(&wipeState, "wipe-state", "", "Wipe Chaos Monkey state by deleting given SimpleDB domain")
	flag.BoolVar(&showVersion, "version", false, "Show program version")
	flag.Parse()

	switch {
	case listGroups:
		groups, err := autoScalingGroups()
		if err != nil {
			abort("failed to get auto scaling groups: %s", err)
		}
		fmt.Println(strings.Join(groups, "\n"))
		return
	case listStrategies:
		for _, s := range chaosmonkey.Strategies {
			fmt.Println(s)
		}
		return
	case wipeState != "":
		if err := deleteSimpleDBDomain(wipeState); err != nil {
			abort("failed to wipe state: %s", err)
		}
		return
	case showVersion:
		fmt.Printf("chaosmonkey %s %s/%s %s\n", Version,
			runtime.GOOS, runtime.GOARCH, runtime.Version())
		return
	}

	client, err := chaosmonkey.NewClient(&chaosmonkey.Config{
		Endpoint:   endpoint,
		Username:   username,
		Password:   password,
		UserAgent:  fmt.Sprintf("chaosmonkey Go client %s", Version),
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

func deleteSimpleDBDomain(domainName string) error {
	var domainExists bool
	svc := simpledb.New(session.New())
	err := svc.ListDomainsPages(nil, func(out *simpledb.ListDomainsOutput, last bool) bool {
		for _, n := range out.DomainNames {
			if aws.StringValue(n) == domainName {
				domainExists = true
			}
		}
		return !last
	})
	if !domainExists {
		return fmt.Errorf("SimpleDB domain %q does not exist", domainName)
	}
	_, err = svc.DeleteDomain(&simpledb.DeleteDomainInput{
		DomainName: aws.String(domainName),
	})
	return err
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

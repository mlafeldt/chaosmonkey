package main

import (
	"fmt"
	"net/http"
	"os"
	"runtime"
	"strings"
	"time"

	"github.com/ryanuber/columnize"
	flag "github.com/spf13/pflag"
	"github.com/spf13/viper"

	chaosmonkey "github.com/mlafeldt/chaosmonkey/lib"
)

func main() {
	var (
		group          = flag.String("group", "", "Name of auto scaling group")
		strategy       = flag.String("strategy", "", "Chaos strategy to use")
		_              = flag.String("endpoint", "", "HTTP endpoint")
		_              = flag.String("username", "", "HTTP username")
		_              = flag.String("password", "", "HTTP password")
		listGroups     = flag.Bool("list-groups", false, "List auto scaling groups")
		listStrategies = flag.Bool("list-strategies", false, "List default chaos strategies")
		wipeState      = flag.String("wipe-state", "", "Wipe Chaos Monkey state by deleting given SimpleDB domain")
		showVersion    = flag.Bool("version", false, "Show program version")
	)
	flag.Parse()

	if flag.NArg() > 0 {
		abort("program expects no arguments, but %d given", flag.NArg())
	}

	viper.AutomaticEnv()
	viper.SetEnvPrefix("chaosmonkey")
	viper.BindPFlag("endpoint", flag.Lookup("endpoint"))
	viper.BindPFlag("username", flag.Lookup("username"))
	viper.BindPFlag("password", flag.Lookup("password"))

	switch {
	case *listGroups:
		groups, err := autoScalingGroups()
		if err != nil {
			abort("failed to get auto scaling groups: %s", err)
		}
		fmt.Println(strings.Join(groups, "\n"))
		return
	case *listStrategies:
		for _, s := range chaosmonkey.Strategies {
			fmt.Println(s)
		}
		return
	case *wipeState != "":
		if err := deleteSimpleDBDomain(*wipeState); err != nil {
			abort("failed to wipe state: %s", err)
		}
		return
	case *showVersion:
		fmt.Printf("chaosmonkey %s %s/%s %s\n", Version,
			runtime.GOOS, runtime.GOARCH, runtime.Version())
		return
	}

	client, err := chaosmonkey.NewClient(&chaosmonkey.Config{
		Endpoint:   viper.GetString("endpoint"),
		Username:   viper.GetString("username"),
		Password:   viper.GetString("password"),
		UserAgent:  fmt.Sprintf("chaosmonkey Go client %s", Version),
		HTTPClient: &http.Client{Timeout: 10 * time.Second},
	})
	if err != nil {
		abort("%s", err)
	}

	if *group != "" {
		event, err := client.TriggerEvent(*group, chaosmonkey.Strategy(*strategy))
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

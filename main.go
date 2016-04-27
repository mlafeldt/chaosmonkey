package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"net/http"
	"os"
)

type Event struct {
	MonkeyType string `json:"monkeyType,omitempty"`
	EventType  string `json:"eventType"`
	// EventTime  time.Date `json:"eventTime,omitempty"`
	Region    string `json:"region,omitempty"`
	GroupType string `json:"groupType"`
	GroupName string `json:"groupName"`
	ChaosType string `json:"chaosType,omitempty"`
}

type Result struct {
	*Event
	Message string `json:"message"`
}

func main() {
	var (
		region    string
		groupName string
		chaosType string
		endpoint  string
	)

	flag.StringVar(&region, "region", "eu-west-1", "AWS region")
	flag.StringVar(&groupName, "group-name", "", "Group name")
	flag.StringVar(&chaosType, "chaos-type", "ShutdownInstance", "Chaos type")
	flag.StringVar(&endpoint, "endpoint", "http://127.0.0.1:8080", "Endpoint")
	flag.Parse()

	if groupName == "" {
		abort("group name missing")
	}

	e := Event{
		MonkeyType: "CHAOS",
		EventType:  "CHAOS_TERMINATION",
		Region:     region,
		GroupType:  "ASG",
		GroupName:  groupName,
		ChaosType:  chaosType,
	}

	data, err := json.Marshal(&e)
	if err != nil {
		panic(err)
	}

	url := endpoint + "/simianarmy/api/v1/chaos"
	resp, err := http.Post(url, "application/json", bytes.NewReader(data))
	if err != nil {
		abort("%s", err)
	}
	defer resp.Body.Close()

	fmt.Println(resp.StatusCode)

	var res Result
	if err := json.NewDecoder(resp.Body).Decode(&res); err != nil {
		abort("%s", err)
	}

	fmt.Printf("%+v\n", res)

	if res.Message != "" {
		abort(res.Message)
	}
}

func abort(format string, a ...interface{}) {
	fmt.Fprintf(os.Stderr, "error: "+format+"\n", a...)
	os.Exit(1)
}

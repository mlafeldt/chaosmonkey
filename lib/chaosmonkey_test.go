package chaosmonkey_test

import (
	"flag"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"reflect"
	"testing"
	"time"

	chaosmonkey "github.com/mlafeldt/chaosmonkey/lib"
)

const newEvent = `
  {
    "monkeyType": "CHAOS",
    "eventId": "i-12345678",
    "eventType": "CHAOS_TERMINATION",
    "eventTime": 1460116927834,
    "region": "eu-west-1",
    "groupType": "ASG",
    "groupName": "SomeAutoScalingGroup",
    "chaosType": "ShutdownInstance"
  }
`

const pastEvents = `[
  {
    "monkeyType": "CHAOS",
    "eventId": "i-12345678",
    "eventType": "CHAOS_TERMINATION",
    "eventTime": 1460116927834,
    "region": "eu-west-1",
    "groupType": "ASG",
    "groupName": "SomeAutoScalingGroup",
    "chaosType": "ShutdownInstance"
  },
  {
    "monkeyType": "CHAOS",
    "eventId": "i-87654321",
    "eventType": "CHAOS_TERMINATION",
    "eventTime": 1460116816326,
    "region": "us-east-1",
    "groupType": "ASG",
    "groupName": "AnotherAutoScalingGroup",
    "chaosType": "BlockAllNetworkTraffic"
  }
]`

var client *chaosmonkey.Client

func TestMain(m *testing.M) {
	flag.Parse()

	// Test server faking Chaos Monkey REST API
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/simianarmy/api/v1/chaos" {
			switch r.Method {
			case "POST":
				fmt.Fprint(w, newEvent)
				return
			case "GET":
				fmt.Fprint(w, pastEvents)
				return
			}
		}
		http.NotFound(w, r)
	}))
	defer ts.Close()

	var err error
	client, err = chaosmonkey.NewClient(&chaosmonkey.Config{Endpoint: ts.URL})
	if err != nil {
		panic(err)
	}

	os.Exit(m.Run())
}

func TestTriggerEvent(t *testing.T) {
	event, err := client.TriggerEvent("SomeAutoScalingGroup", chaosmonkey.StrategyShutdownInstance)
	if err != nil {
		t.Fatal(err)
	}

	expected := &chaosmonkey.Event{
		InstanceID:           "i-12345678",
		AutoScalingGroupName: "SomeAutoScalingGroup",
		Region:               "eu-west-1",
		Strategy:             chaosmonkey.StrategyShutdownInstance,
		TriggeredAt:          time.Unix(1460116927, 0).UTC(),
	}

	if !reflect.DeepEqual(event, expected) {
		t.Fatalf("\ngot:  %+v\nwant: %+v\n", event, expected)
	}
}

func TestEvents(t *testing.T) {
	events, err := client.Events()
	if err != nil {
		t.Fatal(err)
	}

	expected := []chaosmonkey.Event{
		{
			InstanceID:           "i-12345678",
			AutoScalingGroupName: "SomeAutoScalingGroup",
			Region:               "eu-west-1",
			Strategy:             chaosmonkey.StrategyShutdownInstance,
			TriggeredAt:          time.Unix(1460116927, 0).UTC(),
		},
		{
			InstanceID:           "i-87654321",
			AutoScalingGroupName: "AnotherAutoScalingGroup",
			Region:               "us-east-1",
			Strategy:             chaosmonkey.StrategyBlockAllNetworkTraffic,
			TriggeredAt:          time.Unix(1460116816, 0).UTC(),
		},
	}

	if !reflect.DeepEqual(events, expected) {
		t.Fatalf("\ngot:  %+v\nwant: %+v\n", events, expected)
	}
}

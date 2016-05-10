// Package chaosmonkey allows to talk to the Chaos Monkey REST API to trigger
// and retrieve chaos events.
package chaosmonkey

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

const (
	// StrategyShutdownInstance ...
	StrategyShutdownInstance = "ShutdownInstance"

	// StrategyBlockAllNetworkTraffic ...
	StrategyBlockAllNetworkTraffic = "BlockAllNetworkTraffic"

	// StrategyDetachVolumes ...
	StrategyDetachVolumes = "DetachVolumes"

	// StrategyBurnCPU ...
	StrategyBurnCPU = "BurnCpu"

	// StrategyBurnIO ...
	StrategyBurnIO = "BurnIo"

	// StrategyKillProcesses ...
	StrategyKillProcesses = "KillProcesses"

	// StrategyNullRoute ...
	StrategyNullRoute = "NullRoute"

	// StrategyFailEC2 ...
	StrategyFailEC2 = "FailEc2"

	// StrategyFailDNS ...
	StrategyFailDNS = "FailDns"

	// StrategyFailDynamoDB ...
	StrategyFailDynamoDB = "FailDynamoDb"

	// StrategyFailS3 ...
	StrategyFailS3 = "FailS3"

	// StrategyFillDisk ...
	StrategyFillDisk = "FillDisk"

	// StrategyNetworkCorruption ...
	StrategyNetworkCorruption = "NetworkCorruption"

	// StrategyNetworkLatency ...
	StrategyNetworkLatency = "NetworkLatency"

	// StrategyNetworkLoss ...
	StrategyNetworkLoss = "NetworkLoss"
)

// ChaosEvent describes when and how Chaos Monkey terminated an EC2 instance.
type ChaosEvent struct {
	// Name of chaos strategy that was used, e.g. "ShutdownInstance"
	Strategy string

	// Name of auto scaling group containing the terminated EC2 instance
	ASGName string

	// ID of EC2 instance that was terminated
	InstanceID string

	// AWS region of EC2 instance and its auto scaling group
	Region string

	// Time of the chaos event
	Time time.Time
}

type apiRequest struct {
	EventType string `json:"eventType"`
	GroupType string `json:"groupType"`
	GroupName string `json:"groupName"`
	ChaosType string `json:"chaosType,omitempty"`
}

type apiResponse struct {
	*apiRequest

	MonkeyType string `json:"monkeyType"`
	EventID    string `json:"eventId"`
	EventTime  int64  `json:"eventTime"`
	Region     string `json:"region"`
}

// Config is used to configure the creation of the client.
type Config struct {
	// Address of the Chaos Monkey API server
	Endpoint string

	// Optional username for HTTP Basic Authentication
	Username string

	// Optional password for HTTP Basic Authentication
	Password string

	// Custom HTTP client to use (http.DefaultClient by default)
	HTTPClient *http.Client
}

// Client is the client to the Chaos Monkey API. Create a client with NewClient.
type Client struct {
	config *Config
}

// NewClient returns a new client for the given configuration.
func NewClient(c *Config) (*Client, error) {
	if c.Endpoint == "" {
		return nil, fmt.Errorf("Endpoint must not be empty")
	}
	if c.HTTPClient == nil {
		c.HTTPClient = http.DefaultClient
	}
	return &Client{config: c}, nil
}

// TriggerEvent triggers a new chaos event.
func (c *Client) TriggerEvent(asgName, strategy string) (*ChaosEvent, error) {
	url := c.config.Endpoint + "/simianarmy/api/v1/chaos"

	body, err := json.Marshal(apiRequest{
		EventType: "CHAOS_TERMINATION",
		GroupType: "ASG",
		GroupName: asgName,
		ChaosType: strategy,
	})
	if err != nil {
		return nil, err
	}

	var resp apiResponse
	if err := c.sendRequest("POST", url, bytes.NewReader(body), &resp); err != nil {
		return nil, err
	}

	return makeChaosEvent(&resp), nil
}

// GetEvents returns a list of past chaos events.
func (c *Client) GetEvents() ([]ChaosEvent, error) {
	url := c.config.Endpoint + "/simianarmy/api/v1/chaos"

	var resp []apiResponse
	if err := c.sendRequest("GET", url, nil, &resp); err != nil {
		return nil, err
	}

	var events []ChaosEvent
	for _, r := range resp {
		events = append(events, *makeChaosEvent(&r))
	}

	return events, nil
}

func (c *Client) sendRequest(method, url string, body io.Reader, out interface{}) error {
	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return err
	}

	if c.config.Username != "" && c.config.Password != "" {
		req.SetBasicAuth(c.config.Username, c.config.Password)
	}

	resp, err := c.config.HTTPClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return decodeError(resp)
	}

	return json.NewDecoder(resp.Body).Decode(out)
}

func decodeError(resp *http.Response) error {
	var m struct {
		Message string `json:"message"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&m); err == nil && m.Message != "" {
		return fmt.Errorf("%s", m.Message)
	}
	return fmt.Errorf("%s", resp.Status)
}

func makeChaosEvent(in *apiResponse) *ChaosEvent {
	return &ChaosEvent{
		Strategy:   in.ChaosType,
		ASGName:    in.GroupName,
		InstanceID: in.EventID,
		Region:     in.Region,
		Time:       time.Unix(in.EventTime/1000, 0).UTC(),
	}
}

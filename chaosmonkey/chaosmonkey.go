package chaosmonkey

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

type chaosRequest struct {
	EventType string `json:"eventType"`
	GroupType string `json:"groupType"`
	GroupName string `json:"groupName"`
	ChaosType string `json:"chaosType,omitempty"`
}

type chaosResponse struct {
	*chaosRequest

	MonkeyType string `json:"monkeyType"`
	EventID    string `json:"eventId"`
	EventTime  int64  `json:"eventTime"`
	Region     string `json:"region"`
}

type ChaosEvent struct {
	Strategy   string
	ASGName    string
	InstanceID string
	Region     string
	Time       time.Time
}

type Config struct {
	Endpoint string
	Username string
	Password string

	HTTPClient *http.Client
}

type Client struct {
	config *Config
}

func (c *Config) ReadEnvironment() error {
	return nil
}

func NewClient(c *Config) (*Client, error) {
	if c.HTTPClient == nil {
		c.HTTPClient = http.DefaultClient
	}
	return &Client{config: c}, nil
}

func (c *Client) TriggerEvent(asgName, strategy string) (*ChaosEvent, error) {
	body, err := json.Marshal(chaosRequest{
		EventType: "CHAOS_TERMINATION",
		GroupType: "ASG",
		GroupName: asgName,
		ChaosType: strategy,
	})
	if err != nil {
		return nil, err
	}

	url := c.config.Endpoint + "/simianarmy/api/v1/chaos"

	var event chaosResponse
	if err := c.sendRequest("POST", url, bytes.NewReader(body), &event); err != nil {
		return nil, err
	}

	return &ChaosEvent{
		Strategy:   event.ChaosType,
		ASGName:    event.GroupName,
		InstanceID: event.EventID,
		Region:     event.Region,
		Time:       time.Unix(event.EventTime/1000, 0),
	}, nil
}

func (c *Client) GetEvents() ([]ChaosEvent, error) {
	url := c.config.Endpoint + "/simianarmy/api/v1/chaos"

	var events []chaosResponse
	if err := c.sendRequest("GET", url, nil, &events); err != nil {
		return nil, err
	}

	var chaosEvents []ChaosEvent
	for _, event := range events {
		e := ChaosEvent{
			Strategy:   event.ChaosType,
			ASGName:    event.GroupName,
			InstanceID: event.EventID,
			Region:     event.Region,
			Time:       time.Unix(event.EventTime/1000, 0),
		}
		chaosEvents = append(chaosEvents, e)
	}

	return chaosEvents, nil
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

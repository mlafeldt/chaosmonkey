package chaosmonkey

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
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

type Config struct {
	Username string
	Password string
	Endpoint string

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
		c.HTTPClient = &http.Client{
			Timeout: 15 * time.Second,
		}
	}
	return &Client{config: c}, nil
}

func (c *Client) TriggerChaosEvent(groupName, chaosType, region string) error {
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
		return err
	}

	url := c.config.Endpoint + "/simianarmy/api/v1/chaos"
	resp, err := c.SendRequest("POST", url, bytes.NewReader(data))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		fmt.Println(resp.StatusCode)
	}

	var res Result
	if err := json.NewDecoder(resp.Body).Decode(&res); err != nil {
		return err
	}

	fmt.Printf("%+v\n", res)

	if res.Message != "" {
		return fmt.Errorf("%s", res.Message)
	}

	return nil
}

func (c *Client) NewRequest(method, url string, body io.Reader) (*http.Request, error) {
	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return nil, err
	}
	if c.config.Username != "" && c.config.Password != "" {
		req.SetBasicAuth(c.config.Username, c.config.Password)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Add("User-Agent", "havoc")
	return req, nil
}

func (c *Client) SendRequest(method, url string, body io.Reader) (*http.Response, error) {
	req, err := c.NewRequest(method, url, body)
	if err != nil {
		return nil, err
	}
	return c.config.HTTPClient.Do(req)
}

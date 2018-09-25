package airship

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

const (
	// TreatmentOff indicates that an entity is not part of a treatment.
	TreatmentOff = "off"
	// TreatmentOn indicates that an entity is part of a treatment.
	TreatmentOn = "on"

	defaultTimeout = 10 * time.Second
)

// Client is an object that has all the data to "configure" the Airship Go SDK
type Client struct {
	url    string
	client *http.Client
}

// FeatureFlag is an object that represents a flag in the SDK.
type FeatureFlag struct {
	Name   string
	Client *Client
}

type requestDataWrapper struct {
	Flag   string      `json:"flag"`
	Entity interface{} `json:"entity"`
}

type objectValuesContainer struct {
	Treatment  string          `json:"treatment"`
	Payload    json.RawMessage `json:"payload"`
	IsEligible bool            `json:"isEligible"`
	IsEnabled  bool            `json:"isEnabled"`
}

// ClientOption provides optional configuration for a new *Client.
type ClientOption func(*clientParams)

type clientParams struct {
	client *http.Client
}

// WithHTTPClient sets the *http.Client used by the returned Client.
func WithHTTPClient(c *http.Client) ClientOption {
	return func(p *clientParams) {
		p.client = c
	}
}

var defaultClient = &Client{}

// Configure sets up a Airship Go SDK singleton for the airship package.
// Once Configure is called, one may call methods directly on the package.
// E.g., airship.Flag("flag-name")
func Configure(c *Client) {
	defaultClient = c
}

func New(envKey, edgeURL string, opts ...ClientOption) *Client {
	params := &clientParams{}
	for _, opt := range opts {
		opt(params)
	}

	if params.client == nil {
		params.client = &http.Client{
			Timeout: defaultTimeout,
		}
	}

	return &Client{
		url:    edgeURL + "/v2/object-values/" + envKey,
		client: params.client,
	}
}

// Flag returns an FeatureFlag object that represents the flag.
func Flag(flagName string) *FeatureFlag {
	return defaultClient.Flag(flagName)
}

// Flag (a method on an instance of the SDK) returns an FeatureFlag object that represents the flag.
func (c *Client) Flag(flagName string) *FeatureFlag {
	return &FeatureFlag{
		Name:   flagName,
		Client: c,
	}
}

// GetTreatment returns the treatment value or codename for the flag for a particular entity.
func (f *FeatureFlag) GetTreatment(entity interface{}) (string, error) {
	treatment, err := getTreatment(f, f.Client, entity)
	if err != nil {
		return TreatmentOff, fmt.Errorf("airship: %v", err)
	}
	return treatment, nil
}

func getTreatment(flag *FeatureFlag, client *Client, entity interface{}) (string, error) {
	objectValues, err := getObjectValues(flag, client, entity)
	if err != nil {
		return TreatmentOff, err
	}
	return objectValues.Treatment, nil
}

// GetPayload unmarshals the JSON payload value associated with the flag for a particular entity.
// Pass a pointer as the second argument just as you would to json.Unmarshal.
func (f *FeatureFlag) GetPayload(entity interface{}, v interface{}) error {
	err := getPayload(f, f.Client, entity, v)
	if err != nil {
		return fmt.Errorf("airship: %v", err)
	}
	return nil
}

func getPayload(flag *FeatureFlag, client *Client, entity interface{}, v interface{}) error {
	objectValues, err := getObjectValues(flag, client, entity)
	if err != nil {
		return err
	}
	return json.Unmarshal(objectValues.Payload, v)
}

// IsEligible returns whether or not an entity is part of a population (sampled or yet to be sampled) associated with the flag.
func (f *FeatureFlag) IsEligible(entity interface{}) (bool, error) {
	eligible, err := isEligible(f, f.Client, entity)
	if err != nil {
		return false, fmt.Errorf("airship: %v", err)
	}
	return eligible, nil
}

func isEligible(flag *FeatureFlag, client *Client, entity interface{}) (bool, error) {
	objectValues, err := getObjectValues(flag, client, entity)
	if err != nil {
		return false, err
	}
	return objectValues.IsEligible, nil
}

// IsEnabled returns whether or not an entity is sampled inside a population and given a non-off treatment.
func (f *FeatureFlag) IsEnabled(entity interface{}) (bool, error) {
	enabled, err := isEnabled(f, f.Client, entity)
	if err != nil {
		return false, fmt.Errorf("airship: %v", err)
	}
	return enabled, nil
}

func isEnabled(flag *FeatureFlag, client *Client, entity interface{}) (bool, error) {
	objectValues, err := getObjectValues(flag, client, entity)
	if err != nil {
		return false, err
	}
	return objectValues.IsEnabled, nil
}

func getObjectValues(flag *FeatureFlag, client *Client, entity interface{}) (*objectValuesContainer, error) {
	requestObj, err := json.Marshal(&requestDataWrapper{
		Flag:   flag.Name,
		Entity: entity,
	})
	if err != nil {
		return nil, err
	}

	resp, err := client.client.Post(client.url, "application/json", bytes.NewBuffer(requestObj))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected response status: %s", resp.Status)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var objectValues objectValuesContainer
	if err := json.Unmarshal(body, &objectValues); err != nil {
		return nil, err
	}

	return &objectValues, nil
}

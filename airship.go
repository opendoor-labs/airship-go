package airship

import (
	"bytes"
	"time"

	"encoding/json"
	"io/ioutil"
	"net/http"
)

var defaultTimeout = 10 * time.Second

// Client is an object that has all the data to "configure" the Airship Go SDK
type Client struct {
	EnvKey         string
	EdgeURL        string
	RequestTimeout time.Duration
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

var defaultClient = &Client{}

// Configure sets up a Airship Go SDK singleton for the airship package.
// Once Configure is called, one may call methods directly on the package.
// E.g., airship.Flag("flag-name")
func Configure(c *Client) {
	defaultClient = c
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
func (f *FeatureFlag) GetTreatment(entity interface{}) string {
	return getTreatment(f, f.Client, entity)
}

func getTreatment(flag *FeatureFlag, client *Client, entity interface{}) string {
	objectValues, err := getObjectValues(flag, client, entity)
	if err != nil {
		return ""
	}
	return objectValues.Treatment
}

// GetPayload unmarshals the JSON payload value associated with the flag for a particular entity.
// Pass a pointer as the second argument just as you would to json.Unmarshal.
func (f *FeatureFlag) GetPayload(entity interface{}, v interface{}) error {
	return getPayload(f, f.Client, entity, v)
}

func getPayload(flag *FeatureFlag, client *Client, entity interface{}, v interface{}) error {
	objectValues, err := getObjectValues(flag, client, entity)
	if err != nil {
		return err
	}
	return json.Unmarshal(objectValues.Payload, v)
}

// IsEligible returns whether or not an entity is part of a population (sampled or yet to be sampled) associated with the flag.
func (f *FeatureFlag) IsEligible(entity interface{}) bool {
	return isEligible(f, f.Client, entity)
}

func isEligible(flag *FeatureFlag, client *Client, entity interface{}) bool {
	objectValues, err := getObjectValues(flag, client, entity)
	if err != nil {
		return false
	}
	return objectValues.IsEligible
}

// IsEnabled returns whether or not an entity is sampled inside a population and given a non-off treatment.
func (f *FeatureFlag) IsEnabled(entity interface{}) bool {
	return isEnabled(f, f.Client, entity)
}

func isEnabled(flag *FeatureFlag, client *Client, entity interface{}) bool {
	objectValues, err := getObjectValues(flag, client, entity)
	if err != nil {
		return false
	}
	return objectValues.IsEnabled
}

func getObjectValues(flag *FeatureFlag, client *Client, entity interface{}) (*objectValuesContainer, error) {
	requstObj, _ := json.Marshal(&requestDataWrapper{
		Flag:   flag.Name,
		Entity: entity,
	})
	requestTimeout := client.RequestTimeout
	if requestTimeout == 0 {
		requestTimeout = defaultTimeout
	}
	var netClient = &http.Client{
		Timeout: client.RequestTimeout,
	}
	res, err := netClient.Post(client.EdgeURL+"/v2/object-values/"+client.EnvKey, "application/json", bytes.NewBuffer(requstObj))
	if err != nil {
		return nil, err
	}
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	var objectValues objectValuesContainer
	json.Unmarshal(body, &objectValues)
	return &objectValues, nil
}

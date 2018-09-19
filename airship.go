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

// AirshipFlag is an object that represents a flag in the SDK.
type AirshipFlag struct {
	Name   string
	Client *Client
}

type ObjectValuesBody struct {
	Flag   string      `json:"flag"`
	Entity interface{} `json:"entity"`
}

type ObjectValues struct {
	Treatment  string          `json:"treatment"`
	Payload    json.RawMessage `json:"payload"`
	IsEligible bool            `json:"isEligible"`
	IsEnabled  bool            `json:"isEnabled"`
}

var DefaultClient = &Client{}

// Configure sets up a Airship Go SDK singleton for the airship package.
// Once Configure is called, one may call methods directly on the package.
// E.g., airship.Flag("flag-name")
func Configure(c *Client) {
	DefaultClient = c
}

// Flag returns an AirshipFlag object that represents the flag.
func Flag(flagName string) *AirshipFlag {
	return DefaultClient.Flag(flagName)
}

// Flag (a method on an instance of the SDK) returns an AirshipFlag object that represents the flag.
func (c *Client) Flag(flagName string) *AirshipFlag {
	return &AirshipFlag{
		Name:   flagName,
		Client: c,
	}
}

// GetTreatment returns the treatment value or codename for the flag for a particular entity.
func (f *AirshipFlag) GetTreatment(entity interface{}) string {
	return getTreatment(f, f.Client, entity)
}

func getTreatment(flag *AirshipFlag, client *Client, entity interface{}) string {
	objectValues, err := getObjectValues(flag, client, entity)
	if err != nil {
		return ""
	}
	return objectValues.Treatment
}

// GetPayload unmarshals the JSON payload value associated with the flag for a particular entity.
// Pass a pointer as the second argument just as you would to json.Unmarshal.
func (f *AirshipFlag) GetPayload(entity interface{}, v interface{}) error {
	return getPayload(f, f.Client, entity, v)
}

func getPayload(flag *AirshipFlag, client *Client, entity interface{}, v interface{}) error {
	objectValues, err := getObjectValues(flag, client, entity)
	if err != nil {
		return err
	}
	return json.Unmarshal(objectValues.Payload, v)
}

// IsEligible returns whether or not an entity is part of a population (sampled or yet to be sampled) associated with the flag.
func (f *AirshipFlag) IsEligible(entity interface{}) bool {
	return isEligible(f, f.Client, entity)
}

func isEligible(flag *AirshipFlag, client *Client, entity interface{}) bool {
	objectValues, err := getObjectValues(flag, client, entity)
	if err != nil {
		return false
	}
	return objectValues.IsEligible
}

// IsEnabled returns whether or not an entity is sampled inside a population and given a non-off treatment.
func (f *AirshipFlag) IsEnabled(entity interface{}) bool {
	return isEnabled(f, f.Client, entity)
}

func isEnabled(flag *AirshipFlag, client *Client, entity interface{}) bool {
	objectValues, err := getObjectValues(flag, client, entity)
	if err != nil {
		return false
	}
	return objectValues.IsEnabled
}

func getObjectValues(flag *AirshipFlag, client *Client, entity interface{}) (*ObjectValues, error) {
	objJson, _ := json.Marshal(&ObjectValuesBody{
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
	res, err := netClient.Post("http://"+client.EdgeURL+"/v2/object-values/"+client.EnvKey, "application/json", bytes.NewBuffer(objJson))
	if err != nil {
		return nil, err
	}
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	var objectValues ObjectValues
	json.Unmarshal(body, &objectValues)
	return &objectValues, nil
}

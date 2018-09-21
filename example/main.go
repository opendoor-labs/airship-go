package main

import (
	"fmt"
	airship "github.com/airshiphq/airship-go" // we want this folder to just be called airship
	"os"
)

// This is our example app's expected entity
type Entity struct {
	ID string `json:"id"`
}

// This is our example app's expected payload
type Payload struct {
	Foo string `json:"foo"`
}

func main() {
	newInstanceTest()
	singletonTest()
}

func newInstanceTest() {
	envKey := os.Getenv("ENV_KEY")
	edgeURL := "http://localhost:5000"
	airshipInstance := airship.New(envKey, edgeURL)

	airshipInstanceBitcoinPay := airshipInstance.Flag("bitcoin-pay")
	myEntity := &Entity{
		ID: "1",
	}
	fmt.Println(airshipInstanceBitcoinPay.IsEligible(myEntity))
	fmt.Println(airshipInstanceBitcoinPay.IsEnabled(myEntity))

	fmt.Println(airshipInstanceBitcoinPay.GetTreatment(myEntity))

	var myPayload Payload
	err := airshipInstanceBitcoinPay.GetPayload(myEntity, &myPayload)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(myPayload.Foo)
}

func singletonTest() {
	envKey := os.Getenv("ENV_KEY")
	edgeURL := "http://localhost:5000"
	airship.Configure(airship.New(envKey, edgeURL))

	airshipBitcoinPay := airship.Flag("bitcoin-pay")
	myEntity := &Entity{
		ID: "2",
	}
	fmt.Println(airshipBitcoinPay.IsEnabled(myEntity))
}

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
	airshipInstance := &airship.Client{
		EnvKey:  os.Getenv("ENV_KEY"),
		EdgeURL: "localhost:5000",
	}
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
	airship.Configure(&airship.Client{
		EnvKey:  os.Getenv("ENV_KEY"),
		EdgeURL: "localhost:5000",
	})
	airshipBitcoinPay := airship.Flag("bitcoin-pay")
	myEntity := &Entity{
		ID: "2",
	}
	fmt.Println(airshipBitcoinPay.IsEnabled(myEntity))
}

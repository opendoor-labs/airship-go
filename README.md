# Airship Go SDK

## Prerequisite

This SDK works with the [Airship Microservice](https://github.com/airshiphq/airship-microservice). Please refer to the its documenation before proceeding.

## Installation

```
go get https://github.com/airshiphq/airship-go
```

## Usage
```
import (
    "fmt"
    airship "github.com/username/library"
)

type Entity struct {
    ID string `json:"id"`
}

airship.Configure(&airship.Client{
    EnvKey:  "envKey",
    EdgeURL: "http://localhost:5000",
})

airshipBitcoinPay := airship.Flag("bitcoin-pay")

myEntity := &Entity{
    ID: "2",
}

fmt.Println(airshipBitcoinPay.IsEnabled(myEntity))
```

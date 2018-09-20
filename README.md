# Airship Go SDK

## Prerequisite

This SDK works with the [Airship Microservice](https://github.com/airshiphq/airship-microservice). Please refer to the its documenation before proceeding.

### Content
- [01 Installation](#01-installation)
- [02 Key Concepts](#02-key-concepts)
- [03 Configuring Flags](#03-configuring-flags)
- [04 Usage](#04-usage)

## 01 Installation

```
go get https://github.com/airshiphq/airship-go
```

## 02 Key Concepts

```
go get https://github.com/airshiphq/airship-go
```

## 03 Configuring Flags

```
go get https://github.com/airshiphq/airship-go
```

## 04 Usage
```
import (
	"fmt"
	airship "github.com/username/library"
)

type User struct {
	ID string `json:"id"`
}

airship.Configure(&airship.Client{
	EnvKey:  "envKey",
	EdgeURL: "http://localhost:5000",
})

airshipBitcoinPay := airship.Flag("bitcoin-pay")

myUser := &User{
	ID: "2",
}

fmt.Println(airshipBitcoinPay.IsEnabled(myUser))
```

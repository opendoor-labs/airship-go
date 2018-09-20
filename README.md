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

In Airship, feature **flags** control traffic to generic objects (called **entities**). The most common type for entities is `User`, but they can also be other things (i.e. `Page`, `Group`, `Team`, `App`, etc.). By default, all entities have the type `User`.

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

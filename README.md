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

In Go, we define different entity types using `struct`s. (e.g., the `User` struct in the usage section)

## 03 Configuring Flags

To configure Airship, we would need to pass a new Client instance.

```
import (
	airship "github.com/username/library"
)

airship.Configure(&airship.Client{
	EnvKey:  "envKey",
	EdgeURL: "http://localhost:5000",
})

```

Here, `envKey` is the environment key you can get from the [**Airship UI**](https://app.airshiphq.com), and the `EdgeURL` points to your [**Airship Microservice**](https://github.com/airshiphq/airship-microservice) URL.

## 04 Usage
```
import (
	"fmt"
	airship "github.com/username/library"
)

// Do configuration (section 03)

type User struct {
	ID string `json:"id"`
}

airshipBitcoinPay := airship.Flag("bitcoin-pay")

myUser := &User{
	ID: "2",
}

fmt.Println(airshipBitcoinPay.IsEnabled(myUser))
fmt.Println(airshipBitcoinPay.IsEligible(myUser))
```

`IsEnabled` returns whether or not a user or entity has access to the feature. It'll return `false` if the flag is not registered with Airship UI.

`IsEligible` returns whether or not a user could be sampled now or in the future within a population associated with the feature flag. It'll return `false` if the flag is not registered with Airship UI.


```
import (
	"fmt"
	airship "github.com/username/library"
)

// Do configuration (section 03)

type User struct {
	ID string `json:"id"`
}

// This is an example expected payload
type Payload struct {
	Foo string `json:"foo"`
}

airshipBitcoinPay := airship.Flag("bitcoin-pay")

myUser := &User{
	ID: "2",
}

fmt.Println(airshipBitcoinPay.GetTreatment(myUser))

var myPayload Payload
err := airshipInstanceBitcoinPay.GetPayload(myUser, &myPayload)
if err != nil {
	fmt.Println(err)
}
fmt.Println(myPayload.Foo)
```

`GetTreatment` returns the treatment codename that is given to the user or entity. It'll return `""` if the flag is not registered with Airship UI.

`GetPayload` returns the JSON payload associated with the treatment for a flag to the user or entity. The way to get the payload is to define a `struct` and to unmarshal the JSON to that struct. It'll return `nil` if the flag is not registered with Airship UI.

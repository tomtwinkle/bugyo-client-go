# bugyo-client-go
Bugyo Cloud Punchmark Client for go

# How to use

```shell
go get github.com/tomtwinkle/bugyo-client-go
```

```go
package main

import (
	"github.com/tomtwinkle/bugyo-client-go"
	"log"
)

func main() {
	config := &bugyoclient.BugyoConfig{
		TenantCode: "<Your Tenant Code>",
		OBCiD: "<Your OBCID>",
		Password: "<Your Password>",
    }
	client, err := bugyoclient.NewClient(config)
	if err != nil {
		log.Fatal(err)
	}
	if err := client.Login(); err != nil {
		log.Fatal(err)
	}
	if err := client.Punchmark(bugyoclient.ClockTypeClockIn); err != nil {
		log.Fatal(err)
	}
}
```

# How to run

- create `.env` file
```config:.env
TENANTCODE=<You Tenant code>
OBCID=<You OBCID>
PASSWORD=<You Password>
```

- go test login

```shell
go test -v ./client/login_test.go
```

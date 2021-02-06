# bugyo-client-go
Bugyo Cloud Punchmark Client for go

# Bugyo CLI Installation

## Windows
- Install with PowerShell

```poweshell
iwr https://github.com/tomtwinkle/bugyo-client-go/releases/download/v0.10.6/bugyoclient_windows_amd64.zip -OutFile bugyoclient.zip && Expand-Archive -Path bugyoclient.zip && rm bugyoclient.zip
cd bugyoclient
.\bugyoclient.exe help
```

# How to use Bugyo CLI

- 出勤

```shell
.\bugyoclient.exe punchmark --type in
or 
.\bugyoclient.exe pm -t in
```

- 退出

```shell
.\bugyoclient.exe punchmark --type out
or 
.\bugyoclient.exe pm -t out
```

- 外出

```shell
.\bugyoclient.exe punchmark --type go
or 
.\bugyoclient.exe pm -t go
```

- 再入

```shell
.\bugyoclient.exe punchmark --type return
or 
.\bugyoclient.exe pm -t return
```

# How to use client

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

# How to develop

- download go modules

```shell
go mod download
```

- tag release

```shell
git fetch
git switch main
git pull origin main

git tag -a vX.X.X -m "release comment"
git push origin vX.X.X
```

# How to run test

- create `.env` file
```config:.env
TENANTCODE=<You Tenant code>
OBCID=<You OBCID>
PASSWORD=<You Password>
```

- go test login

```shell
go test login_test.go
```

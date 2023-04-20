# oav

oav is OpenAPI Validation tool.

inspired code: https://zenn.dev/podhmo/scraps/5dbfa70654f9f0

## install

cli

```shell
go install github.com/kijimaD/oav@main
```

library

```shell
go get github.com/kijimaD/oav@main
```

## Usage

```go
import (
    "github.com/kijimaD/oav/oa"
)

func TestSchema() {}
	file, err := os.Open("openapi.yml")
	if err != nil {
		panic(err)
	}

	baseURL, err := url.Parse("http://localhost:8080") # serversに登録されているホスト名である必要がある
	if err != nil {
		panic(err)
	}

	c := oa.New(os.Stdout, file, *baseURL)
	err = c.Run("/pets", "GET", "{}", "", 200)
	if err != nil {
		log.Fatalf("!! %+v", err)
	}
	err = c.Run("/users", "GET", "{}", "", 200)
	if err != nil {
		log.Fatalf("!! %+v", err)
	}
```

## command

dump schema routes.

```shell
docker-compose up -d
go run . openapi.yml

Endpoint        Method          ID
──────────      ──────────      ──────────
/pets           Get             list_pets
```

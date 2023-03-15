# oav

oav is OpenAPI Validation tool.

inspired code: https://zenn.dev/podhmo/scraps/5dbfa70654f9f0

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

	baseURL, err := url.Parse("http://localhost:8080")
	if err != nil {
		panic(err)
	}

	c := oa.New(os.Stdout, file, *baseURL)
	err = c.Run("/pets")
	if err != nil {
		log.Fatalf("!! %+v", err)
	}
	err = c.Run("/users")
	if err != nil {
		log.Fatalf("!! %+v", err)
	}
```

- チェックするのはserversに登録されているアドレスである必要がある。

## command

dump schema routes.

```shell
docker-compose up -d
go run . openapi.yml
```

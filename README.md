# oav

Portable OpenAPI validation tool.

inspired code: https://zenn.dev/podhmo/scraps/5dbfa70654f9f0

- チェックするのはserversに登録されているアドレスである必要がある。

## Check

```shell
docker-compose up -d
go run . ./openapi.yml http://localhost:8080
```

# go-idempotency

[![CI](https://github.com/tech-nimble/go-idempotency/actions/workflows/ci.yml/badge.svg)](https://github.com/tech-nimble/go-idempotency/actions/workflows/ci.yml)

Idempotency middleware for the [Gin framework](https://github.com/gin-gonic/gin).
Repeated requests carrying the same idempotency key return the cached response
instead of executing the handler again. Backed by
[redis cache](https://github.com/go-redis/cache).

## Install

```sh
go get github.com/tech-nimble/go-idempotency
```

## Usage

```go
import (
	"github.com/redis/go-redis/v9"
	"github.com/tech-nimble/go-idempotency"
)

func main() {
	rdb := redis.NewClient(&redis.Options{Addr: "localhost:6379"})

	checker := idempotency.Initialize(rdb)
	router.POST("/orders", checker.Api, createOrder)
}
```

Clients pass a unique UUID v4 in the `X-Idempotency-Key` header. A served cached
response is marked with `X-Idempotency-Cache: HIT`.

### Custom options

```go
m := idempotency.NewIdempotency(
	idempotency.WithStorage(customStorage),
	idempotency.WithHeaderKey("X-Request-Id"),
)
```

Implement `idempotency.Storage` to plug in a different backend.

## License

[MIT](LICENSE) © Nimble Tech

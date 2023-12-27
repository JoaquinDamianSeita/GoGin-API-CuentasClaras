# GoGin-API-CuentasClaras

run with:

``` bash
go run .
```

test with:

``` bash
go test ./...
```

clean test cache:

``` bash
go clean -testcache
```

Generate wire file:
``` bash
wire gen GoGin-API-CuentasClaras/config
```


.env:

``` go

# Port
PORT=8080

# Application
APPLICATION_NAME=GoGin-API-CuentasClaras
ENVIRONMENT=development | test | production

# Database
DB_DSN="host=HOST user=USER password=PASSWORD dbname=DBNAME port=PORT"
DB_DSN_TEST="host=HOST user=USER password=PASSWORD dbname=DBNAME port=PORT"

# Logging
LOG_LEVEL=DEBUG

# Secret JWT key
SECRET_JWT_KEY="SECRET_JWT_KEY"
```

Live Reload Golang Development With Gin:

``` bash
gin --appPort 3000 --port 8080 --immediate
```


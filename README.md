# Aggregating multiple json sources in Bigtable

This repository contains a small PoC that uses google cloud bigtable to save and aggregate multiple data sources.

# how to run?

```
$ docker compose up -d
$ go run cmd/intaker/main.go
```

# endpoints

- `GET /products/generate`
  - generates random products and saves them to the bigtable emulator
  - `amount` query parameter can be used to change amount of products generated. default is 10.
- `GET /products`
  - returns all products currently in the bigtable emulator
  - `limit` query parameter can be used to limit amount of products returned.

# looking into the bigtable emulator

The `cbt` tool from `gcloud components` can be used to check the data in bigtable

make sure your terminal session is connecting to the emulator
```
 export BIGTABLE_EMULATOR_HOST=localhost:8086
```

See all products
```
cbt -project fake-local-project -instance products read products 
```

See product by id
```
cbt -project fake-local-project -instance products lookup products <product-id>
```
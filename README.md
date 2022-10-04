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

# Sketch of bigtable products table

The table below shows how this PoC stores different things.

Every product gets a row with a key like `product#<product-id>`.

Sources that have a `one-to-one` relation with products, would get a column in the `products` column family. In this example it's only the `product` column containing base product information.

Sources that have a `many-to-one` relation with products, like offers and reviews. Get their own column family and a column per unique entity. So for example the column `offer_1` contains the offer information of the offer with id `1`. This works fine because bigtable is sparsely populated, so all columns that don't have a cell value for a given row don't take up any space.

Sources that don't have a direct relation with products (or `many-to-many` in some way) are a bit more difficult. In the example below this is shown with the `products<->groups` relation. A group contains many products, but a product can also be in many groups.

I'm not sure yet how to save those, but my first attempt is going to be as shown. `products` get a list of group ids and groups get their own rows with keys like `group#<group-id>` and a column `group` in the `product` column family. While reading a problem this means multiple rows need to be read which is less than ideal, but I do want to test it. Alternatively, the group data could be duplicated for every product that is part of the group. But that feels very inefficient in a case where a group contains millions of products and a group gets deleted for example.

| family         | product                               |                                      | offers                           |                                   |                                   | reviews                                         |                        |
|----------------|---------------------------------------|--------------------------------------|----------------------------------|-----------------------------------|-----------------------------------|-------------------------------------------------|------------------------|
|                | product                               | group                                | offer_1                          | offer_2                           | offer_3                           | review_1                                        | review_2               |
| product#123456 | { id: "123456", groupIDs: ["1","2"] } |                                      | { id: "1", price_in_cents: 999 } |                                   |                                   |                                                 | { id: "2", rating: 3 } |
| product#234567 | { id: "234567", groupIDs: ["2"] }     |                                      |                                  | { id: "2", price_in_cents: 1295 } | { id: "3", price_in_cents: 1195 } | { id: "1", rating: 5, message: "cool product" } |                        |
| group#1        |                                       | { id: "1", name: "product group 1" } |                                  |                                   |                                   |                                                 |                        |
| group#2        |                                       | { id: "2", name: "product group 2" } |                                  |                                   |                                   |                                                 |                        |

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

Delete all products
```
cbt -project fake-local-project -instance products deleteallrows products
```
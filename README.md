
Cart API
========

## Overview

This repository contains a Go programming exercise for interview candidates. The
goal is to see how you approach a real world programming task. In this case,
you'll be developing an API for an online shopping cart in the Go programming language.


## Process

1. Create a feature branch.
2. Implement the requirements below as if this were a real Go project to be deployed.
3. Submit a PR to this repository.
4. Discuss your approach and tradeoffs via PR review.


## Requirements

This is a REST API for basic CRUD operations for an online shopping cart. Data
should be persisted in a storage layer which can use Postgres, MongoDB, or a Go
storage library such as [BoltDB](https://github.com/boltdb/bolt) or
[Badger](https://github.com/dgraph-io/badger).

Please include tests for your application.

### Domain Types

The Cart API consists of two simple types: `Cart` and `CartItem`. The `Cart`
holds zero or more `CartItem` objects.


### Create Cart

A new cart should be created and an ID generated. The new empty cart should be returned.

```sh
$ curl -X POST http://localhost:3000/carts -d '{}'
```

```json
{
	"id": 1,
	"items": []
}
```

### Add to cart

A new item should added to an existing cart. Should fail if the cart does not
exist, if the product name is blank, or if the quantity is non-positive. The
new item should be returned.

```sh
$ curl -X POST http://localhost:3000/carts/1/items -d '{
	"product": "Shoes",
	"quantity": 10
}'
```

```json
{
	"id": 1,
	"cart_id": 1,
	"product": "Shoes",
	"quantity": 10
}
```

### Remove from cart

An existing item should be removed from a cart. Should fail if the cart does not
exist or if the item does not exist.

```sh
$ curl -X DELETE http://localhost:3000/carts/1/items/1
```

```json
{}
```


### View cart

An existing cart should be able to be viewed with its items. Should fail if the
cart does not exist.

```sh
$ curl http://localhost:3000/carts/1
```

```json
{
	"id": 1,
	"items": [
		{
			"id": 1,
			"cart_id": 1,
			"product": "Shoes",
			"quantity": 10
		},
		{
			"id": 2,
			"cart_id": 1,
			"product": "Socks",
			"quantity": 5
		}
	]
}
```


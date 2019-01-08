# Topaz Quick Start Guide

To get up and running with Topaz API as quickly as possible, follow this guide.

You'll create a new account, log in, create a new app, generate a token for that app, and use it to interact with Topaz.

# User Registration

## Create your new account

### Request

`POST /users`

#### Headers

* `'Content-Type: application/json'`

#### Body

```json
{
	"email": "adam@topaz.io",
	"password": "hunter2",
	"name": "Adam Gall"
}
```

### Response

```json
{
    "ID": 3,
    "CreatedAt": "2018-12-18T15:15:01.069166-05:00",
    "UpdatedAt": "2018-12-18T15:15:01.069166-05:00",
    "DeletedAt": null,
    "name": "Adam Gall",
    "email": "adam@topaz.io",
    "password": "$2a$14$AZ3T4SiWjN4Yf3wALygXVusRm17LYTJ5FtFKI4625auYXKC9DigI6"
}
```

# User Authentication

## Log in to your new account

### Request

`POST /auth/admin/token`

#### Headers

* `'Content-Type: application/json'`

#### Body

```json
{
	"email": "adam@topaz.io",
	"password": "hunter2"
}
```

### Response

```json
{
    "token": "eyJhbGciOiJSUzUxMiIsInR5cCI6IkpXVCJ9.eyJ1c2VySWQiOiIzIiwiZXhwIjoxNTQ1NDI1MTY0LCJpYXQiOjE1NDUxNjU5NjQsInN1YiI6IjMifQ.ZpArTvZig7f9Tj9Ztb-K3IzpPE4rP15wb1_fl8FtaUGXuHYhiS9X4XomPP6WqUYRDd1QHlUBaw3Xd36aPYsRd5sNBSeyf4fv23iYAV5zzmpl5mylyZDSu-9aN7ttYTqxgd9bu5ck_nppFGmmM3cTnkTqsUoVJR-TlSpc7pTNvScX_RyZ8gk-KrDQq9xGEOZw3WQ48FcN0pLZlnu6_e84-OTwpIAoSvxIWYZyr-3i3DCl1ZzxsFJeX0Cu9Txs3dbTq6tisHTnJPpf9vqxS38Koc-PLkVEzmckIu3yavKJH7FbEB1ZImNZPRLbxGEZh9Mce0TG9drafIw7X4nvLfpY7g"
}
```

You'll need this token to create new apps, and create new app tokens, as described below.

# App Creation

## Create a new app context

### Request

`POST /apps`

#### Headers

* `'Content-Type: application/json'`
* `'Authorization: Bearer eyJhbGciOiJSUzUxMiIsInR5cCI6IkpXVCJ9.eyJ1c2VySWQiOiIzIiwiZXhwIjoxNTQ1NDI1MTY0LCJpYXQiOjE1NDUxNjU5NjQsInN1YiI6IjMifQ.ZpArTvZig7f9Tj9Ztb-K3IzpPE4rP15wb1_fl8FtaUGXuHYhiS9X4XomPP6WqUYRDd1QHlUBaw3Xd36aPYsRd5sNBSeyf4fv23iYAV5zzmpl5mylyZDSu-9aN7ttYTqxgd9bu5ck_nppFGmmM3cTnkTqsUoVJR-TlSpc7pTNvScX_RyZ8gk-KrDQq9xGEOZw3WQ48FcN0pLZlnu6_e84-OTwpIAoSvxIWYZyr-3i3DCl1ZzxsFJeX0Cu9Txs3dbTq6tisHTnJPpf9vqxS38Koc-PLkVEzmckIu3yavKJH7FbEB1ZImNZPRLbxGEZh9Mce0TG9drafIw7X4nvLfpY7g'`

#### Body

```json
{
	"name": "Supply Chain Workflow",
	"interval": 3600
}
```

### Response

```json
{
    "ID": 1,
    "CreatedAt": "2018-12-18T16:03:33.113418-05:00",
    "UpdatedAt": "2018-12-18T16:03:33.113418-05:00",
    "DeletedAt": null,
    "interval": 3600,
    "name": "Supply Chain Workflow",
    "ethAddress": "0x8aba912417dE237b7Df401C437cCad0846a2ef76",
    "userId": 1
}
```

Take note of this App's `ID`: `1`, which will be used to generate an app-specific token, as described in the next section.

# App Authentication

## Create a new app token

### Request

`POST /auth/app/token`

#### Headers

* `'Content-Type: application/json'`
* `'Authorization: Bearer eyJhbGciOiJSUzUxMiIsInR5cCI6IkpXVCJ9.eyJ1c2VySWQiOiIzIiwiZXhwIjoxNTQ1NDI1MTY0LCJpYXQiOjE1NDUxNjU5NjQsInN1YiI6IjMifQ.ZpArTvZig7f9Tj9Ztb-K3IzpPE4rP15wb1_fl8FtaUGXuHYhiS9X4XomPP6WqUYRDd1QHlUBaw3Xd36aPYsRd5sNBSeyf4fv23iYAV5zzmpl5mylyZDSu-9aN7ttYTqxgd9bu5ck_nppFGmmM3cTnkTqsUoVJR-TlSpc7pTNvScX_RyZ8gk-KrDQq9xGEOZw3WQ48FcN0pLZlnu6_e84-OTwpIAoSvxIWYZyr-3i3DCl1ZzxsFJeX0Cu9Txs3dbTq6tisHTnJPpf9vqxS38Koc-PLkVEzmckIu3yavKJH7FbEB1ZImNZPRLbxGEZh9Mce0TG9drafIw7X4nvLfpY7g'`

#### Body

```json
{
	"ID": 1
}
```

### Response

```json
{
    "token": "eyJhbGciOiJSUzUxMiIsInR5cCI6IkpXVCJ9.eyJhcHBJZCI6IjIiLCJleHAiOjE1NDU0MjY2MTEsImlhdCI6MTU0NTE2NzQxMSwic3ViIjoiMiJ9.L6fYvVSp1WNmazPk4Rwo7pLxiXIJiYv0U5vc2hhHWf7zk7f3L7kCsVwE7EJUFFINqneQ0EW5gklBthEaVWl3Ven10dvnpGNgL5MtlXyzdXnRf5duc2qeVBLRUD8V8JJsAt28EVBu-rU27thWAtod0kLgDnSmaoOmqEAF4uizD5dvOcKAH9-rLwEDsiYFrsO8AI23Wdjcg_w7AVYz_lZteZXk9J5KKEmohv3a6nlOblFdHBGrsv8kgnyX4OYB9wfJOXCvuD5a_WGbfjX590iVe9pR7Z7WaYUd5gmRSe0uhOWRYpT5O72rQcvcv-FT0pa59SFM6HZb1kYQJGE5RRg_fw"
}
```

This token will be used for all subsequent calls to Topaz API, in order to use the power of public blockchains to timestamp and report on data within a specific app context.

# App Usage

## Trust your data

Use the `/trust` endpoint to send a hash of your business-valuable data to Topaz to be processed. The POST status code is used for new objects which have not yet been seen by Topaz, and does not require a `UUID` parameter.

### Request

`POST /trust`

#### Headers

* `'Authorization: Bearer eyJhbGciOiJSUzUxMiIsInR5cCI6IkpXVCJ9.eyJhcHBJZCI6IjIiLCJleHAiOjE1NDU0MjY2MTEsImlhdCI6MTU0NTE2NzQxMSwic3ViIjoiMiJ9.L6fYvVSp1WNmazPk4Rwo7pLxiXIJiYv0U5vc2hhHWf7zk7f3L7kCsVwE7EJUFFINqneQ0EW5gklBthEaVWl3Ven10dvnpGNgL5MtlXyzdXnRf5duc2qeVBLRUD8V8JJsAt28EVBu-rU27thWAtod0kLgDnSmaoOmqEAF4uizD5dvOcKAH9-rLwEDsiYFrsO8AI23Wdjcg_w7AVYz_lZteZXk9J5KKEmohv3a6nlOblFdHBGrsv8kgnyX4OYB9wfJOXCvuD5a_WGbfjX590iVe9pR7Z7WaYUd5gmRSe0uhOWRYpT5O72rQcvcv-FT0pa59SFM6HZb1kYQJGE5RRg_fw'`

#### Body

```json
{
	"hash": "17fee8f6b4c18e3ceb93362e67551aadb3b5772264e6c7523613f87f10342592"
}
```

### Response

```json
{
    "ID": 71,
    "CreatedAt": "2018-12-18T16:16:53.061753-05:00",
    "UpdatedAt": "2018-12-18T16:16:53.061753-05:00",
    "DeletedAt": null,
    "uuid": "cb01ed8c-2fb4-4b42-b449-0e24c4782c83",
    "hash": "17fee8f6b4c18e3ceb93362e67551aadb3b5772264e6c7523613f87f10342592",
    "unixTimestamp": 1545167813,
    "appId": 1,
    "proofId": null
}
```

## Update your previously trusted data

Use the `/trust` endpoint to send a hash of your business-valuable data to Topaz to be processed. The PUT/PATCH status codes are used for existing objects which have already been previously trusted by Topaz, and require the UUID of the object to be passed to the endpoint.

### Request

`PUT/PATCH /trust`

#### Headers

* `'Authorization: Bearer eyJhbGciOiJSUzUxMiIsInR5cCI6IkpXVCJ9.eyJhcHBJZCI6IjIiLCJleHAiOjE1NDU0MjY2MTEsImlhdCI6MTU0NTE2NzQxMSwic3ViIjoiMiJ9.L6fYvVSp1WNmazPk4Rwo7pLxiXIJiYv0U5vc2hhHWf7zk7f3L7kCsVwE7EJUFFINqneQ0EW5gklBthEaVWl3Ven10dvnpGNgL5MtlXyzdXnRf5duc2qeVBLRUD8V8JJsAt28EVBu-rU27thWAtod0kLgDnSmaoOmqEAF4uizD5dvOcKAH9-rLwEDsiYFrsO8AI23Wdjcg_w7AVYz_lZteZXk9J5KKEmohv3a6nlOblFdHBGrsv8kgnyX4OYB9wfJOXCvuD5a_WGbfjX590iVe9pR7Z7WaYUd5gmRSe0uhOWRYpT5O72rQcvcv-FT0pa59SFM6HZb1kYQJGE5RRg_fw'`

#### Body

```json
{
    "uuid": "cb01ed8c-2fb4-4b42-b449-0e24c4782c83",
	"hash": "85cf5d1d911e1dde12d8f701b85c69591e1e19e1b1c642d54b4a57fc6a5fbee7"
}
```

### Response

```json
{
    "ID": 72,
    "CreatedAt": "2018-12-18T16:16:53.061753-05:00",
    "UpdatedAt": "2018-12-18T16:16:53.061753-05:00",
    "DeletedAt": null,
    "uuid": "cb01ed8c-2fb4-4b42-b449-0e24c4782c83",
    "hash": "85cf5d1d911e1dde12d8f701b85c69591e1e19e1b1c642d54b4a57fc6a5fbee7",
    "unixTimestamp": 1545169876,
    "appId": 1,
    "proofId": null
}
```

## Verify your data

Use the `/verify` endpoint to check whether or not a particular hash of data has been seen by Topaz API.

### Request

`GET /verify/{hash}`

example: `/verify/17fee8f6b4c18e3ceb93362e67551aadb3b5772264e6c7523613f87f10342592`

#### Headers

* `'Authorization: Bearer eyJhbGciOiJSUzUxMiIsInR5cCI6IkpXVCJ9.eyJhcHBJZCI6IjIiLCJleHAiOjE1NDU0MjY2MTEsImlhdCI6MTU0NTE2NzQxMSwic3ViIjoiMiJ9.L6fYvVSp1WNmazPk4Rwo7pLxiXIJiYv0U5vc2hhHWf7zk7f3L7kCsVwE7EJUFFINqneQ0EW5gklBthEaVWl3Ven10dvnpGNgL5MtlXyzdXnRf5duc2qeVBLRUD8V8JJsAt28EVBu-rU27thWAtod0kLgDnSmaoOmqEAF4uizD5dvOcKAH9-rLwEDsiYFrsO8AI23Wdjcg_w7AVYz_lZteZXk9J5KKEmohv3a6nlOblFdHBGrsv8kgnyX4OYB9wfJOXCvuD5a_WGbfjX590iVe9pR7Z7WaYUd5gmRSe0uhOWRYpT5O72rQcvcv-FT0pa59SFM6HZb1kYQJGE5RRg_fw'`

### Response

```json
[
    {
        "ID": 71,
        "CreatedAt": "2018-12-18T16:16:53.061753Z",
        "UpdatedAt": "2018-12-18T16:25:14.690197Z",
        "DeletedAt": null,
        "uuid": "cb01ed8c-2fb4-4b42-b449-0e24c4782c83",
        "hash": "17fee8f6b4c18e3ceb93362e67551aadb3b5772264e6c7523613f87f10342592",
        "unixTimestamp": 1545167813,
        "appId": 1,
        "proofId": 24,
        "proof": {
            "ID": 24,
            "CreatedAt": "2018-12-18T16:25:14.67507Z",
            "UpdatedAt": "2018-12-18T16:25:14.67507Z",
            "DeletedAt": null,
            "merkleRoot": "QmYoB7DqNkQ5aaSuJYVeNATeWYdaSk3trugK7X6SwGKBdp",
            "ethTransaction": "0xd4fd388f808993612627644add6d4ade7591865be858b009741f010c2ca2a852",
            "validStructure": true,
            "batchId": 8894,
            "batch": {
                "ID": 8894,
                "CreatedAt": "2018-12-18T16:25:14.657868Z",
                "UpdatedAt": "2018-12-18T16:25:14.657868Z",
                "DeletedAt": null,
                "unixTimestamp": 1545168314,
                "appId": 1
            },
            "objects": [
                {
                    "ID": 71,
                    "CreatedAt": "2018-12-18T16:16:53.061753Z",
                    "UpdatedAt": "2018-12-18T16:25:14.690197Z",
                    "DeletedAt": null,
                    "uuid": "cb01ed8c-2fb4-4b42-b449-0e24c4782c83",
                    "hash": "17fee8f6b4c18e3ceb93362e67551aadb3b5772264e6c7523613f87f10342592",
                    "unixTimestamp": 1545167813,
                    "appId": 1,
                    "proofId": 24
                }
            ]
        }
    }
]
```

## Trace your data

Use the `/trace` endpoint to track the changes of a UUID over time

### Request

`GET /trace/{uuid}`

example: `/trace/cb01ed8c-2fb4-4b42-b449-0e24c4782c83`

#### Headers

* `'Authorization: Bearer eyJhbGciOiJSUzUxMiIsInR5cCI6IkpXVCJ9.eyJhcHBJZCI6IjIiLCJleHAiOjE1NDU0MjY2MTEsImlhdCI6MTU0NTE2NzQxMSwic3ViIjoiMiJ9.L6fYvVSp1WNmazPk4Rwo7pLxiXIJiYv0U5vc2hhHWf7zk7f3L7kCsVwE7EJUFFINqneQ0EW5gklBthEaVWl3Ven10dvnpGNgL5MtlXyzdXnRf5duc2qeVBLRUD8V8JJsAt28EVBu-rU27thWAtod0kLgDnSmaoOmqEAF4uizD5dvOcKAH9-rLwEDsiYFrsO8AI23Wdjcg_w7AVYz_lZteZXk9J5KKEmohv3a6nlOblFdHBGrsv8kgnyX4OYB9wfJOXCvuD5a_WGbfjX590iVe9pR7Z7WaYUd5gmRSe0uhOWRYpT5O72rQcvcv-FT0pa59SFM6HZb1kYQJGE5RRg_fw'`

### Response

```json
[
    {
        "ID": 71,
        "CreatedAt": "2018-12-18T16:16:53.061753Z",
        "UpdatedAt": "2018-12-18T16:25:14.690197Z",
        "DeletedAt": null,
        "uuid": "cb01ed8c-2fb4-4b42-b449-0e24c4782c83",
        "hash": "17fee8f6b4c18e3ceb93362e67551aadb3b5772264e6c7523613f87f10342592",
        "unixTimestamp": 1545167813,
        "appId": 1,
        "proofId": 24,
        "proof": {
            "ID": 24,
            "CreatedAt": "2018-12-18T16:25:14.67507Z",
            "UpdatedAt": "2018-12-18T16:25:14.67507Z",
            "DeletedAt": null,
            "merkleRoot": "QmYoB7DqNkQ5aaSuJYVeNATeWYdaSk3trugK7X6SwGKBdp",
            "ethTransaction": "0xd4fd388f808993612627644add6d4ade7591865be858b009741f010c2ca2a852",
            "validStructure": true,
            "batchId": 8894,
            "batch": {
                "ID": 8894,
                "CreatedAt": "2018-12-18T16:25:14.657868Z",
                "UpdatedAt": "2018-12-18T16:25:14.657868Z",
                "DeletedAt": null,
                "unixTimestamp": 1545168314,
                "appId": 1
            },
            "objects": [
                {
                    "ID": 71,
                    "CreatedAt": "2018-12-18T16:16:53.061753Z",
                    "UpdatedAt": "2018-12-18T16:25:14.690197Z",
                    "DeletedAt": null,
                    "uuid": "cb01ed8c-2fb4-4b42-b449-0e24c4782c83",
                    "hash": "17fee8f6b4c18e3ceb93362e67551aadb3b5772264e6c7523613f87f10342592",
                    "unixTimestamp": 1545167813,
                    "appId": 1,
                    "proofId": 24
                }
            ]
        }
    },
    {
        "ID": 72,
        "CreatedAt": "2018-12-18T21:51:16.061753-05:00",
        "UpdatedAt": "2018-12-18T21:51:16.061753-05:00",
        "DeletedAt": null,
        "uuid": "cb01ed8c-2fb4-4b42-b449-0e24c4782c83",
        "hash": "85cf5d1d911e1dde12d8f701b85c69591e1e19e1b1c642d54b4a57fc6a5fbee7",
        "unixTimestamp": 1545169876,
        "appId": 1,
        "proofId": 25,
        "proof": {
            "ID": 25,
            "CreatedAt": "2018-12-18T21:16:00.00007Z",
            "UpdatedAt": "2018-12-18T21:16:00.00007Z",
            "DeletedAt": null,
            "merkleRoot": "QmYoB7DqNkQ5aaSuJYVeNATeWYdaSk3trugK7X6SwGKBdp",
            "ethTransaction": "0x1a62c1fbbeb49946cb13337971056aa0e58c46a77192d8eb4bc5f23a64200fbe",
            "validStructure": true,
            "batchId": 8895,
            "batch": {
                "ID": 8895,
                "CreatedAt": "2018-12-18T21:16:00.00007Z",
                "UpdatedAt": "2018-12-18T21:16:00.00007Z",
                "DeletedAt": null,
                "unixTimestamp": 1545168314,
                "appId": 1
            },
            "objects": [
                {
                    "ID": 72,
                    "CreatedAt": "2018-12-18T21:51:16.061753-05:00",
                    "UpdatedAt": "2018-12-18T21:51:16.061753-05:00",
                    "DeletedAt": null,
                    "uuid": "cb01ed8c-2fb4-4b42-b449-0e24c4782c83",
                    "hash": "85cf5d1d911e1dde12d8f701b85c69591e1e19e1b1c642d54b4a57fc6a5fbee7",
                    "unixTimestamp": 1545169876,
                    "appId": 1,
                    "proofId": 25,
                }
            ]
        }
    }
]
```

## Report on your data

Run time-based reports with the `/report` endpoint.

### Request

`POST /report`

#### Headers

* `'Content-Type: application/json'`
* `'Authorization: Bearer eyJhbGciOiJSUzUxMiIsInR5cCI6IkpXVCJ9.eyJhcHBJZCI6IjIiLCJleHAiOjE1NDU0MjY2MTEsImlhdCI6MTU0NTE2NzQxMSwic3ViIjoiMiJ9.L6fYvVSp1WNmazPk4Rwo7pLxiXIJiYv0U5vc2hhHWf7zk7f3L7kCsVwE7EJUFFINqneQ0EW5gklBthEaVWl3Ven10dvnpGNgL5MtlXyzdXnRf5duc2qeVBLRUD8V8JJsAt28EVBu-rU27thWAtod0kLgDnSmaoOmqEAF4uizD5dvOcKAH9-rLwEDsiYFrsO8AI23Wdjcg_w7AVYz_lZteZXk9J5KKEmohv3a6nlOblFdHBGrsv8kgnyX4OYB9wfJOXCvuD5a_WGbfjX590iVe9pR7Z7WaYUd5gmRSe0uhOWRYpT5O72rQcvcv-FT0pa59SFM6HZb1kYQJGE5RRg_fw'`

#### Body

Use Unix timestamps to create a range in which to run reports.

```json
{
	"start": 0,
	"end": 1545167814
}
```

### Response

```json
[
    {
        "ID": 72,
        "CreatedAt": "2018-12-18T21:51:16.061753-05:00",
        "UpdatedAt": "2018-12-18T21:51:16.061753-05:00",
        "DeletedAt": null,
        "uuid": "cb01ed8c-2fb4-4b42-b449-0e24c4782c83",
        "hash": "85cf5d1d911e1dde12d8f701b85c69591e1e19e1b1c642d54b4a57fc6a5fbee7",
        "unixTimestamp": 1545169876,
        "appId": 1,
        "proofId": 25,
        "proof": {
            "ID": 25,
            "CreatedAt": "2018-12-18T21:16:00.00007Z",
            "UpdatedAt": "2018-12-18T21:16:00.00007Z",
            "DeletedAt": null,
            "merkleRoot": "QmYoB7DqNkQ5aaSuJYVeNATeWYdaSk3trugK7X6SwGKBdp",
            "ethTransaction": "0x1a62c1fbbeb49946cb13337971056aa0e58c46a77192d8eb4bc5f23a64200fbe",
            "validStructure": true,
            "batchId": 8895,
            "batch": {
                "ID": 8895,
                "CreatedAt": "2018-12-18T21:16:00.00007Z",
                "UpdatedAt": "2018-12-18T21:16:00.00007Z",
                "DeletedAt": null,
                "unixTimestamp": 1545168314,
                "appId": 1
            }
        }
    }
]
```

This response returns an array of objects which were sent to Topaz API in the given timeframe. Nested within the objects are their corresponding on-chain proofs (who's merkle root is the key identifier which is stored on-chain), and all of the objects which contributed to that specific proof.

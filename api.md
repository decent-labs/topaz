# Topaz Quick Start Guide

To get up and running with Topaz API as quickly as possible, follow this guide.

You'll create a new account, log in, create a new app, create objects, create hashes, and view proofs.

## Note

All API requests should be prefixed with `/v1`, indicating that you're targeting version 1 of our API.

Topaz API will follow Semantic Versioning (https://semver.org), so expect that any breaking changes will be versioned under a new route prefix.

## API Tokens

Get an API token at topaz.io. Create a new account log in, and generate a new Token. A valid API token is necessary to access the API.

# Apps

## Create a new app

To create a new app, post to the `/apps` endpoint. You'll need to name the app, and pass in an "interval", in seconds. This interval represents the amount of time between blockchain transactions that will occur as you're adding hashes to objects in this app.

In effect, it's the "resolution" at which you're comfortable proving that your data exists.

### Request

`POST /apps`

#### Headers

* `'Content-Type: application/json'`
* `'Authorization: Bearer eyJhbGciOiJSUzUxMiIsInR5cCI6IkpXVCJ9...'`

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
    "id": "d122b8de-c48f-4586-a2b6-42cb337341b6",
    "interval": 3600,
    "name": "Supply Chain Workflow",
    "userId": "3210c665-2dab-41b8-b688-83a73fc2f127"
}
```

## Get all apps

This endpoint will return all apps registered for a user

### Request

`GET /apps`

#### Headers

* `'Content-Type: application/json'`
* `'Authorization: Bearer eyJhbGciOiJSUzUxMiIsInR5cCI6IkpXVCJ9...'`

### Response

```json
[
    {
        "id": "d122b8de-c48f-4586-a2b6-42cb337341b6",
        "interval": 3600,
        "name": "Supply Chain Workflow",
        "userId": "3210c665-2dab-41b8-b688-83a73fc2f127"
    },
    {
        "id": "1780363b-8dd9-48d0-963a-3fb7e13130f8",
        "interval": 3600,
        "name": "Supply Chain Workflow 2",
        "userId": "3210c665-2dab-41b8-b688-83a73fc2f127"
    }
]
```

## Get a single app

This endpoint will return details about a single app, given the `appId` passed in as the request parameter

### Request

`GET /apps/{appId}`

#### Headers

* `'Content-Type: application/json'`
* `'Authorization: Bearer eyJhbGciOiJSUzUxMiIsInR5cCI6IkpXVCJ9...'`

### Response

```json
{
    "id": "d122b8de-c48f-4586-a2b6-42cb337341b6",
    "interval": 3600,
    "name": "Supply Chain Workflow",
    "userId": "3210c665-2dab-41b8-b688-83a73fc2f127"
}
```

# Objects

Objects are created within the context of an app, so all requests in this section are prefixed with `/apps/{appId}`.

## Create a new object

This endpoint doesn't take any body data in its request.

### Request

`POST /apps/{appId}/objects`

#### Headers

* `'Content-Type: application/json'`
* `'Authorization: Bearer eyJhbGciOiJSUzUxMiIsInR5cCI6IkpXVCJ9...'`

### Response

```json
{
    "id": "566f8a83-bdf9-4397-8ca4-698c29ca656d",
    "appId": "d122b8de-c48f-4586-a2b6-42cb337341b6"
}
```

## Get all objects

This endpoint will return all objects registered for an app

### Request

`GET /apps/{appId}/objects`

#### Headers

* `'Content-Type: application/json'`
* `'Authorization: Bearer eyJhbGciOiJSUzUxMiIsInR5cCI6IkpXVCJ9...'`

### Response

```json
[
    {
        "id": "566f8a83-bdf9-4397-8ca4-698c29ca656d",
        "appId": "d122b8de-c48f-4586-a2b6-42cb337341b6"
    },
    {
        "id": "64e6b247-4043-46d7-bc62-3d5372995551",
        "appId": "d122b8de-c48f-4586-a2b6-42cb337341b6"
    }
]
```

## Get a single object

This endpoint will return details about a single object, given the `objectId` passed in as the request parameter

### Request

`GET /apps/{appId}/objects/{objectId}`

#### Headers

* `'Content-Type: application/json'`
* `'Authorization: Bearer eyJhbGciOiJSUzUxMiIsInR5cCI6IkpXVCJ9...'`

### Response

```json
{
    "id": "566f8a83-bdf9-4397-8ca4-698c29ca656d",
    "appId": "d122b8de-c48f-4586-a2b6-42cb337341b6"
}
```

# Hashes

Hashes are created within the context of an object, so all requests in this section are prefixed with `/apps/{appId}/objects/{objectId}`.

## Create a new hash

This endpoint inputs a json object in it's body which defines a `data` property. This `data` should be a hex-encoded SHA256 hash.

### Request

`POST /apps/{appId}/objects/{objectId}/hashes`

#### Headers

* `'Content-Type: application/json'`
* `'Authorization: Bearer eyJhbGciOiJSUzUxMiIsInR5cCI6IkpXVCJ9...'`

### Request

```json
{
	"hash": "e17b691a22460b7540a75f524e4f093faa0f2df08a47b3fb3d0c220606db6e7d"
}
```

### Response

```json
{
    "id": "bed9d7cd-9ea7-4d81-a735-18e5239e8010",
    "unixTimestamp": 1550610197,
    "objectId": "566f8a83-bdf9-4397-8ca4-698c29ca656d",
    "proofId": null,
    "hash": "e17b691a22460b7540a75f524e4f093faa0f2df08a47b3fb3d0c220606db6e7d"
}
```

Notice how the `proofId` is null. It will become populated during the next time that this app's `interval` ticks over, and a blockchain transaction is created.

## Get all hashes

This endpoint will return all hashes which have been associated with an object

### Request

`GET /apps/{appId}/objects/{objectId}/hashes`

#### Headers

* `'Content-Type: application/json'`
* `'Authorization: Bearer eyJhbGciOiJSUzUxMiIsInR5cCI6IkpXVCJ9...'`

### Response

```json
[
    {
        "id": "bed9d7cd-9ea7-4d81-a735-18e5239e8010",
        "unixTimestamp": 1550610197,
        "objectId": "566f8a83-bdf9-4397-8ca4-698c29ca656d",
        "proofId": "be199086-670b-48e8-a6a0-7fd506f97ff0",
        "hash": "e17b691a22460b7540a75f524e4f093faa0f2df08a47b3fb3d0c220606db6e7d"
    },
    {
        "id": "1ea94df9-4069-404c-82e6-718850798354",
        "unixTimestamp": 1550610332,
        "objectId": "566f8a83-bdf9-4397-8ca4-698c29ca656d",
        "proofId": null,
        "hash": "b7b85ea166f915f514e1950e3195ea87cfedb87934eee8b6c4a80891e190f819"
    }
]
```

## Get a single hash

This endpoint will return details about a single hash, given the `hashId` passed in as the request parameter

### Request

`GET /apps/{appId}/objects/{objectId}/hashes/{hashId}`

#### Headers

* `'Content-Type: application/json'`
* `'Authorization: Bearer eyJhbGciOiJSUzUxMiIsInR5cCI6IkpXVCJ9...'`

### Response

```json
{
    "id": "bed9d7cd-9ea7-4d81-a735-18e5239e8010",
    "unixTimestamp": 1550610197,
    "objectId": "566f8a83-bdf9-4397-8ca4-698c29ca656d",
    "proofId": "be199086-670b-48e8-a6a0-7fd506f97ff0",
    "hash": "e17b691a22460b7540a75f524e4f093faa0f2df08a47b3fb3d0c220606db6e7d"
}
```

See how the `proofId` is now populated. This means that a proof has been created and stamped into the blockchain. Congratulations, your data is now provable for all of eternity!

# Proofs

Proofs are automatically created within the context of an app, so all requests in this section are prefixed with `/apps/{appId}`.

## Get all proofs

This endpoint will return all proofs associated with an app

### Request

`GET /apps/{appId}/proofs`

#### Headers

* `'Content-Type: application/json'`
* `'Authorization: Bearer eyJhbGciOiJSUzUxMiIsInR5cCI6IkpXVCJ9...'`

### Response

```json
[
    {
        "id": "be199086-670b-48e8-a6a0-7fd506f97ff0",
        "merkleRoot": "QmRXPkj6np5mJe9o6a7NJRzaY9esaJ7yvbtnk4XpJd5dgX",
        "ethTransaction": "0xcc9b0aae939acd8aa727b71864e8a803888bd6f2e91b5aa5b4a351dbce9f5eeb",
        "unixTimestamp": 1550610221,
        "appId": "d122b8de-c48f-4586-a2b6-42cb337341b6"
    },
    {
        "id": "01041352-6faa-425e-ac2c-babcc464f2ff",
        "merkleRoot": "QmaFtr1XH8XfpmBqJkoUXqG3jbd452YtWdwcEWK9Py6QAx",
        "ethTransaction": "0xa3939cbf1ae04451cd59ddfd3d49477506d7e320a05762d231d7fecbff2d8fa0",
        "unixTimestamp": 1550610551,
        "appId": "d122b8de-c48f-4586-a2b6-42cb337341b6"
    }
]
```

A proof has a `merkleRoot`, and an `ethTransaction`. These two pieces of data allow you to confirm that your data is secured in the blockchain.

We combine all hashes sent to Topaz for a given app during an `interval` timeframe, create a Merkel Tree out of those hashes, compute the root hash (`merkleRoot`), craft an Ethereum transaction using that `merkleRoot` as input data, and broadcast that transaction to be mined into a block, solidifying it's existence into an untamperable data structure.

## Get a single proof

This endpoint will return details about a single proof, given the `proofId` passed in as the request parameter

### Request

`GET /apps/{appId}/proofs/{proofId}`

#### Headers

* `'Content-Type: application/json'`
* `'Authorization: Bearer eyJhbGciOiJSUzUxMiIsInR5cCI6IkpXVCJ9...'`

### Response

```json
{
    "id": "be199086-670b-48e8-a6a0-7fd506f97ff0",
    "merkleRoot": "QmRXPkj6np5mJe9o6a7NJRzaY9esaJ7yvbtnk4XpJd5dgX",
    "ethTransaction": "0xcc9b0aae939acd8aa727b71864e8a803888bd6f2e91b5aa5b4a351dbce9f5eeb",
    "unixTimestamp": 1550610221,
    "appId": "d122b8de-c48f-4586-a2b6-42cb337341b6"
}
```
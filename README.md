# Topaz

Topaz is a platform for api-based blockchain interaction

## Getting Started

### Sync dependencies

This project uses golang modules, so simply building or running the project(s) will download and install the correct version of dependencies, based on `go.mod` and `go.sum`.

### Set up your `.env`

```sh
$ cp .env.example .env
```

Update appropriately.

### Spin up external services

```sh
$ docker-compose up -d
```

The ganache-cli node will save all local blockchain and account state into the `/chainstate` folder of this repo, the contents of which are gitignored. This lets it stay persisted between `docker-compose` runs. If you need to wipe your local blockchain, delete all the files in `/chainstate` (except for `.keep`)

### Run DB migrations if you have to

```sh
$ cd migrate
$ make run
```

### Start the API

```sh
$ cd api
$ make run
```

### Begin the batch process

In a new terminal window

```sh
$ cd batch
$ make run
```

### Use the API

Check out `api.md` to learn how the API functions.

## Where you at?

Made with :heart: in:
* British Airways Flight 1575
* Miami, Florida
* Positano, Italy
* Santa Cruz, California
* Cleveland, Ohio
* Boston, Massachusetts
* Copenhagen, Denmark
* Prague, Czech Republic
* Stockholm, Sweden
* Berlin, Germany
* Amsterdam, Netherlands

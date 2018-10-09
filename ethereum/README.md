# topaz-ethereum

Connecting Topaz with Ethereum. Used for deploying and interacting with Capture Contracts.

## Getting Started

* Download and install the latest version of [docker](https://www.docker.com/get-started)
* Download and install the latest version of [go](https://golang.org/dl/)
* Download and install the latest version of [terraform](https://terraform.io)

Wherever possible, this application reflects the following standards:

* [Standard Go Project Layout](https://github.com/golang-standards/project-layout)
* [The Twelve-Factor App](https://12factor.net/)
* [Go: Best Practices for Production Environments](https://peter.bourgon.org/go-in-production/)

## Local Development

```
$ CONN=https://ropsten.infura.io PRIVKEY=your-private-key-here go run *.go
```

will start the API, accepting POST requests at /deploy

## Test

```
$ cd contracts
$ go test
```

## Where You At

When you make a commit to the project from somewhere new, add it to the list!

Made with :heart: in:
* Cleveland, OH

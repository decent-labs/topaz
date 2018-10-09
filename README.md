# topaz

## getting started

```sh
$ docker-compose up
```

```sh
$ curl --request POST -d "your-data-here" http://localhost:8081
```

## todo

### big stuff

* replace HTTP wherever possible
* project-wide best practices such as error handling

### deployments

* terraform
* k8s

### dispatch

* decide on long term timing (cron) solution
* deal with postgres race condition (ask Parker about this)
* activate 'flush' service asynchronously via worker pool

### ethereum

* add 'ethereum' service to `docker-compose.yml`
* integrate in stack

### flush

* deploy as worker

### ipfs

* add cluster / resilience infrastructure

### postgres

* decide on migrations system (ask Parker about this)

### store

* add authentication

## where you at?

Made with :heart: in:
* British Airways Flight 1575
* Miami, FL
* Positano, IT
* Santa Cruz, CA
* Cleveland, OH

# chaosmonkey

[![Build Status](https://travis-ci.org/mlafeldt/chaosmonkey.svg?branch=master)](https://travis-ci.org/mlafeldt/chaosmonkey)
[![GoDoc](https://godoc.org/github.com/mlafeldt/chaosmonkey?status.svg)](https://godoc.org/github.com/mlafeldt/chaosmonkey)

Go client to the [Chaos Monkey REST API](https://github.com/Netflix/SimianArmy/wiki/REST) that can be used to trigger and retrieve chaos events.

This project was started for the purpose of controlled failure-injection during GameDay events.

## Prerequisites

First of all, you need a running [Simian Army](https://github.com/Netflix/SimianArmy) that exposes its REST API via HTTP.

In order to trigger chaos events via the API, Chaos Monkey must be unleashed and on-demand termination must be enabled via these configuration properties:

```
simianarmy.chaos.leashed = false
simianarmy.chaos.terminateOndemand.enabled = true
```

## Go library

To install the `chaosmonkey` library from source:

```bash
go get github.com/mlafeldt/chaosmonkey
```

For usage and examples, see the [Godoc documentation](https://godoc.org/github.com/mlafeldt/chaosmonkey).

## CLI

In addition to the library, the project provides the `chaosmonkey` command-line tool, which you can install this way:

```bash
go get github.com/mlafeldt/chaosmonkey/cmd/chaosmonkey
```

Use the tool to:

* Trigger a new chaos event:

    ```bash
    chaosmonkey -endpoint http://chaosmonkey.example.com:8080 \
        -strategy ShutdownInstance -group ExampleAutoScalingGroup
    ```

* Get a list of past chaos events:

    ```bash
    chaosmonkey -endpoint http://chaosmonkey.example.com:8080
    ```

* List all auto scaling groups for a given AWS account, which you may then pass to `-group`:

    ```bash
    export AWS_ACCESS_KEY_ID...
    export AWS_SECRET_ACCESS_KEY=...
    export AWS_REGION=...
    chaosmonkey -list-groups
    ```

* List available chaos strategies, which you may pass to `-strategy`:

    ```bash
    chaosmonkey -list-strategies
    ```

As always, invoke `chaosmonkey -h` for a list of all available options.

## Use with Docker

[This Docker image](https://github.com/mlafeldt/docker-simianarmy) allows you to deploy Chaos Monkey with a single command:

```bash
docker run -it --rm -p 8080:8080 \
    -e SIMIANARMY_CLIENT_AWS_ACCOUNTKEY=$AWS_ACCESS_KEY_ID \
    -e SIMIANARMY_CLIENT_AWS_SECRETKEY=$AWS_SECRET_ACCESS_KEY \
    -e SIMIANARMY_CLIENT_AWS_REGION=$AWS_REGION \
	-e SIMIANARMY_CHAOS_LEASHED=false \
	-e SIMIANARMY_CHAOS_TERMINATEONDEMAND_ENABLED=true \
    mlafeldt/simianarmy
```

Afterwards, you can use `chaosmonkey` to talk to the dockerized Chaos Monkey:

```bash
chaosmonkey -endpoint http://$DOCKER_HOST_IP:8080 ...
```

## Author

This project is being developed by [Mathias Lafeldt](https://twitter.com/mlafeldt).

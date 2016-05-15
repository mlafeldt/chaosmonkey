# chaosmonkey

[![Build Status](https://travis-ci.org/mlafeldt/chaosmonkey.svg?branch=master)](https://travis-ci.org/mlafeldt/chaosmonkey)
[![GoDoc](https://godoc.org/github.com/mlafeldt/chaosmonkey/chaosmonkey?status.svg)](https://godoc.org/github.com/mlafeldt/chaosmonkey/chaosmonkey)

Go client to the [Chaos Monkey REST API](https://github.com/Netflix/SimianArmy/wiki/REST) that can be used to trigger and retrieve chaos events.

I started the project for the purpose of controlled failure-injection during GameDay events.

## Prerequisites

In order to trigger chaos events, Chaos Monkey must be unleashed and on-demand termination must be enabled via these configuration properties:

```
simianarmy.chaos.leashed = false
simianarmy.chaos.terminateOndemand.enabled = true
```

## Go library

Run this command to install the `chaosmonkey` library from source:

```
$ go get github.com/mlafeldt/chaosmonkey/chaosmonkey
```

For usage and examples, see the [Godoc documentation](https://godoc.org/github.com/mlafeldt/chaosmonkey/chaosmonkey).

## CLI

In addition to the library, the project provides the `chaosmonkey` command-line tool, which you can install this way:

```
$ go get github.com/mlafeldt/chaosmonkey/cmd/chaosmonkey
```

Use the tool to trigger a new chaos event:

```
$ chaosmonkey -endpoint http://chaosmonkey.example.com:8080 \
    -strategy ShutdownInstance -group ExampleAutoScalingGroup
```

Or to get a list of past chaos events:

```
$ chaosmonkey -endpoint http://chaosmonkey.example.com:8080
```

Run `chaosmonkey -h` for a list of all available options.

## Use with Docker

[This Docker image](https://github.com/mlafeldt/docker-simianarmy) allows you to deploy Chaos Monkey using a single command:

```
docker run -it --rm -p 8080:8080 \
    -e SIMIANARMY_CLIENT_AWS_ACCOUNTKEY=$AWS_ACCESS_KEY_ID \
    -e SIMIANARMY_CLIENT_AWS_SECRETKEY=$AWS_SECRET_ACCESS_KEY \
    -e SIMIANARMY_CLIENT_AWS_REGION=$AWS_REGION \
	-e SIMIANARMY_CHAOS_LEASHED=false \
	-e SIMIANARMY_CHAOS_TERMINATEONDEMAND_ENABLED=true \
    mlafeldt/simianarmy
```

You can then use `chaosmonkey` to talk to the dockerized Chaos Monkey:

```
chaosmonkey -endpoint http://$DOCKER_HOST_IP:8080 ...
```

## Author

This project is being developed by [Mathias Lafeldt](https://twitter.com/mlafeldt).

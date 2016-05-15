# chaosmonkey

[![Build Status](https://travis-ci.org/mlafeldt/chaosmonkey.svg?branch=master)](https://travis-ci.org/mlafeldt/chaosmonkey)
[![GoDoc](https://godoc.org/github.com/mlafeldt/chaosmonkey/chaosmonkey?status.svg)](https://godoc.org/github.com/mlafeldt/chaosmonkey/chaosmonkey)

Go client to the [Chaos Monkey REST API](https://github.com/Netflix/SimianArmy/wiki/REST) that can be used to trigger and retrieve chaos events.

I started the project for the purpose of controlled failure-injection during GameDay events (in combination with [this Docker image](https://github.com/mlafeldt/docker-simianarmy)).

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

## Prerequisites

Note that in order to trigger chaos events, Chaos Monkey must be unleashed and on-demand termination must be enabled via these configuration properties:

```
simianarmy.chaos.leashed = false
simianarmy.chaos.terminateOndemand.enabled = true
```

## Author

This project is being developed by [Mathias Lafeldt](https://twitter.com/mlafeldt).

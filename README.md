# chaosmonkey

[![Build Status](https://travis-ci.org/mlafeldt/chaosmonkey.svg?branch=master)](https://travis-ci.org/mlafeldt/chaosmonkey)
[![GoDoc](https://godoc.org/github.com/mlafeldt/chaosmonkey/lib?status.svg)](https://godoc.org/github.com/mlafeldt/chaosmonkey/lib)

Go client to the [Chaos Monkey REST API](https://github.com/Netflix/SimianArmy/wiki/REST) that can be used to trigger and retrieve chaos events.

This project was started for the purpose of controlled failure injection during GameDay events.

## Prerequisites

First of all, you need a running [Simian Army](https://github.com/Netflix/SimianArmy) that exposes its REST API via HTTP.

In order to trigger chaos events via the API, Chaos Monkey must be unleashed and on-demand termination must be enabled via these configuration properties:

```
simianarmy.chaos.leashed = false
simianarmy.chaos.terminateOndemand.enabled = true
```

## CLI

### Installation

If you're on Mac OS X, the easiest way to get the `chaosmonkey` command-line tool is via Homebrew:

```bash
brew tap mlafeldt/formulas
brew install chaosmonkey
```

You can also build the tool from source, provided you have Go installed:

```bash
go get -u github.com/mlafeldt/chaosmonkey
```

### Usage

Use the tool to:

* Trigger a new chaos event:

    ```bash
    chaosmonkey -endpoint http://example.com:8080 \
        -group ExampleAutoScalingGroup -strategy ShutdownInstance
    ```

* Trigger the same event 5 times at intervals of 10 seconds with a probability of 20%:

    ```bash
    chaosmonkey -endpoint http://example.com:8080 \
        -group ExampleAutoScalingGroup -strategy ShutdownInstance \
        -count 5 -interval 10s -probability 0.2
    ```

    This is useful to terminate more than one EC2 instance of an auto scaling group.

* Get a list of past chaos events:

    ```bash
    chaosmonkey -endpoint http://example.com:8080
    ```

* List available chaos strategies, which you may pass to `-strategy`:

    ```bash
    chaosmonkey -list-strategies
    ```

* List all auto scaling groups for a given AWS account, which you may then pass to `-group`:

    ```bash
    export AWS_ACCESS_KEY_ID=...
    export AWS_SECRET_ACCESS_KEY=...
    export AWS_REGION=...
    chaosmonkey -list-groups
    ```

* Wipe state of Chaos Monkey by deleting its SimpleDB domain (named `SIMIAN_ARMY` by default):

    ```bash
    export AWS_ACCESS_KEY_ID=...
    export AWS_SECRET_ACCESS_KEY=...
    export AWS_REGION=...
    chaosmonkey -wipe-state SIMIAN_ARMY
    ```

    Warning: Requires a restart of Chaos Monkey.

As always, invoke `chaosmonkey -h` for a list of all available options.

In addition to command-line options, the tool also understands these environment variables:

* `CHAOSMONKEY_ENDPOINT` - the same as `-endpoint`
* `CHAOSMONKEY_USERNAME` - the same as `-username`
* `CHAOSMONKEY_PASSWORD` - the same as `-password`

### Use with Docker

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

## Go library

In addition to the CLI tool, the project also provides the `chaosmonkey` Go library for use in other Go projects. To install it from source:

```bash
go get -u github.com/mlafeldt/chaosmonkey/lib
```

For usage and examples, see the [Godoc documentation](https://godoc.org/github.com/mlafeldt/chaosmonkey/lib).

## Further resources

* [Talk: Chaos Monkey for Fun and Profit](https://mlafeldt.github.io/blog/chaos-monkey-for-fun-and-profit/)
* [Article: A Little Story about Amazon ECS, systemd, and Chaos Monkey](https://medium.com/production-ready/a-little-story-about-amazon-ecs-systemd-and-chaos-monkey-8bc7d1b5778)

## Author

This project is being developed by [Mathias Lafeldt](https://twitter.com/mlafeldt).

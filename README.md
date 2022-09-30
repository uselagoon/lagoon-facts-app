# amazeeio-facts-audit

This application is run as part of the bootstrapping process of the Lagoon CLI pods.
It is used to gather facts about the running environment and push them back up to the Lagoon Insights system.

Each fact is determined by a "gatherer", essentially an object that implements the `Gatherer` interface and registers itself with the application.

We run through each fact Gatherer, which determines whether it is able to be run in the current environment, and, if so, it gathers its facts. These are then written back to Lagoon.

## Table of Contents

* [Setup](#setup)
* [Extending](#extending)
    * [Adding a new gatherer](#adding-a-new-gatherer)
* [Running](#running)
* [Arbitrary Facts](#arbitrary-facts)
* [Building](#building)

## Setup

The application is written in Go (GoLang), for more information about setting this up please see [https://golang.org/](https://golang.org/).

## Extending

### Adding a new gatherer

As mentioned above, creating a new gatherer is fairly simple.

Look at the requirements of the [`Gatherer` interface](https://github.com/bomoko/amazeeio-facts-audit/blob/main/gatherers/defs.go#L14) and implement the `AppliesToEnvironment` and `GatherFacts` functions on a new structure.
Then register the new gatherer by calling `RegisterGatherer` (for example, [here](https://github.com/bomoko/amazeeio-facts-audit/blob/main/gatherers/DrushGatherer.go#L56))

## Running

To run the application you can use the following:

```bash
$ go run main.go [command] 
$ go run main.go gather # run the gather command
```

However, most normal circumstances this will fail to run, as the gather will try to look for environment variables.

You can preface the command with these environment variables, this will allow it to continue.

```bash
$ LAGOON_PROJECT=[project_name] LAGOON_GIT_BRANCH=[git_branch] go run main.go gather
```

Although you will be doing a scan of your local machine, and will send this information as facts up to which ever project has been set.

**Not recommended** but might be required for some testing.

The recommended way to use this application is from inside a lagoon container, assuming this container does not have go installed you will need to [build](#building) the application.

## Arbitrary Facts

If you would like to generate facts outside of the provided fact gatherers, you can write them to a file, and then have the `file gatherer` parse and write the facts to the backend.

The structure of the file should follow closely the definition of a `GatheredFact` in `/gathers/defs.go`.

The gatherer expects the incoming file to contain an array of facts, of the following format

```
[
    {
        name: "fact name",
        value: "value stored against name",
        source: "Where this fact is sourced from",
        environment: <environment ID>,
        description: "description of fact",
        keyFact: true|false,
        category: "See /gatheres/categories.go for values"
    }
]

```

## Building

To build the project for use inside a container you will need a statically linked binary.

**MacOS**
```bash
$ env CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 go build -a -o builds/lagoon-facts-macos
```

**Linux**
```bash
$ env CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -o builds/lagoon-facts-linux
```

Once built from inside a container you can use the application:

```bash
$ ./builds/lagoon-facts-linux [command]
$ ./builds/lagoon-facts-linux gather # run the gather command
```


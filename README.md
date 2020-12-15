# amazeeio-facts-audit

This application is run as part of the bootstrapping process of the Lagoon CLI pods.
It is used to gather facts about the running environment and push them back up to the Lagoon Insights system.

Each fact is determined by a "gatherer", essentially an object that implements the `Gatherer` interface and registers itself with the application.

We run through each fact Gatherer, which determines whether it is able to be run in the current environment, and, if so, it gathers its facts. These are then written back to Lagoon.

### Extending - adding a new gatherer

As mentioned above, creating a new gatherer is fairly simple.

Look at the requirements of the (`Gatherer` interface)[https://github.com/bomoko/amazeeio-facts-audit/blob/main/gatherers/defs.go#L14] and implement the `AppliesToEnvironment` and `GatherFacts` functions on a new structure.
Then register the new gatherer by calling `RegisterGatherer` (for example, (here)[https://github.com/bomoko/amazeeio-facts-audit/blob/main/gatherers/DrushGatherer.go#L56])

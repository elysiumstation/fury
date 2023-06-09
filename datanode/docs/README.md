# Fury core architecture

Data node is a stand alone product that is built on the top of Fury core product.
It consumes stream of events from core Fury via socket using [Broker](./broker.md) then aggregates the events and save them to storage.

## Component relationships

The following diagram shows how the various components of this implementation interact with each other at a high level.

![Fury core protocol architecture](diagrams/design-architecture-2023-01-26.svg "Fury core protocol architecture")

## Modelling the domain

Some subdirectories contain Golang packages which represent a discrete domain or concept from the whitepaper.
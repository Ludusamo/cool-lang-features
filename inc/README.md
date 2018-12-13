# Internode Communication

## Purpose

Library for handling the management of communication between an arbitrary set of nodes.

## Usage

A `ConnectionManager` needs to be instantiated on each node that wants to communicate.

A node can then ask the `ConnectionManager` for a channel to write messages to a particular TCP address.

The `ConnectionManager` will then return a pair of `chan interface{}` to which a node or goroutine can then send objects to be JSON serialized and sent.

For every operation, the second of the pair of channels will be messages sent to the process.

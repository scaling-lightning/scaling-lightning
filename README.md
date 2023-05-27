# Scaling Lightning - A stress-testing tool-kit for the Lightning Network

This is an initiative to build a stress-testing tool-kit for the Lightning Network protocol, implementations
and applications that depends on the Lightning Network.

The goal is to collaborate as an industry to help scale the Lightning Network and the applications that depend on it.

## Why is this important?

Currently, there are unknowns and untested assumptions about how the Lightning Network and applications will 
react to shocks in transaction volume, channels, nodes, gossip messages, etc.

Having a set of tools and a signet Lightning network will help:

* Developers to test their applications.
* Researchers to test their assumptions.
* Operators to test their infrastructure.
* Novices to learn how the Lightning Network and various applications works in a somewhat realistic environment 
  without risking real coins.

## How will it work?

We are still in the early stages of planning, but the first tool we are building will be a tool to quickly generate 
one or more Lightning Nodes that can connect either to a public Signet Lightning Network or 
a private Regtest Network Lightning for any combination of LN implementations (CLN, LND, LDK, Acinq etc.).

Other tools made specifically for testing isolated parts of the protocol is also relevant. This can help developers 
and researchers test their assumptions in an isolated environment. 
An example of this is [The Million Channels Project](https://github.com/rustyrussell/million-channels-project-data) 
developed by Rusty Russell to test gossip.

## How can you help?

* We need developers to contribute to the tool-kit.
* Operators to run nodes on the Signet Lightning Network.
* Researchers to help design general tools relevant to them.

![ScalingLN Twitter Banner](https://github.com/ohenrik/scaling-lightning/assets/647617/8511c586-7549-4e2b-ad6d-bf87419a624c)

# Scaling Lightning - A Testing Toolkit for the Lightning Network

This initiative aims to build a testing toolkit for the Lightning Network protocol, its implementations, and 
applications that depend on the Lightning Network.

The goal is to collaborate as an industry to help scale the Lightning Network and the applications that depend on it.

## Why is this important?

Currently, there are unknowns and untested assumptions about how the Lightning Network and its applications will react 
to shocks in transaction volume, channels, nodes, gossip messages, etc.

Having a set of tools and a signet Lightning Network will help:

* Developers test their applications.
* Researchers verify their assumptions.
* Operators test their infrastructure.
* Novices learn how the Lightning Network and various applications work in a somewhat realistic environment without 
  risking real coins.

## How will it work?

We are still in the early stages of planning, but the first tool we are building will be a tool to quickly generate one 
or more Lightning Nodes. These nodes can connect either to a public signet Lightning Network or a private Regtest 
Lightning Network for any combination of LN implementations (CLN, LND, LDK, Acinq etc.).

Other tools, made specifically for testing isolated parts of the protocol, are also relevant. These can help developers
and researchers test their assumptions in an isolated environment. An example of this is 
[The Million Channels Project](https://github.com/rustyrussell/million-channels-project-data) developed by Rusty Russell 
to test gossip.

## Project milestones

This is an outline of the project's milestones. We will further detail these milestones using the features of GitHub's milestones and project management tools:

* [ ] Create a Helm chart to initiate a Regtest Lightning Network, making it easier for developers to start development.
* [ ] Create a Helm chart for one or more Signet Bitcoin and Lightning Network nodes.
* [ ] Develop a tool for generating or simulating channel openings/closings and traffic between nodes.
* [ ] Develop a tool to simplify interaction with your nodes in your signet or regtest network.
* [ ] Thoroughly document the project and provide instructions for use.
* [ ] Create a website to host the documentation, using Docusaurus.
* [ ] Document and provide links to supporting development resources.

## How can you help?

* We need developers to contribute to the toolkit.
* We invite operators to run nodes on the Signet Lightning Network.
* We encourage researchers to help design general tools relevant to them.
* Donate to the project to help fund development and maintain the signet Lightning Network.

## How can you reach us?

Telgram: https://t.me/+AytRsS0QKH5mMzM8
Twitter: https://twitter.com/ScalingLN

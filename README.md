## Setup

- Install Encore
  https://encore.dev/docs/install
- Ensure Temporal Server is running, easiest way to setup is to use docker-compose from Temporal
  https://github.com/temporalio/docker-compose

## Start Server

- run `encore run dev`

## Temporal as Backend for low code

This repository contains code for a custom workflow engine built on top of Temporal.

It includes a YAML parser to define workflows, a topological sort algorithm to determine the order of execution, and activity functions to perform specific tasks.

The workflow engine executes each node in the workflow as an activity and passes the output of the previous node as input to the next node.

The code also includes a main function to initialize the workflow engine and start the workflow execution.

Custom DSL is defined in [workflow.yaml](dsl/workflow.yaml)

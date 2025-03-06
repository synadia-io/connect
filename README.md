# Connect
Connect is a NATS based data pipeline that allows you to easily connect data sources and sinks to NATS. Connect is built on 
top of the NATS ecosystem and uses the NATS server as a message broker. Connect is designed to be easy to use and 
scalable, allowing you to connect multiple data sources and sinks to NATS.

> [!IMPORTANT]
> This project merely provides the connect CLI and SDK to interact with the connect service. The connect service is 
> currently hosted by Synadia during a private beta, but we do want to let you all know what we are working on 
> and get feedback.
> If you have any questions, feel free to reach out in the #connectors channel on the NATS slack.

## Installation
Download the `connect` binary from the [releases page](https://github.com/synadia-io/connect/releases) and place them in your PATH.

Make sure you are using the right nats context
```shell
nats context select
```

To make sure your binary is correct and everything is up and running, try listing the available components:
```shell
connect library ls
```

Once you have the binary correctly installed, you may want to move on to create your first connector. Take a look at
the [getting started guide](docs/getting-started.md) for more information on how to accomplish that.

Another useful resource is the [connector specification](spec/schemas/connector-spec.schema.json) which details the 
structure of a connector.

Further reference documentation can be found in the  [Connector Reference documentation](docs/reference/connector.md)

## Usage
You can use the `connect` command to interact with the connect service and manage connectors or explore the library.

Take a look at the help page to see what you can do:
```shell
connect --help
```

## Contributing
We love to get feedback, bug reports, and contributions from our community. If you have any questions or want to
contribute, feel free to reach out in the #connectors channel on the NATS slack.

Take a look at the [Contributing](CONTRIBUTING.md) guide to get started.
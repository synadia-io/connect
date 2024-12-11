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

### Install as a plugin in NATS CLI
```shell
nats plugins register connect <path to your connect binary>
```

## Documentation
Documentation for connect is rather sparse for the time being. We are working on improving the documentation and adding
more examples. If you have any questions, feel free to ask in the #connector channel on the Synadia slack.

Reference:
- [Connector Reference documentation](docs/reference/connector.md)

## Usage
You can use the `connect` command to interact with the connect service and manage connectors or explore the library.

Take a look at the help page to see what you can do:
```shell
connect --help
```

## Getting Started
To get started with connect, take a look at the [Getting Started](docs/getting-started.md) guide.

## Contributing
We love to get feedback, bug reports, and contributions from our community. If you have any questions or want to
contribute, feel free to reach out in the #connectors channel on the NATS slack.

Take a look at the [Contributing](CONTRIBUTING.md) guide to get started.
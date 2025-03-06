## Getting Started
So off to the horses then! Let's get started creating your first inlet.

### Private Beta
Connect is currently in private beta. If you would like to participate in the beta, please reach out to us on the 
Synadia slack so we can enable the connect feature for you.

### Inlet vs Outlet
Connectors come in two distinct flavors: inlets and outlets. Inlets are connectors that read data from an external system
and write it to NATS. Outlets are connectors that read data from NATS and write it to an external system.

Inlets make use of a `source` to read data from an external system while outlets make use of a `sink` to write data to
an external system. To know which sources are available, you can use the `connect component search` command.

### Create a new inlet
Let's create a new inlet. In this case we will create an inlet that reads from a `generate` source and publishes those
generated messages to NATS.

`generate` is a source that generates messages at a specified interval. Not only is this useful for testing, but it can
also be used to generate messages at a specific moment in time, configured through a cron expression.

To create a new inlet, run the following command:
```shell
connect connector create hello
```
The `-i` in this case tells the create command we want to interactively create the connector. We will be asked which
type of connector we want to create. Select `inlet` and press enter. The CLI will generate a template inlet definition
and open your default editor.

### Configure the inlet
To makes things easier, we already have a template for you to use. Edit the template to look like this:
```yaml
description: "A simple hello world inlet"
runtime_id: wombat
steps:
  source:
    type: nats
    config:
      subject: "testing.hello.input"
      url: nats://demo.nats.io:4222
  producer:
    core:
      subject: "testing.hello.output"
    nats_config:
      url: nats://demo.nats.io:4222
```

Let's run through the different fields in this file:
- `description`: A short description of what the connector does
- `runtime_id`: The id of the runtime providing the building blocks for the connector 
- `steps`: The steps that make up the connector. In this case, we have two steps: `source` and `producer`
- `source.type`: The type of source. Take a look at the `connect component search` command to see the available sources.
- `source.config`: The configuration for the source. This depends on the type of source being used.
- `producer.core.subject`: The NATS subject on which the message will be published.
- `producer.nats_config.url`: The URL of the NATS server we want to publish the data to.

Save the inlet and exit the editor. You can run `connect connectors list` (or `connect c ls`) to see the newly created inlet.

### Start the inlet
You can now deploy the connector by running the following command:
```shell
connect connector start hello
```

You should now see messages being sent to the `testing.hello.*` subject on the NATS server:
```shell
nats -s "nats://demo.nats.io:4222" sub 'testing.hello.>'
```

### Cleaning up
Use `connect connector stop hello` to stop the connector. You can then delete the connector using the 
`connect connector delete hello` command.
# Mapping Transformer
A mapping transformer allows you to change the structure and contents of a message by applying a mapping to the message.

## Example
```yaml
mapping:
  sourcecode: |-
    root = this.content().uppercase()
```

In this example, the `mapping` transformer is used to uppercase the content of the message.

## Fields
| Field        | Type   | Default | Required | Description                 |
|--------------|--------|---------|----------|-----------------------------|
| `sourcecode` | string |         | yes      | The source code to execute. |
```
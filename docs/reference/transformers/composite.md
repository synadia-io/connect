# Composite Transformer
The composite transformer allows you to chain multiple transformers together. This is useful when you need to apply 
multiple transformations to a message. Each transformer will be executed in sequence with the output of the previous
transformer being passed as the input to the next.

## Example
```yaml
composite:
  sequential:
    - mapping:
        sourcecode: |-
          root = this.content().uppercase()
    - mapping:
        sourcecode: |-
          root = this.content().replace("FOO", "BAR")
```

In this example, the `composite` transformer is used to chain two `mapping` transformers together. The first transformer
will uppercase the content of the message, and the second transformer will replace the string `FOO` with `BAR`.

## Fields
| Field        | Type                                      | Default | Required | Description                                      |
|--------------|-------------------------------------------|---------|----------|--------------------------------------------------|
| `sequential` | array of [Transformer](../transformer.md) |         | yes      | An array of transformers to execute in sequence. |

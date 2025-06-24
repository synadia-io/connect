# Source
A source describes how to connect to an external system and read data from it.

## Example
```yaml
type: "aws_sqs"
config:
  url: "https://sqs.us-east-2.amazonaws.com/563342913055/connect-test-1"
  region: "us-east-2"
  credentials:
    id: "YOUR_AWS_ACCESS_KEY_ID"
    secret: "YOUR_AWS_SECRET_ACCESS_KEY"
```

# Fields
## `type`
The type of the source to use. The name can be found by querying the library using
`connect library list --kind=source`

## `config`
The configuration for the source. The fields in this section are specific to the source type. Find
the fields by querying the library using `connect library get <runtime> source <type>`
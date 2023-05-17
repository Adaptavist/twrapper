# Integration Tests

Because integrations tests require access to remote systems, we have used build tags to isolate
these tests to those working explicitly on these integrtions.

Please add `integration_test` to your __Go Build Tags__ in your prefered IDE to be able to debug the
tests.

## Workstation setup

To make setup as simple as posssible, we've opted to load global configuration file from your
__$HOME__ directory. This avoid leaking information about our and your organization by accidentally
hardcoding them into tests.

The schema is identical to __terraform.twrapper.yml__ we've just segmented the configuration for
testing purposes.

## Running integration tests

```bash
# run integration tests
make integration_tests
```

### Configuration

Do as much configuring in your environment as possible, however, do not set `TW_WORKSPACE_ID` and `TW_BACKEND_KEY` as
the integration tests will generate values for these. You need to make sure the backend, bucket, region, table are set.
And the cloud organisation is set.

```yaml
# $HOME/test.twrapper.yml
backend:
  type: s3
  props:
    bucket: ${TW_BACKEND_PROPS_BUCKET}
    region: ${TW_BACKEND_PROPS_REGION}
    dynamodb_table: ${TW_BACKEND_PROPS_DYNAMODB_TABLE}
    encrypted: true
    key: ${TW_BACKEND_KEY}.tfstate
cloud:
  organization: ${TW_CLOUD_ORGANIZATION}
  project: twrapper-integration-tests
  workspace: twrapper-integration-test-${TW_WORKSPACE_ID}
```
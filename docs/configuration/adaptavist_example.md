# Adaptavist Example

This is roughly how we implement Twrapper in our self-service repos

## Backend configuration

This is representative of what our backend configuration will look like after this release.

```yaml
# test twrapper doc
backend:
  type: s3
  props:
    bucket: ${TW_BACKEND_BUCKET}
    dynamodb_table: ${TW_BACKEND_TABLE}
    encrypt: true
    region: ${TW_BACKEND_REGION}
    key: ${TW_BACKEND_KEY}.tfstate
```

## Cloud Configuration

This is what we are moving towards, and will likely need to overload with the backend config to migrate the state.

```yaml
cloud:
  organization: ${TW_CLOUD_ORG}
  project: ${PROJECT}
  # workspaces must be globally unique in the organization, so we build something unique enough to avoid collisions
  # we even encoure uses to use a section of a v4 UUID as the TW_WORKSPACE_ID component
  workspace: ${TEAM}-${PROJECT}-${MODULE}-${SUFFIX}-${TW_WORKSPACE_ID}
```

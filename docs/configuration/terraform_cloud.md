# Terraform Cloud

This is now our proference over the [AWS Backend](aws.md) configuration. So much so that if you configure both a
backend and cloud. Twrapper will migrate your state to Terraform Cloud and remove it from the backend.

We support the same configuration as Terraform. E.G using `~/.terraform.d/credentials.json` or setting `TF_TOKEN_app_terraform_io` in the environment. Avoiding another method of configuring for the same service.

## Configuration

Again we prefer to pass variables in via the environment but clearly define the config.

```yaml
# terraform.twrapper.yaml
cloud:
    organization: ${ORG}
    project: ${PROJECT}
    workspace: ${WORKSPACE}
```

### Project and workspace creation

We will automatically create project sand workspaces before running any Terraform operations. This is to help keep
organizations tidy by assigning workspaces to projects early in the lifecycle.

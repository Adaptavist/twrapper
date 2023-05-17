# Configurations

## Table of Contents

- [Configurations](#configurations)
  - [Table of Contents](#table-of-contents)
  - [WARNING: Backend and Cloud configuration presedence](#warning-backend-and-cloud-configuration-presedence)
  - [Where can configuration be stored](#where-can-configuration-be-stored)
  - [Using the environment for configuration](#using-the-environment-for-configuration)
  - [Environment variables](#environment-variables)

Please also checkout specific configuration for [AWS](aws.md), [Terraform Cloud](terraform_cloud.md), and
[Adaptavists Example](adaptavist_example.md).

## WARNING: Backend and Cloud configuration presedence

While partial backend and cloud configurations can exist for the backend and cloud actions to work, both
configurations cannot be used when running the terraform cli wrapping commands. This is because Terraform needs
only one backend configuration which `cloud` overlaps.

**If both configurations are found, Twrapper will automatically migrate from your local/s3 backend to cloud.**

## Where can configuration be stored

Twrapper will look for `terraform.twrapper.yml` in the following locations:

- In your current working directory (ideal when working within a project)
- In your home directory (idea when working accross projects

E.G terraform/workspaces/my-workspace/twrapper.terraform.yml or ~/twrapper.terraform.yml

## Using the environment for configuration

There are two methods of utilising the environment for configuration. Once is explicitly setting a variable in the
environment that matches the structure of configuration with a `TW_` prefix, which is spoken more about in the 
[Environment variables section](#environment-variables).

Or there is the substitution method, which allows you to lookup environment variables from the configuration file
itself, using `${VAR_NAME}`. Be warned Twrapper will panic if the variable is not found.

For example:

```yaml
aws:
    account_id: ${AWS_ACCOUNT_ID}
```

## Environment variables

Environment variables you can use to override any settings, these will take precedence over the configuration file.

| Variable name             | Description                                                                                                      |
| ------------------------- | ---------------------------------------------------------------------------------------------------------------- |
| TW_ENV                    | Useful for debugging, set to development for increased verbosity                                                 |
| TW_AWS_ACCOUNT_ID         | Environment var account ID is available                                                                          |
| TW_AWS_TRY_ROLE_NAMES     | Role names to enumerate, needs account ID. comma delimited. E.G `role-1,role-2,role-3`                           |
| TW_AWS_ROLE_ARN           | Hard-code role ARN to assume                                                                                     |
| TW_AWS_ROLE_TF_VAR        | Terraform variable to pass a role arn for assuming, requires `role_arn` or `account_id_var` and `try_role_names` |
| TW_BACKEND_TYPE           | Backend you want twrapper to configure, s3 or local are supported                                                |
| TW_BACKEND_PROPS_BUCKET   | `backend_type=s3` S3 Bucket of the backend                                                                       |
| TW_BACKEND_PROPS_REGION   | `backend_type=s3` Region of the backend                                                                          |
| TW_BACKEND_PROPS_ENCYTPED | `backend_type=s3` Is the backend encrypted?                                                                      |
| TW_BACKEND_PROPS_KEY      | `backend_type=s3` Object key                                                                                     |
| TW_BACKEND_PROPS_PATH     | `backend_type=local` Statete path                                                                                |
| TW_CLOUD_ORGANIZATION     | Terraform enterprise org                                                                                         |
| TW_CLOUD_WORKSPACE        | Terraform enterprise workspace                                                                                   |
| TW_REQUIRED_VARS          | Variables that are required by the module. E.G `TF_VAR_var_name, TF_VAR_name`                                    |

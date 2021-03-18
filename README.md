# Twrapper Terraform Wrapper

## About

This is not an alternative to or a replacement for Terragrunt!

At adaptavist we have numerous private GIT repos with Terraform modules ready to be remotely 
triggered enabling some level of self-service. However, we need to wrap terraform to automate 
backend configuration and ensuring a suitable role can be assumed. We already had something
written in python but we want to replace it with something lighter weight and pre-bundled in a 
container.

Our intentions with this wrapper

- Have configuration file per module
- Check for environment variable existance as we often use them for variable assignment
- Restrict a module to a role
- Allow searching through a list of roles to find a suitable one to run the module
- Ensure Account IDs are specified for the roles we use when working accross accounts
- Force a unique backend key by using UUIDv4
- Enable running most Terraform commands

## Configuration

```yaml
# module/twrapper.yaml
# error if these vars are not set or are empty
required_vars:
  - TF_VAR_example_var_name
# Tells Twrapper we are working on AWS
aws:
  # If you want to provide a role to terraform, you must tell it the variable used to apply it.
  role_tf_var: aws_assume_role_arn
  # role_arn pins your module to a specific role
  role_arn: arn:aws:iam::00000000:role/AdminisratorAccess
  # try these roles if role name is not set, only used if AccountIDVar is set
  try_role_names: [ OrganizationAccountAccessRole ]
  # when using try_role_names, ACCOUNT_ID must be present in the env
  account_id_var: ACCOUNT_ID
# Tell Twrapper what environment variable should be used for the key_id
backend_key_id_var: KEY_ID # VAR used to get the key id
# Also support local for development, may add support for more in time
backend_type: s3
# optional - defaults to {key_id}.tfstate
backend_key_template: "module-{key_id}.tfstate"
backend_config:
  bucket: "backend-bucket-name"
  region: "us-east-1"
  dynamodb_table: "backend-table-name"
  encrypt: true
```

### Example Configs

#### Required Terraform Variables

You may want Terraform to validate your variables but you may be using a third-party module. So Twrapper
provides a real basic feature for checking variables are set and not empty, but it will not validate
the value beyond that.

```yaml
required_vars:
  - TF_VAR_my_var_1
  - TF_VAR_my_var_2
```

#### Backend configuration

```yaml
# What type of backend
backend_type: s3
# The environment variable which has the backend key
backend_key_id_var: KEY_ID
# override the key/path format for your backend
backend_key_template: "module-{key_id}.tfstate"
# Main backend configure, this is generally static
backend_config:
  bucket: "backend-bucket-name"
  region: "us-east-1"
  dynamodb_table: "backend-table-name"
  encrypt: true
```

#### AWS Configuration

Chances are you're deploying to some sort of IaaS platform. and you need an appropriate session. We
aim to resolve this by trying to resolve which roles can be used. We understand that as busineses
grow the role names may not be consistent and therefore the ability to try others is helpful.

```yaml
# Tell the AWS package to try and assume one of the two roles and assign the ARN of the working
# role to terraform as "assume_role_arn"
aws:
  role_tf_var: assume_role_arn
  account_id_var: AWS_ACCOUNT_ID
  try_role_names:
    - MyAdminRole
    - MyOrganizationRole

# Fix your module to a single role on a single account, and assign to the "assume_role_arn" variable
aws:
  role_tf_var: assume_role_arn
  role_arn: arn:aws:iam::00000000:role/RoleName
```


## Usage

```bash
twrapper [terraform args...]
twrapper plan # run a terraform plan
twrapper apply # runs an apply but adds -auto-approve
twrapper destroy # runs a destroy but add -force
```
# Home

- [Home](#home)
  - [Developer setup](#developer-setup)
    - [Software dependencies](#software-dependencies)
    - [Accounts, and APIs](#accounts-and-apis)
      - [AWS](#aws)
      - [Terraform Cloud Platform](#terraform-cloud-platform)
  - [Testing](#testing)
    - [Backend tests](#backend-tests)

## Developer setup

### Software dependencies

- Golang 1.20+
- Terrform 1.4+

### Accounts, and APIs

If you're working with backends and APIs you may need write or administrative access to the
following platforms:

#### AWS

Twrapper supports utilises the AWS SDK, so you can make use of AWS profiles and SSO.

```bash
# Set the profile
export AWS_PROFILE=my-profile
# Create a session with SSO
aws sso login
# Get to work :smile:
twrapper ...
```

#### Terraform Cloud Platform

We also support Terraform Cloud Platform but only HashiCorps SaaS offering. We support the same configuration as
Terraform. E.G using `~/.terraform.d/credentials.json` or setting `TF_TOKEN_app_terraform_io` in the environment.

## Testing

### Backend tests

To test certain backend operations that twrapper automates, you should checkout the
[Backends](backends/index.md) section.

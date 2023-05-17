# Twrapper Terraform Wrapper

## NOTICE

We will be deprecating, this project very soon. Once we've migrated to Terraform Cloud, we'll no longer maintain this
project.

## About

This is not an alternative to or a replacement for Terragrunt!

At adaptavist we have numerous private GIT repos with Terraform modules ready to be remotely triggered enabling some
level of self-service. However, we need to wrap terraform to automate backend configuration and ensuring a suitable role
can be assumed. We already had something written in python but we want to replace it with something lighter weight and
pre-bundled in a container.

Our intentions with this wrapper

- Have configuration file per module
- Check for environment variable existance as we often use them for variable assignment
- Restrict a module to a role
- Allow searching through a list of roles to find a suitable one to run the module
- Ensure Account IDs are specified for the roles we use when working accross accounts
- Enable running most Terraform commands
- Migrating to Terraform Cloud

## Usage

Please check out the [documentation](docs/) for usage.
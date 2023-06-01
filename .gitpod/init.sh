#!/usr/bin/env sh

# AWS

[[ -d /home/gitpod/.aws ]] || mkdir /home/gitpod/.aws
cat <<- AWSFILE > /home/gitpod/.aws/config
[default]
cli_pager=
output=json

# profile pinned in environment (~/.bashrc.d/300-aws)
[profile main]
sso_start_url=${AWS_SSO_START_URL}
sso_region=${AWS_SSO_REGION}
sso_account_id=${AWS_SSO_ACCOUNT_ID}
sso_role_name=${AWS_SSO_ROLE_NAME}
region=${AWS_REGION}
AWSFILE

echo "Generated: /home/gitpod/.aws/config"

# TERRAFORM

[[ -d /home/gitpod/.terraform.d ]] || mkdir /home/gitpod/.terraform.d
cat <<- TFRCFILE > ~/.terraform.d/credentials.tfrc.json
{
  "credentials": {
    "app.terraform.io": {
      "token": ${TF_TOKEN_app_terraform_io}
    }
  }
}
TFRCFILE

echo "Generated: /home/gitpod/.terraform.d/credentials.tfrc.json"

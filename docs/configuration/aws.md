# AWS Configuration

TL/DR:

- Twrapper uses the AWS SDK, so it will pickup your AWS configuration exactly the same way as the AWS CLI.
- Twrapper needs access to the backend, if you're using an S3 Backend with DynamoDB.
- Twrapper needs permissions to assume roles, if you're planning on assuming roles

## Assuming roles

This example will try to assume each role in the account specified in the environment variable __AWS_ACCOUNT_ID__, on
the first possitive result, it will pass the ARN to Terraform in the __role_arn__ variable.

```yaml
# terraform.twrapper.yaml
aws:
    account_id: ${AWS_ACCOUNT_ID}
    try_role_names: 
    - role-1
    - role-2
    - role-3
    role_tf_var: role_arn
```

This next example, will assume a fixed role, but with a dynamic account number.

```yaml
# terraform.twrapper.yaml
aws:
    role_arn: arn:aws:${AWS_ACCOUNT_ID}:iam:role/myRole
    role_tf_var: role_arn
```

### Backend

When creating Terraform modules for self-service, it is rather easy to hard-code most of the backend configuration.
However, we prefer to keep the configuration within environment variables of our CI/CD platforms. Then look them up in
Twrappers configuration file.

```yaml
backend:
    type: s3
    props:
        # using variables shared from the CI/CD tooling
        bucket: ${BACKEND_BUCKET}
        dynamodb_table: ${BACKEND_DYNAMODB_TABLE}
        encrypted: ${BACKEND_ENCRYPTED}
        # key is actually provided by the user, rather than hard-coded in the environment
        key: ${BACKEND_KEY}
```

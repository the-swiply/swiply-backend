#!/bin/bash

set -eou pipefail

LOCKBOX_AWS_SECRET_ID=e6qfkbv38cqf78chgphl

AWS_ACCESS_KEY_ID=$(yc lockbox payload get --id=$LOCKBOX_AWS_SECRET_ID --key=aws_access_key_id --profile=swiply)
export AWS_ACCESS_KEY_ID

AWS_SECRET_ACCESS_KEY=$(yc lockbox payload get --id=$LOCKBOX_AWS_SECRET_ID --key=aws_secret_access_key --profile=swiply)
export AWS_SECRET_ACCESS_KEY

# Postgres users login and password section.
LOCKBOX_POSTGRES_ADMIN_SECRET_ID=e6qgabf1uqm0v7hnu2ht
export POSTGRES_ADMIN_LOGIN=$(yc lockbox payload get --id=$LOCKBOX_POSTGRES_ADMIN_SECRET_ID --key=login --profile=swiply)
export POSTGRES_ADMIN_PASSWORD=$(yc lockbox payload get --id=$LOCKBOX_POSTGRES_ADMIN_SECRET_ID --key=password --profile=swiply)

# Initialize terraform variables.
cat templates/template.tfvars | envsubst > swiply.tfvars

export YC_TOKEN=$(yc iam create-token --profile=swiply)
terraform "$@" -var-file=swiply.tfvars

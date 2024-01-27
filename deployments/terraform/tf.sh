#!/bin/bash

set -eou pipefail

LOCKBOX_AWS_SECRET_ID=e6qfkbv38cqf78chgphl

AWS_ACCESS_KEY_ID=$(yc lockbox payload get --id=$LOCKBOX_AWS_SECRET_ID --key=aws_access_key_id --profile=swiply)
export AWS_ACCESS_KEY_ID

AWS_SECRET_ACCESS_KEY=$(yc lockbox payload get --id=$LOCKBOX_AWS_SECRET_ID --key=aws_secret_access_key --profile=swiply)
export AWS_SECRET_ACCESS_KEY

export YC_TOKEN=$(yc iam create-token --profile=swiply)
terraform "$@"

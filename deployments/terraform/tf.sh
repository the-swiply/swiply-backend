#!/bin/bash

set -eou pipefail

LOCKBOX_AWS_SECRET_ID=e6qfkbv38cqf78chgphl

AWS_ACCESS_KEY_ID=$(yc lockbox payload get --id=$LOCKBOX_AWS_SECRET_ID --key=aws_access_key_id --profile=swiply)
export AWS_ACCESS_KEY_ID

AWS_SECRET_ACCESS_KEY=$(yc lockbox payload get --id=$LOCKBOX_AWS_SECRET_ID --key=aws_secret_access_key --profile=swiply)
export AWS_SECRET_ACCESS_KEY

# Postgres users login and password section.
LOCKBOX_POSTGRES_PROFILE_OWNER_SECRET_ID=e6qgabf1uqm0v7hnu2ht
export POSTGRES_PROFILE_OWNER_LOGIN=$(yc lockbox payload get --id=$LOCKBOX_POSTGRES_PROFILE_OWNER_SECRET_ID --key=login --profile=swiply)
export POSTGRES_PROFILE_OWNER_PASSWORD=$(yc lockbox payload get --id=$LOCKBOX_POSTGRES_PROFILE_OWNER_SECRET_ID --key=password --profile=swiply)

LOCKBOX_POSTGRES_NOTIFICATION_OWNER_SECRET_ID=e6q6k6se7v3ebng55mj3
export POSTGRES_NOTIFICATION_OWNER_LOGIN=$(yc lockbox payload get --id=$LOCKBOX_POSTGRES_NOTIFICATION_OWNER_SECRET_ID --key=login --profile=swiply)
export POSTGRES_NOTIFICATION_OWNER_PASSWORD=$(yc lockbox payload get --id=$LOCKBOX_POSTGRES_NOTIFICATION_OWNER_SECRET_ID --key=password --profile=swiply)

LOCKBOX_POSTGRES_CHAT_OWNER_SECRET_ID=e6qjsphklgju06leifcq
export POSTGRES_CHAT_OWNER_LOGIN=$(yc lockbox payload get --id=$LOCKBOX_POSTGRES_CHAT_OWNER_SECRET_ID --key=login --profile=swiply)
export POSTGRES_CHAT_OWNER_PASSWORD=$(yc lockbox payload get --id=$LOCKBOX_POSTGRES_CHAT_OWNER_SECRET_ID --key=password --profile=swiply)

LOCKBOX_POSTGRES_EVENT_OWNER_SECRET_ID=e6qajklqsuoot9l8ti7m
export POSTGRES_EVENT_OWNER_LOGIN=$(yc lockbox payload get --id=$LOCKBOX_POSTGRES_EVENT_OWNER_SECRET_ID --key=login --profile=swiply)
export POSTGRES_EVENT_OWNER_PASSWORD=$(yc lockbox payload get --id=$LOCKBOX_POSTGRES_EVENT_OWNER_SECRET_ID --key=password --profile=swiply)

LOCKBOX_POSTGRES_RECOMMENDATION_OWNER_SECRET_ID=e6q4sut3m4gn23qbllom
export POSTGRES_RECOMMENDATION_OWNER_LOGIN=$(yc lockbox payload get --id=$LOCKBOX_POSTGRES_RECOMMENDATION_OWNER_SECRET_ID --key=login --profile=swiply)
export POSTGRES_RECOMMENDATION_OWNER_PASSWORD=$(yc lockbox payload get --id=$LOCKBOX_POSTGRES_RECOMMENDATION_OWNER_SECRET_ID --key=password --profile=swiply)

LOCKBOX_POSTGRES_RANDOMCOFFEE_OWNER_SECRET_ID=e6qghuqj7v0ajndteo2k
export POSTGRES_RANDOMCOFFEE_OWNER_LOGIN=$(yc lockbox payload get --id=$LOCKBOX_POSTGRES_RANDOMCOFFEE_OWNER_SECRET_ID --key=login --profile=swiply)
export POSTGRES_RANDOMCOFFEE_OWNER_PASSWORD=$(yc lockbox payload get --id=$LOCKBOX_POSTGRES_RANDOMCOFFEE_OWNER_SECRET_ID --key=password --profile=swiply)

LOCKBOX_POSTGRES_BACKOFFICE_OWNER_SECRET_ID=e6qsmtlrepblr8a90kob
export POSTGRES_BACKOFFICE_OWNER_LOGIN=$(yc lockbox payload get --id=$LOCKBOX_POSTGRES_BACKOFFICE_OWNER_SECRET_ID --key=login --profile=swiply)
export POSTGRES_BACKOFFICE_OWNER_PASSWORD=$(yc lockbox payload get --id=$LOCKBOX_POSTGRES_BACKOFFICE_OWNER_SECRET_ID --key=password --profile=swiply)

# Redis users password section.
LOCKBOX_REDIS_ADMIN_SECRET_ID=e6qj94omr0htgvi1i7mf
export REDIS_ADMIN_PASSWORD=$(yc lockbox payload get --id=$LOCKBOX_REDIS_ADMIN_SECRET_ID --key=password --profile=swiply)

# Initialize terraform variables.
cat templates/template.tfvars | envsubst > swiply.tfvars

export YC_TOKEN=$(yc iam create-token --profile=swiply)
terraform "$@" -var-file=swiply.tfvars

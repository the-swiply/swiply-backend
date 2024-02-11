import os

ENV_STAGE_KEY = "STAGE"

ENV_PROD = "prod"
ENV_DEV = "dev"


def is_prod():
    return os.environ.get(ENV_STAGE_KEY).lower() == ENV_PROD


def is_dev():
    return os.environ.get(ENV_STAGE_KEY).lower() == ENV_DEV

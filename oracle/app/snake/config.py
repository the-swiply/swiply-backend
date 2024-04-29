import os.path
import logging

import yaml

from app.snake import stage

DEV_CONFIG_NAME = "values-dev"
PROD_CONFIG_NAME = "values-prod"

__instance = None


def parse_yaml():
    if stage.is_dev():
        cfg_path = os.path.join("configs", DEV_CONFIG_NAME + ".yaml")
    elif stage.is_prod():
        cfg_path = os.path.join("configs", PROD_CONFIG_NAME + ".yaml")
    else:
        raise KeyError("unknown stage")

    with open(cfg_path, 'r') as f:
        global __instance
        __instance = yaml.full_load(f)
    print(__instance)


def get(key: str):
    if not __instance:
        logging.error("config instance not set")
    return __instance.get(key)

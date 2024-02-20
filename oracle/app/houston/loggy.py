import logging

__instance: logging.Logger = logging.getLogger(__name__)


def init(app_name: str):
    logger = logging.getLogger(app_name)
    logger.setLevel(logging.INFO)
    handler = logging.StreamHandler()
    handler.setFormatter(logging.Formatter('%(asctime)s - %(name)s - %(levelname)s - %(message)s'))
    logger.addHandler(handler)
    global __instance
    __instance = logger


def info(msg: str):
    __instance.info(msg)

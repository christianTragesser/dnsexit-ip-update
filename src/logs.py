import logging
import sys
from pythonjsonlogger import jsonlogger


def get_console_handler():
    console_handler = logging.StreamHandler(sys.stdout)
    formatter = jsonlogger.JsonFormatter()
    console_handler.setFormatter(formatter)
    return console_handler


def logger(name):
    logger = logging.getLogger(name)
    logger.setLevel(logging.INFO)
    logger.addHandler(get_console_handler())
    return logger

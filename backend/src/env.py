import os
from enum import Enum


class Env(Enum):
    MODEL = "MODEL"


class Config:

    @staticmethod
    def get(key: Env) -> str:
        return str(os.getenv(key.value))

import os
from enum import Enum


class Env(Enum):
    MODEL = "MODEL"
    LLM_URL = "LLM_URL"


class Config:

    @staticmethod
    def get(key: Env) -> str:
        return str(os.getenv(key.value))

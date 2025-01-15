from langchain_ollama.chat_models import ChatOllama

from env import Config, Env

llm_url = "http://host.docker.internal:11434"


class SutramAgent:

    @staticmethod
    def new() -> ChatOllama:

        model = Config.get(Env.MODEL)

        return ChatOllama(
            model=model,
            temperature=0,
            base_url=llm_url,
        )

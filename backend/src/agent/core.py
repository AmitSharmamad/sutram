from langchain_ollama.chat_models import ChatOllama

from env import Config, Env


class SutramAgent:

    @staticmethod
    def new() -> ChatOllama:

        model = Config.get(Env.MODEL)
        llm_url = Config.get(Env.LLM_URL)

        return ChatOllama(
            model=model,
            temperature=0,
            base_url=llm_url,
        )

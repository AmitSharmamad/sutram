from langchain_ollama.chat_models import ChatOllama

llm_url = "http://host.docker.internal:11434"


class SutramAgent:

    @staticmethod
    def new() -> ChatOllama:
        return ChatOllama(
            model="llama3.2",
            temperature=0,
            base_url=llm_url,
        )

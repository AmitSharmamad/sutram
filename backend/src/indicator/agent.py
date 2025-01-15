from langchain.agents import AgentExecutor, create_tool_calling_agent
from langchain_core.messages import HumanMessage
from langchain_core.prompts.chat import ChatPromptTemplate
from langchain_ollama.chat_models import ChatOllama
from langgraph.prebuilt import create_react_agent

from indicator.ema import calculate_ema
from indicator.sma import calculate_sma

llm_url = "http://host.docker.internal:11434"


def start():
    model = ChatOllama(
        model="llama3.2",
        temperature=0,
        base_url=llm_url,
    )

    tools = [calculate_sma, calculate_ema]

    expression = "ema(13) AND sma(10)"

    prompt = ChatPromptTemplate.from_messages(
        [
            (
                "system",
                "You are an expert Algo Trading Engine with extensive knowledge in Technical Indicators. Only invoke functions once. Don't describe any functionality",
            ),
            ("human", "{input}"),
            ("placeholder", "{agent_scratchpad}"),
        ]
    )

    agent = create_tool_calling_agent(model, tools, prompt)

    agent_executor = AgentExecutor(agent=agent, tools=tools, verbose=True)

    print(agent_executor.invoke({"input": expression}))

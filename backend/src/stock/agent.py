from langchain.agents import AgentExecutor, create_tool_calling_agent
from langchain_core.prompts.chat import ChatPromptTemplate

from agent.core import SutramAgent
from stock.tool import get_tools


def start():
    model = SutramAgent.new()

    tools = get_tools()

    expression = "1m timeframe on symbol AAPL"

    system_prompt = """
        You are an expert ClickHouse SQL developer with extensive knowledge in ClickHouse docs. Following is the schema of the stock table given below in ```. Project only the close value unless the other columns are specified. Don't give any explanation. Only make the function calls once.

        ```
        stock_(SYMBOL) [ stock_AAPL, stock_TSLA ]
            - timestamp [ PK ]
            - open
            - high
            - close
            - low
            - volume
            - created_at
            - updated_at
        ```

    """

    prompt = ChatPromptTemplate.from_messages(
        [
            (
                "system",
                system_prompt,
            ),
            ("human", "{input}"),
            ("placeholder", "{agent_scratchpad}"),
        ]
    )

    agent = create_tool_calling_agent(model, tools, prompt)

    agent_executor = AgentExecutor(agent=agent, tools=tools, verbose=True)

    print(agent_executor.invoke({"input": expression}))

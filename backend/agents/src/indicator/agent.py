from langchain.agents import AgentExecutor, create_tool_calling_agent
from langchain_community.llms import Ollama
from langchain_core.prompts.chat import ChatPromptTemplate

from agent.core import SutramAgent
from indicator.tool import get_tools


def start():
    model = SutramAgent.new()

    llm = Ollama(model="llama3.2", base_url="http://host.docker.internal:11434")

    prompt = """
        You are an Algo Testing expert.
        You have to convert the below expression in "" following the below schema in ```.

        expression: "EMA(13) Crosses Above SMA(9)"

        ```
{
  "$schema": "http://json-schema.org/draft-07/schema#",
  "title": "ExpressionStatement",
  "type": "object",
  "properties": {
    "type": {
      "type": "string",
      "enum": [
        "ExpressionStatement"
      ]
    },
    "returns": {
      "type": "string",
      "enum": [
        "bool"
      ]
    },
    "value": {
      "type": "object",
      "properties": {
        "type": {
          "type": "string",
          "enum": [
            "LogicalExpression"
          ]
        },
        "operator": {
          "type": "string",
          "enum": [
            "AND",
            "OR",
            "NA"
          ]
        },
        "returns": {
          "type": "string",
          "enum": [
            "bool"
          ]
        },
        "lhs": {
          "type": "object",
          "properties": {
            "type": {
              "type": "string",
              "enum": [
                "IndicatorExpression"
              ]
            },
            "operator": {
              "type": "string",
              "enum": [
                "CrossesAbove"
              ]
            },
            "returns": {
              "type": "string",
              "enum": [
                "bool"
              ]
            },
            "lhs": {
              "type": "object",
              "properties": {
                "type": {
                  "type": "string",
                  "enum": [
                    "Indicator"
                  ]
                },
                "name": {
                  "type": "string",
                  "enum": [
                    "EMA"
                  ]
                },
                "returns": {
                  "type": "string",
                  "enum": [
                    "float"
                  ]
                },
                "value": {
                  "type": "object",
                  "properties": {
                    "period": {
                      "type": "integer"
                    },
                    "offset": {
                      "type": "integer"
                    }
                  },
                  "required": [
                    "period",
                    "offset"
                  ]
                }
              },
              "required": [
                "type",
                "name",
                "returns",
                "value"
              ]
            },
            "rhs": {
              "type": "object",
              "properties": {
                "type": {
                  "type": "string",
                  "enum": [
                    "Number"
                  ]
                },
                "name": {
                  "type": "string",
                  "enum": [
                    "Number"
                  ]
                },
                "returns": {
                  "type": "string",
                  "enum": [
                    "float"
                  ]
                },
                "value": {
                  "type": "object",
                  "properties": {
                    "value": {
                      "type": "number"
                    }
                  },
                  "required": [
                    "value"
                  ]
                }
              },
              "required": [
                "type",
                "name",
                "returns",
                "value"
              ]
            }
          },
          "required": [
            "type",
            "operator",
            "returns",
            "lhs",
            "rhs"
          ]
        },
        "rhs": {
          "type": "object",
          "properties": {
            "type": {
              "type": "string",
              "enum": [
                "IndicatorExpression"
              ]
            },
            "operator": {
              "type": "string",
              "enum": [
                "CrossesAbove"
              ]
            },
            "returns": {
              "type": "string",
              "enum": [
                "bool"
              ]
            },
            "lhs": {
              "type": "object",
              "properties": {
                "type": {
                  "type": "string",
                  "enum": [
                    "Indicator"
                  ]
                },
                "name": {
                  "type": "string",
                  "enum": [
                    "EMA"
                  ]
                },
                "returns": {
                  "type": "string",
                  "enum": [
                    "float"
                  ]
                },
                "value": {
                  "type": "object",
                  "properties": {
                    "value": {
                      "type": "object",
                      "properties": {
                        "period": {
                          "type": "integer"
                        },
                        "offset": {
                          "type": "integer"
                        }
                      },
                      "required": [
                        "period",
                        "offset"
                      ]
                    }
                  },
                  "required": [
                    "value"
                  ]
                }
              },
              "required": [
                "type",
                "name",
                "returns",
                "value"
              ]
            },
            "rhs": {
              "type": "object",
              "properties": {
                "type": {
                  "type": "string",
                  "enum": [
                    "Number"
                  ]
                },
                "name": {
                  "type": "string",
                  "enum": [
                    "Number"
                  ]
                },
                "returns": {
                  "type": "string",
                  "enum": [
                    "float"
                  ]
                },
                "value": {
                  "type": "object",
                  "properties": {
                    "value": {
                      "type": "number"
                    }
                  },
                  "required": [
                    "value"
                  ]
                }
              },
              "required": [
                "type",
                "name",
                "returns",
                "value"
              ]
            }
          },
          "required": [
            "type",
            "operator",
            "returns",
            "lhs",
            "rhs"
          ]
        }
      },
      "required": [
        "type",
        "operator",
        "returns",
        "lhs"
      ]
    }
  },
  "required": [
    "type",
    "returns",
    "value"
  ]
}

example:

{
  "type": "ExpressionStatement",
  "returns": "bool",
  "value": {
    "type": "LogicalExpression",
    "operator": "AND",
    "returns": "bool",
    "lhs": {
      "type": "IndicatorExpression",
      "operator": "CrossesAbove",
      "returns": "bool",
      "lhs": {
        "type": "Indicator",
        "name": "EMA",
        "point": "close",
        "returns": "float",
        "value": {
          "period": 5,
          "offset": 0
        }
      },
      "rhs": {
        "type": "Number",
        "name": "Number",
        "point": "close",
        "returns": "float",
        "value": {
          "value": 21000
        }
      }
    },
    "rhs": {
      "type": "IndicatorExpression",
      "operator": "CrossesAbove",
      "returns": "bool",
      "lhs": {
        "type": "Indicator",
        "name": "EMA",
        "point": "close",
        "returns": "float",
        "value": {
          "value": {
            "period": 5,
            "offset": 0
          }
        }
      },
      "rhs": {
        "type": "Number",
        "name": "Number",
        "point": "close",
        "returns": "float",
        "value": {
          "value": 21000
        }
      }
    }
  }
}

        ```

    """

    response = llm.invoke(prompt)

    print(response)

    return

    tools = get_tools()

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

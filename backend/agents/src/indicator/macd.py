from langchain_core.tools import tool


@tool
def calculate_macd():
    """Calculates the MACD"""
    print("calculating the MACD")

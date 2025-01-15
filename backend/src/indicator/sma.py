from langchain_core.tools import tool


@tool
def calculate_sma(period: int):
    """Calculates the SMA or Simple Moving Average"""
    print("calculating the sma with period ", period)

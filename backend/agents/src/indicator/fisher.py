from langchain_core.tools import tool


@tool
def calculate_fisher():
    """Calculates the Fisher"""
    print("calculating the Fisher")

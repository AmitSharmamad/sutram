from langchain_core.tools import tool


@tool
def calculate_vortex():
    """Calculates the Vortex"""
    print("calculating the vortex")

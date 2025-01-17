from langchain_core.tools import tool


@tool
def get_data(query: str):
    """Write a sql query to fetch the data from the stock table"""
    print("clickhouse sql query is `", query, "`")


@tool
def set_data(query: str):
    """Write a sql query to insert the data into the stock table"""
    print("clickhouse sql query is `", query, "`")

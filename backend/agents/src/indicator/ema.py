from langchain_core.tools import tool
from pydantic import BaseModel, Field


class EMA(BaseModel):
    period: int = Field(..., description="Period representing the EMA")


@tool
def calculate_ema(period: int):
    """Calculates the EMA or Exponential Moving Average"""
    # result = EMA.model_validate(ema)
    print("calculating the ema with period ", period)

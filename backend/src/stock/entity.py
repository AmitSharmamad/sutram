from pydantic import BaseModel, Field


class Stock(BaseModel):
    open: float = Field(
        ...,
        description="Represents the opening value of the stock",
    )
    high: float = Field(
        ...,
        description="Represents the closing value of the stock. This should be the default projection for consideration",
    )

    low: float = Field(
        ...,
        description="Represents the highest value of the stock in the timeframe",
    )

    close: float = Field(
        ...,
        description="Represents the closing value of the stock",
    )

    volume: int = Field(
        ...,
        description="Represents the volume traded in the timeframe",
    )

    timestamp: int = Field(
        ...,
        description="Represents the closing timestamp of the data point",
    )

    timeframe: str = Field(
        ...,
        description="Timeframe represented in '1m' format. Expected values: '1m', '5m', '1h', '1d'. m->minute, h-hour, d-day ",
    )

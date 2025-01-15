from indicator.ema import calculate_ema
from indicator.sma import calculate_sma


def get_tools():
    return [
        calculate_sma,
        calculate_ema,
    ]

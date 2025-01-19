from indicator.ema import calculate_ema
from indicator.fisher import calculate_fisher
from indicator.macd import calculate_macd
from indicator.sma import calculate_sma
from indicator.vortex import calculate_vortex


def get_tools():
    return [
        calculate_sma,
        calculate_ema,
        calculate_macd,
        calculate_fisher,
        calculate_vortex,
    ]

create a module `migrate`, make one file `migrate.go` that does the following

make an array that stores the hardcoded values in the following format
```
{
"securityId": "1333",
"exchangeSegment": "NSE_EQ",
"instrument": "EQUITY",
"interval": "1",
"oi": false,
"fromDate": "2024-09-11 09:30:00",
"toDate": "2024-09-15 13:00:00"
}
```

use the below api to fetch the data
```
curl --request POST \
    --url https://api.dhan.co/v2/charts/intraday \
    --header 'Accept: application/json' \
    --header 'Content-Type: application/json' \
    --header 'access-token: ' \
    --data '{}'
```

response structure
```
{
    "open":   [
        3750,3757.85,3751.2
    ],
    "high":   [
        3750,3757.85,3751.2
    ],
    "low":      [
        3750,3757.85,3751.2
    ],
    "close":    [
        3750,3757.85,3751.2
    ],
    "volume":   [
        166,53629,34592,20802
    ],
    "timestamp":[
        1328845020,1328845500,1328845560,1328845620
    ],
    "open_interest":[
        0,0,0,0
    ]
}
```

for each of these values in the array do the following

1. fetch data from 2020 jan till today
2. fetch data in range of 2 months at once per request
3. once one year data is accumulated, write it to a file with the format `{securityId}_${year}.csv`
4. create the headers exchange_segment_id, open, high, low, close, volume, timestamp, created_at, updated_at
5. dates should be formatted in "2025-05-21 11:36:42" format
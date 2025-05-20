## Use the below libraries to build the project

```
github.com/ClickHouse/clickhouse-go/v2 v2.30.1
github.com/gin-contrib/cors v1.7.4
github.com/gin-gonic/gin v1.10.0
github.com/golang-jwt/jwt v3.2.2+incompatible
github.com/golang-jwt/jwt/v5 v5.2.1
github.com/google/uuid v1.6.0
github.com/google/wire v0.6.0
github.com/gorilla/websocket v1.5.3
github.com/jmoiron/sqlx v1.4.0
github.com/joho/godotenv v1.5.1
github.com/lib/pq v1.10.9
github.com/xeipuuv/gojsonschema v1.2.0
```

## Import Rules
1. Use algo/ as relative qualifier for imports
2. Don't use absolute imports for project files

## Struct Tags
1. Add json and db tags with snake_case for all the relevant model structs

## No Code Duplication
1. Do not create the same interface or struct twice.
2. For Example, if backtesting module needs StockService, do not create it inhouse. Instead, import it from Stock Module.
3. Before creating a interface or struct, check if it already exists

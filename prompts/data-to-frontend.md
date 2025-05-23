Generate a Go (Golang) HTTP server application.

**Requirements:**

1.  **Project Setup:** Create a `main.go` file and a basic `go.mod` file (e.g., `module trading-backend`).
2.  **Protobuf Integration:**
    * Assume there is a `proto/trading_data.proto` file (based on the `MarketData` and `ChartDataResponse` messages).
    * Provide instructions on how to install the necessary Go Protobuf tools (`protoc-gen-go`) and how to generate the Go code from the `.proto` file into a `gen/trading` directory.
    * Use the generated Go Protobuf code.
3.  **HTTP Server:**
    * Use the standard `net/http` library.
    * Create an HTTP server listening on port `8080`.
    * Implement an endpoint `/api/chart-data`.
4.  **Endpoint Logic (`/api/chart-data`):**
    * When this endpoint is requested:
        * Create a sample `ChartDataResponse` instance with at least 3-4 `MarketData` points (use example values).
        * Serialize this `ChartDataResponse` instance into binary Protobuf format using `proto.Marshal`.
        * Set the `Content-Type` header to `application/octet-stream`.
        * Set up CORS headers to allow GET requests from `http://localhost:3000` (typical Next.js dev port), including `Access-Control-Allow-Origin`, `Access-Control-Allow-Methods`, and `Access-Control-Allow-Headers`.
        * Write the serialized binary data as the HTTP response.
        * Include basic error handling and logging.
5.  **Structure:** Ensure the code is clean, well-structured, and includes necessary imports.
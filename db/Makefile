run:
	docker run -d --rm \
		-e CLICKHOUSE_USER=admin \
		-e CLICKHOUSE_DEFAULT_ACCESS_MANAGEMENT=1 \
		-e CLICKHOUSE_PASSWORD=admin \
		-p 9000:9000/tcp \
		-p 8123:8123/tcp \
		-v ./mount/data:/var/lib/clickhouse/ \
    -v ./mount/logs:/var/log/clickhouse-server/ \
    -v ./sql:/docker-entrypoint-initdb.d/ \
		--name sutramdb \
		clickhouse/clickhouse-server


gen:
	go run ./dev/main.go


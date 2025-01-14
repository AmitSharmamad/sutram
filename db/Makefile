build:
	docker build -t sutram_db .

# run:
# 	docker run -it --name sutram_db sutram_db

run:
	docker run -d --rm -e CLICKHOUSE_DB=sutram \
		-e CLICKHOUSE_USER=admin \
		-e CLICKHOUSE_DEFAULT_ACCESS_MANAGEMENT=1 \
		-e CLICKHOUSE_PASSWORD=admin \
		-p 9000:9000/tcp \
		-p 8123:8123/tcp \
		--name sutram_db \
		clickhouse/clickhouse-server

	docker cp ./sutram.sql sutram_db:/sutram.sql

	docker exec sutram_db clickhouse-client --multiquery < sutram.sql

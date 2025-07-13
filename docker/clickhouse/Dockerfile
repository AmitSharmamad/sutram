# Custom ClickHouse image with pre-configured settings
FROM clickhouse/clickhouse-server:latest

# Create necessary directories
RUN mkdir -p /etc/clickhouse-server/users.d/ \
    && mkdir -p /etc/clickhouse-server/config.d/ \
    && mkdir -p /var/lib/clickhouse \
    && mkdir -p /var/log/clickhouse-server

# Copy custom configuration files
COPY docker/clickhouse/config/users.xml /etc/clickhouse-server/users.d/
COPY docker/clickhouse/config/config.xml /etc/clickhouse-server/config.d/

# Copy initialization scripts
COPY docker/clickhouse/init/ /docker-entrypoint-initdb.d/

# Set proper permissions
RUN chown -R clickhouse:clickhouse /var/lib/clickhouse \
    && chown -R clickhouse:clickhouse /var/log/clickhouse-server \
    && chown -R clickhouse:clickhouse /etc/clickhouse-server

# Expose ports
EXPOSE 8123 9000 9009

# Use the default entrypoint from the base image
ENTRYPOINT ["/entrypoint.sh"]

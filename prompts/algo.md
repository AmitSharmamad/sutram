add the following columns to the strategy table

1. algo_id string -> exchange's algo_id
2. algo_id_sent_datetime timestamp -> our sent time
3. algo_id_received_datetime timestamp -> our received time
4. algo_id_status string -> exchange's status -> approved, rejected
5. algo_id_process_status enum(SENT, RECEIVED) -> to track the internal transition

fix the sql, repos, services and apis as needed
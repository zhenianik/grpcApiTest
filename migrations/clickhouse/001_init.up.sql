CREATE TABLE users (
        user_id UInt64,
        time DateTime,
        name String
) Engine = MergeTree
PARTITION BY toYYYYMM(time)
ORDER BY (user_id, time);

CREATE TABLE users_queue (
        user_id UInt64,
        time DateTime,
        name String
) ENGINE = Kafka
SETTINGS kafka_broker_list = 'kafka:9093',
       kafka_topic_list = 'users',
       kafka_group_name = 'users_consumer_group1',
       kafka_format = 'JSONEachRow',
       kafka_max_block_size = 1048576;
CREATE MATERIALIZED VIEW users_queue_mv TO users AS
SELECT user_id, time, name
FROM users_queue;
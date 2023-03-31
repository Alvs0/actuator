CREATE TABLE IF NOT EXISTS sensor
(
    first_id CHAR,
    second_id VARCHAR(10),
    timestamps TIMESTAMP,
    sensor_value FLOAT,
    sensor_type TEXT,

    PRIMARY KEY(first_id, second_id)
);
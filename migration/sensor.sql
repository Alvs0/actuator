CREATE TABLE IF NOT EXISTS sensor
(
    first_id CHAR,
    second_id CHAR,
    timestamps TIMESTAMP,
    sensor_value FLOAT,
    sensor_type TEXT,

    PRIMARY KEY(first_id, second_id)
);
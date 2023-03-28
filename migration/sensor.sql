CREATE TABLE IF NOT EXISTS "sensor"
(
    first_id CHAR,
    second_id CHAR,
    timestamp TIMESTAMPTZ,
    sensor_value FLOAT,

    PRIMARY KEY (first_id, second_id)
);

CREATE INDEX sensor_timestamp ON sensor(first_id, second_id, timestamp);

COMMENT ON TABLE "sensor" IS 'Holds sensor values information';
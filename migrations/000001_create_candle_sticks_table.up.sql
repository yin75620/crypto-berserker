CREATE TABLE candle_sticks (
    time_interval ENUM('1m', '3m', '5m', '15m', '30m', '1h', '2h', '4h', '6h', '8h', '12h', '1d', '3d', '1w') NOT NULL,

    open_time BIGINT NOT NULL,
    open DECIMAL(20,8) NOT NULL,
    high DECIMAL(20,8) NOT NULL,
    low DECIMAL(20,8) NOT NULL,
    close DECIMAL(20,8) NOT NULL,
    volume DECIMAL(20,8) NOT NULL,
    close_time BIGINT NOT NULL,
    quote_asset_volume DECIMAL(20,8),
    trade_count INT,
    taker_buy_asset_volume DECIMAL(20,8),
    taker_buy_quote_asset_volume DECIMAL(20,8),
    ignore_field DECIMAL(20,8),

    PRIMARY KEY (time_interval, open_time)
);

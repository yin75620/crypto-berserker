CREATE TABLE candlesticks (
    id INT AUTO_INCREMENT PRIMARY KEY,
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
    ignore_field DECIMAL(20,8)
);

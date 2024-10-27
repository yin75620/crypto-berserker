## golang-migrate
```
brew install golang-migrate
```

### create migration file
```
mkdir migrations
migrate create -ext sql -dir ./migrations -seq create_candlesticks_table
```

### execute migration file
```
migrate -path ./migrations -database "mysql://root:password@tcp(127.0.0.1:3306)/exchange_data" up
```

`source .env`
```
migrate -path ./migrations -database "mysql://$MYSQL_USER:$MYSQL_PWD@tcp($MYSQL_HOST:$MYSQL_PORT)/$MYSQL_DB" up
```
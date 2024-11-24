module github.com/Meduzz/abstractions.go

go 1.20

require github.com/go-redis/redis/v8 v8.11.5

require github.com/golang-jwt/jwt/v5 v5.2.1

require (
	filippo.io/edwards25519 v1.1.0 // indirect
	github.com/go-sql-driver/mysql v1.8.1 // indirect
	github.com/google/uuid v1.6.0 // indirect
	github.com/jinzhu/inflection v1.0.0 // indirect
	github.com/jinzhu/now v1.1.5 // indirect
	github.com/mattn/go-sqlite3 v1.14.22 // indirect
	golang.org/x/text v0.20.0 // indirect
	gorm.io/driver/mysql v1.5.7 // indirect
)

require (
	github.com/Meduzz/helper v0.0.0-20240730101358-04a510a685f3
	github.com/cespare/xxhash/v2 v2.3.0 // indirect
	github.com/dgryski/go-rendezvous v0.0.0-20200823014737-9f7001d12a5f // indirect
	gorm.io/datatypes v1.2.4
	gorm.io/driver/sqlite v1.5.6
	gorm.io/gorm v1.25.12
)

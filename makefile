server:
	go run main.go

db:
	docker run --env=MYSQL_ROOT_PASSWORD=password --env=MYSQL_DATABASE=blognado -p 4001:3306 -d mysql:8.0.31

redis:
	docker run --name redis -p 6379:6379 -d redis:7.2-rc1-alpine

ping-redis:
	docker exec -it redis redis-cli ping

migrateup:
	migrate -path db/migration -database "mysql://root:password@tcp(127.0.0.1:4001)/blognado?parseTime=true" -verbose up

migratedown:
	migrate -path db/migration -database "mysql://root:password@tcp(127.0.0.1:4001)/blognado?parseTime=true" -verbose down

.PHONY: server db redis migrateup migratedown
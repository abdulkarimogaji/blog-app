server:
	go run main.go

db:
	docker run --env=MYSQL_ROOT_PASSWORD=password --env=MYSQL_DATABASE=blognado -p 4001:3306 -d mysql:8.0.31
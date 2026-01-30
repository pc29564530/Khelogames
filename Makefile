DB_URL=postgresql://root:Bharat@12@localhost:5432/khelogames?sslmode=disable

#network:
#	docker network create bank-network
#
#postgres:
#	docker run --name postgres14 -p 5432:5432 -e POSTGRES_USER=root -e POSTGRES_PASSWORD=Bharat@12 -v ~/database/khelogames:/var/lib/postgresql/data -d postgres:14-alpine
#
#mysql:
#	docker run --name mysql8 -p 3306:3306  -e MYSQL_ROOT_PASSWORD=secret -d mysql:8
#
createDb:
	docker exec -it postgres createdb --username=root --owner=root Khelogames

dropDb:
	docker exec -it postgres dropdb Khelogames

migrateUp1: 
	migrate -path db/migration -database postgresql://root:Bharat@12@localhost:5432/khelogames?sslmode=disable -verbose u

migrateUp:
	migrate -path db/migration -database "$(DB_URL)" -verbose up

migrateUp1:
	migrate -path db/migration -database "$(DB_URL)" -verbose up 1

migrateDown:
	migrate -path db/migration -database "$(DB_URL)" -verbose down

migrateDown1:
	migrate -path db/migration -database "$(DB_URL)" -verbose down 1

migrateUp2:
	migrate -path db/migration -database "$(DB_URL)" -verbose up 2

migrateDown2:
	migrate -path db/migration -database "$(DB_URL)" -verbose down 2

migrateUp3:
	migrate -path db/migration -database "$(DB_URL)" -verbose up 3

migrateDown3:
	migrate -path db/migration -database "$(DB_URL)" -verbose down 3

migrateUp4:
	migrate -path db/migration -database "$(DB_URL)" -verbose up 4

migrateDown4:
	migrate -path db/migration -database "$(DB_URL)" -verbose down 4

new_migration:
	migrate create -ext sql -dir db/migration -seq $(name)

db_docs:
	dbdocs build doc/db.dbml

db_schema:
	dbml2sql --postgres -o doc/schema.sql doc/db.dbml

sqlc:
	sqlc generate

test:
	go test -v -cover -short ./...

server:
	go run main.go

mock:
	mockgen -package mockdb -destination db/mock/store.go github.com/techschool/simplebank/db/sqlc Store
	mockgen -package mockwk -destination worker/mock/distributor.go github.com/techschool/simplebank/worker TaskDistributor

proto:
	rm -f pb/*.go
	rm -f doc/swagger/*.swagger.json
	protoc --proto_path=proto --go_out=pb --go_opt=paths=source_relative \
	--go-grpc_out=pb --go-grpc_opt=paths=source_relative \
	--grpc-gateway_out=pb --grpc-gateway_opt=paths=source_relative \
	--openapiv2_out=doc/swagger --openapiv2_opt=allow_merge=true,merge_file_name=simple_bank \
	proto/*.proto
	statik -src=./doc/swagger -dest=./doc

evans:
	evans --host localhost --port 9090 -r repl

redis:
	docker run --name redis -p 6379:6379 -d redis:7-alpine

rabbitmq:
	docker run -d --hostname rabbit --name rabbit-mq -p 5672:5672 -p 15672:15672 rabbitmq:3-management

# Docker Compose commands
docker-build:
	docker build -t khelogames-backend .

docker-up:
	docker-compose up -d

docker-down:
	docker-compose down

docker-down-volumes:
	docker-compose down -v

docker-logs:
	docker-compose logs -f backend

docker-logs-all:
	docker-compose logs -f

docker-restart:
	docker-compose restart backend

docker-ps:
	docker-compose ps

docker-exec:
	docker-compose exec backend sh

docker-db:
	docker-compose exec postgres psql -U root -d Khelogames

docker-backup-db:
	docker-compose exec postgres pg_dump -U root Khelogames > backup_$(shell date +%Y%m%d_%H%M%S).sql

docker-clean:
	docker-compose down -v
	docker system prune -f

.PHONY: network postgres createdb dropdb migrateup migratedown migrateUp1 migrateDown1 new_migration db_docs db_schema sqlc test server mock proto evans redis migrateUp2 migrateUp3 migrateUp4 docker-build docker-up docker-down docker-down-volumes docker-logs docker-logs-all docker-restart docker-ps docker-exec docker-db docker-backup-db docker-clean
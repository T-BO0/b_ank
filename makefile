.PHONY:
	create_db drop_db start_postgres_docker_container stop_postgres_docker_container start_new_postgres_container migrate_up migrate_down sqlc_g test

create_db:
	docker exec -it postgres12 createdb --username=root --owner=root bank

drop_db:
	docker exec -it postgres12 dropdb bank

start_postgres_docker_container:
	docker start postgres12

stop_postgres_docker_container:
	docker stop postgres12

start_new_postgres_container:
	docker run --name postgres12 -e POSTGRES_USER=root -e POSTGRES_PASSWORD=secret -p 5432:5432 -d 7c8f48705831

migrate_up:
	migrate  -path db/migrations -database "postgresql://root:secret@localhost:5432/bank?sslmode=disable" -verbose up

migrate_down:
	migrate  -path db/migrations -database "postgresql://root:secret@localhost:5432/bank?sslmode=disable" -verbose down

sqlc_g:
	sqlc generate

test:
	go test -v -cover ./...

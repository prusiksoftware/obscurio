psql_chinook:
	docker compose exec db psql -U user -d chinook -h localhost -p 5432 -W

psql_proxy:
	docker compose exec db psql -U user -d chinook -h proxy -p 5432 -W

psql_proxy_user1:
	docker compose exec db psql -U user -d chinook -h proxy -p 5432 -W

psql_proxy_user2:
	docker compose exec db psql -U user2 -d chinook -h proxy -p 5432 -W

reset_db:
	docker compose stop db
	sudo rm -rf ./example_db/data
	docker compose up -d db

test:
	docker compose exec proxy go test ./... -v

coverage:
	docker compose up -d
	docker compose exec proxy go test ./... -v -coverprofile=./cover.out -timeout 2s -covermode=atomic -coverpkg=./...
	go tool cover -html=cover.out -o coverage.html

pg_dump:
	pg_dump -U user -h localhost -p 5432 -d chinook -vv -f ./dump.sql

purge_docker_logs:
	sudo find /var/lib/docker/containers/ -type f -name "*.log" -exec rm -f {} \;

find_deadcode:
	staticcheck ./...
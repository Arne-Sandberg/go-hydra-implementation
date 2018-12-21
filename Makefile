docker_compose_file = build/docker-compose.yml
project_name = go-hydra-implementation
server_name = go-hydra-server

default:
	docker-compose -p $(project_name) -f $(docker_compose_file) build

up build:
	docker-compose -p $(project_name) -f $(docker_compose_file) up # --build

down:
	docker-compose -p $(project_name) -f $(docker_compose_file) down

restart:
	make down && make up

hard_restart:
	make down && make && make up

logs:
	docker logs -f $(server_name)

shell:
	docker exec -it $(server_name) sh

fmt:
	docker exec -it $(server_name) gofmt -s -w .

test:
	docker exec -it $(server_name) go test -v
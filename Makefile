.PHONY: setup start stop restart build logs clean seed health ps shell

setup:
	@./scripts/setup.sh

start:
	docker-compose up -d

stop:
	docker-compose down

restart: stop start

build:
	docker-compose build

logs:
	docker-compose logs -f

logs-svc:
	docker-compose logs -f $(SVC)

clean:
	docker-compose down -v --remove-orphans

seed:
	@./scripts/seed-data.sh

health:
	@./scripts/health-check.sh

ps:
	docker-compose ps

shell:
	docker-compose exec $(SVC) /bin/sh || docker-compose exec $(SVC) /bin/bash

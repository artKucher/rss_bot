# Константы

TAIL=100

# Переменные, которые используются в скриптах, для использование нужно указать ее при запуске
# Например make logs c=celery


define set-default-container
	ifndef c
	c = rss_bot
	else ifeq (${c},all)
	override c=
	endif
endef


set-container:
	$(eval $(call set-default-container))


build:
	docker compose -f docker-compose.dev.yml build
up:
	docker compose -f docker-compose.dev.yml up --remove-orphans  -d $(c)
down:
	docker compose -f docker-compose.dev.yml down
logs: set-container
	docker compose -f docker-compose.dev.yml logs --tail=$(TAIL) -f $(c)
restart: set-container
	docker compose -f docker-compose.dev.yml restart $(c)
exec: set-container
	docker compose -f docker-compose.dev.yml exec $(c) /bin/bash
remove: set-container
	docker compose -f docker-compose.dev.yml rm -fs $(c)
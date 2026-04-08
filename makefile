# Переменные для удобства (легко менять пути в одном месте)
CONFIG_PATH=./config/local.yaml
APP_PATH=cmd/sso/main.go

.PHONY: run
# Команда для запуска приложения
run:
	go run $(APP_PATH) --config=$(CONFIG_PATH)

.PHONY: build
# Команда для сборки бинарного файла (полезно для проверки ошибок компиляции)
build:
	go build -o sso $(APP_PATH)

.PHONY: help
# Подсказка по доступным командам
help:
	@echo "Доступные команды:"
	@echo "  make run   - запустить приложение с локальным конфигом"
	@echo "  make build - собрать исполняемый файл"

.PHONY: migrat

migrate:
	go run ./cmd/migrator \
		-migrations-path=./migrations \
		-host=localhost \
		-port=5433 \
		-user=postgres_auth \
		-password=1234 \
		-dbname=sso
		
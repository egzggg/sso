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

migrat:
	go run ./cmd/migrator --storage-path=./storage/sso.db --migrations-path=./migrations
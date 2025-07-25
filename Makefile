# Makefile для pm — простого пакетного менеджера

# Имя бинарника
BINARY := pm

# Пути к конфигам (можно изменить)
CREATE_CONFIG ?= ./config/packet.json
UPDATE_CONFIG ?= ./config/packages.json

# Цель по умолчанию
.DEFAULT_GOAL := help

# Проверка, установлен ли Go
ifeq ($(shell which go),)
  $(error "Go не установлен. Установите Go, чтобы продолжить.")
endif

# Цели

## Собрать бинарник
build:
	@echo "Собираем бинарник..."
	@go build -o $(BINARY)
	@sudo mv pm /usr/local/bin/
	@echo "Бинарник готов: ./${BINARY}"
	@echo "Добавлен в usr/bin"
	

## Установить зависимости и почистить go.mod
mod:
	@echo "Обновляем зависимости..."
	@go mod tidy
	@echo "go.mod и go.sum обновлены"


update: build
	@echo "🔁 Запускаем update..."
	@if [ -f "$(UPDATE_CONFIG)" ]; then \
        ./$(BINARY) update $(UPDATE_CONFIG); \
    else \
        echo "❌ Файл $(UPDATE_CONFIG) не найден"; \
        exit 1; \
    fi


clean:
	@rm -f $(BINARY)
	@echo "🗑️  Бинарник удалён"


help:
	@echo ""
	@echo "Доступные команды:"
	@echo ""
	@echo "  make build        — Собрать бинарник"
	@echo "  make mod          — Обновить зависимости (go mod tidy)"


	@echo "  make clean        — Удалить бинарник"
	@echo "  make help         — Показать это сообщение"
	@echo ""
	@echo ""

# Сокращённые цели
.PHONY: build mod create update clean help
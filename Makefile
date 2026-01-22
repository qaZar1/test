PROJECT := $(shell cat go.mod | grep module | awk -F ' ' '{print $$2}' | awk -F '/' '{print $$NF}')

REPORTS := .reports
VENDOR := vendor
AUTOGEN := autogen

dependencies:
	[ -x "$$(command -v mockgen)" ] || go install github.com/golang/mock/mockgen@latest
	[ -x "$$(command -v oapi-codegen)" ] || go install github.com/deepmap/oapi-codegen/cmd/oapi-codegen@latest
	[ -x "$$(command -v redoc-cli)" ] || npm i -g redoc-cli

$(REPORTS):
	mkdir -p $(REPORTS)

help:
	@egrep -h '\s##\s' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

setup: dependencies $(REPORTS) ## Настроить рабочее окружение

clean: ## Очистить рабочее окружение
	rm -rf $(VENDOR)
	rm -rf $(REPORTS)
	rm -rf $(AUTOGEN)
	go clean -r -i -testcache

mock-gen: setup ## Генерация mock
	go generate ./...

tests: mock-gen ## Запуск авто-тестов
	go test -race -cover -coverprofile $(REPORTS)/coverage.out ./...
	go tool cover -html=$(REPORTS)/coverage.out -o $(REPORTS)/coverage.html

oapi-yaml-gen: setup ## Генерация YAML из OpenAPI спецификации
	mkdir -p $(AUTOGEN)/docs
	sed -e 's/^/    /' README.md > $(AUTOGEN)/docs/README.md
	sed -e '/    %README.md%/{' -e "r $(AUTOGEN)/README.md" -e 'd' -e '}' oapi3.yaml > $(AUTOGEN)/docs/oapi3.yaml

oapi-doc-gen: setup oapi-yaml-gen ## Генерация документации из OpenAPI спецификации
	redoc-cli build $(AUTOGEN)/docs/oapi3.yaml -o $(AUTOGEN)/docs/user-doc.html

oapi-code-gen: setup oapi-doc-gen ## Генерация кода из OpenAPI спецификации
	oapi-codegen -o $(AUTOGEN)/server.go -old-config-style -package autogen -generate chi-server $(AUTOGEN)/docs/oapi3.yaml
	oapi-codegen -o $(AUTOGEN)/types.go  -old-config-style -package autogen -generate types $(AUTOGEN)/docs/oapi3.yaml
	oapi-codegen -o $(AUTOGEN)/client.go -old-config-style -package autogen -generate client $(AUTOGEN)/docs/oapi3.yaml

up: ## Запуск приложения в docker
	docker compose -f 'docker-compose.yml' --env-file ./cmd/wallet/config.env up -d --build

all: oapi-code-gen tests ## Последовательный запуск основных команд

.DEFAULT_GOAL := help
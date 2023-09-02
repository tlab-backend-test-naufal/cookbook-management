export LINTER_VERSION ?= 1.43.0

GO_PACKAGES ?= $(shell go list ./... | grep -v 'examples\|qtest\|mock')
TMP_DIR     := $(shell mktemp -d)
MODULES      = $(shell cd module && \ls -d */)
UNAME       := $(shell uname)

# Default postgres migration settings
export POSTGRES_USER 	?= tlab
export POSTGRES_PASS 	?= tlab
export POSTGRES_HOST 	?= localhost
export POSTGRES_PORT 	?= 5439
export POSTGRES_DB   	?= cookbook_management
export POSTGRES_SSLMODE ?= disable

# If the first argument is "migrate"...
ifeq (migrate,$(firstword $(MAKECMDGOALS)))
    # use the rest as arguments for "migrate"
    MIGRATE_ARGS := $(wordlist 2,$(words $(MAKECMDGOALS)),$(MAKECMDGOALS))
    # ...and turn them into do-nothing targets
    $(eval $(MIGRATE_ARGS):;@:)
endif

all: check test

tool-lint: bin
	@test -e ./bin/golangci-lint || curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b ./bin v${LINTER_VERSION}

check: tool-lint
	./bin/golangci-lint run -v --timeout 3m0s

tool-migrate: bin
ifeq ($(UNAME), Linux)
	@curl -sSfL https://github.com/golang-migrate/migrate/releases/download/v4.15.1/migrate.linux-amd64.tar.gz | tar zxf - --directory /tmp \
	&& cp /tmp/migrate bin/
else ifeq ($(UNAME), Darwin)
	@curl -sSfL https://github.com/golang-migrate/migrate/releases/download/v4.15.1/migrate.darwin-amd64.tar.gz | tar zxf - --directory /tmp \
	&& cp /tmp/migrate bin/
else
	@echo "Your OS is not supported."
endif

test:
	@go test -race -v ./module/...

cover:
	@go tool cover -html=coverage.out

migrate: tool-migrate
	@$(foreach module, $(MODULES), cp module/$(module)/db/migrations/*.sql $(TMP_DIR) 2>/dev/null;)
	@bin/migrate -source file://$(TMP_DIR) -database "postgres://$(POSTGRES_USER):$(POSTGRES_PASS)@$(POSTGRES_HOST):$(POSTGRES_PORT)/$(POSTGRES_DB)?sslmode=$(POSTGRES_SSLMODE)" $(MIGRATE_ARGS)

compose-up:
	@docker-compose up

compose-down:
	@docker-compose down
PKG_VERSION  = $(shell cat VERSION)
TEST_FOLDERS = $(shell go list ./... | grep "journal-migrator/test")

export DATABASE_USER=journal_user
export DATABASE_PASS=journal_pass
export DATABASE_HOST=localhost


.PHONY: check
check:
	@echo "Checking code format"
	@gofmt -l .
	@test -z "$(shell gofmt -l .)"
	@echo "Checking dependencies"
	@test -z "$(shell go mod tidy -diff)"


.PHONY: install-dev
install-dev:
	@echo "Installing development packages"
	@go mod download


.PHONY: services-up
services-up:
	@echo "Starting services"
	@docker compose --file "compose/$(COMPOSE_FILE)" up --wait


.PHONY: services-down
services-down:
	@echo "Stopping services"
	@docker compose --file "compose/$(COMPOSE_FILE)" down


.PHONY: test-unit
test-unit:
	@echo "Running unit tests"
	@go test -short $(TEST_FOLDERS)


.PHONY: test-integration
test-integration:
	@echo "Running integration tests"
	@go test $(TEST_FOLDERS)


.PHONY: tag
tag:
	@echo "Tagging current commit"
	@git tag --annotate "v$(PKG_VERSION)" --message "Tag v$(PKG_VERSION)"
	@git push --follow-tags

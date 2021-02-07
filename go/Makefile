# Sets the project dir correctly when run in Docker
PROJECT_DIR	= ${CURDIR}/..

# Default GO_BIN to Go binary in PATH
DOCKER_BIN			?= docker
DOCKER_TARGET_CMD	= $(word 2, $(subst /, ,$(@)))

# Default GO_BIN to Go binary in PATH
GO_BIN				?= go
GO_VERSION			?= 1.14
GO_DOCKER_CMD		?= $(DOCKER_BIN) run --rm -it -v ${PROJECT_DIR}:/app -w /app/go golang:${GO_VERSION}

TEST_PATTERN ?=.
TEST_OPTIONS ?=
SOURCE_FILES ?= ./...

TEST_FLAGS += -failfast
TEST_FLAGS += -race
TEST_TIMEOUT ?= 10m

GO_TEST ?= test $(TEST_OPTIONS) $(TEST_FLAGS) $(SOURCE_FILES) -run $(TEST_PATTERN) -timeout=$(TEST_TIMEOUT)

.PHONY: go-get
go-get:
	@printf '\n================================================================\n'
	@printf 'Target: go-get'
	@printf '\n================================================================\n'
	$(GO_BIN) mod vendor
	@echo '[go-get] Done.'

.PHONY: docker/go-get
docker/go-get:
	@($(GO_DOCKER_CMD) make $(DOCKER_TARGET_CMD))

.PHONY: test-coverage
test-coverage: TEST_FLAGS += -covermode=atomic -coverprofile=coverage.out
test-coverage: go-get
	@printf '\n================================================================\n'
	@printf 'Target: test-coverage'
	@printf '\n================================================================\n'
	@echo '[test] Testing packages: $(SOURCE_FILES)'
	$(GO_BIN) $(GO_TEST)

.PHONY: docker/test-coverage
docker/test-coverage:
	@($(GO_DOCKER_CMD) make $(DOCKER_TARGET_CMD))

.PHONY: test
test: go-get
	@printf '\n================================================================\n'
	@printf 'Target: test'
	@printf '\n================================================================\n'
	$(GO_BIN) $(GO_TEST)

.PHONY: docker/test
docker/test:
	@($(GO_DOCKER_CMD) make $(DOCKER_TARGET_CMD))

.PHONY: test-only
test-only: TEST_FLAGS += -v
test-only: fmt
	@printf '\n================================================================\n'
	@printf 'Target: test-only'
	@printf '\n================================================================\n'
	$(GO_BIN) $(GO_TEST)

.PHONY: docker/test-only
docker/test-only:
	@($(GO_DOCKER_CMD) make $(DOCKER_TARGET_CMD))

.PHONY: fmt
fmt:
	@printf '\n================================================================\n'
	@printf 'Target: fmt'
	@printf '\n================================================================\n'
	$(GO_BIN) fmt ./...

.PHONY: docker/fmt
docker/fmt:
	@($(GO_DOCKER_CMD) make $(DOCKER_TARGET_CMD))
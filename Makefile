.PHONY: test test-unit test-integration test-e2e test-load

test: test-unit test-integration test-e2e

test-unit:
    go test -v ./tests/unit/...

test-integration:
    go test -v ./tests/integration/...

test-e2e:
    go test -v ./tests/e2e/...

test-load:
    node ./tests/load/notification_load_test.js 
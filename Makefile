.PHONY: mock
mock:
	@mockgen -package=redismocks -destination=internal/mocks/redismocks/cmdable.mock.go github.com/redis/go-redis/v9 Cmdable
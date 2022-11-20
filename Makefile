run-local:
	go run main.go start --config=config/local.toml

start-docker:
	docker-compose -f ./docker/stack.yml up

stop-docker:
	docker-compose -f ./docker/stack.yml down
.PHONY: run
# run
run:
	go build main.go

.PHONY: docs
# run
docs:
	swag init


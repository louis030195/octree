bench:
	go test ./... -bench=. -benchmem -benchtime 1000000x

test:
	go test ./...
	@echo 'Test passed'

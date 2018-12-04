.PHONY: fotmat check-style lint check-error

format:
	@find $(CURDIR) -mindepth 1 -maxdepth 1 -type d -not -name vendor -not -name .git -print0 | xargs -0 gofmt -s -w
	@find $(CURDIR) -maxdepth 1 -type f -name '*.go' -print0 | xargs -0 gofmt -s -w

check-style:
	@golangci-lint run --disable-all -E gofmt ./...

lint:
	@golangci-lint run ./...

check-error:
	@golangci-lint run --disable-all -E errcheck ./...

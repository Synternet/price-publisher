BINARY_NAME=price-publisher
LDFLAGS="-w -s"

.PHONY: build build-static fmt watch

build:
	go build -o . ./...

build-static:
	CGO_ENABLED=1 go build -race -v -o $(BINARY_NAME) -a -installsuffix cgo -ldflags $(LDFLAGS) ./...

fmt:
	gofumpt -l -w .

watch:
	@if [ -x "$(GOPATH)/bin/air" ]; then \
	    "$(GOPATH)/bin/air"; \
		@echo "Watching...";\
	else \
	    read -p "air is not installed. Do you want to install it now? (y/n) " choice; \
	    if [ "$$choice" = "y" ]; then \
			go install github.com/cosmtrek/air@latest; \
	        "$(GOPATH)/bin/air"; \
				@echo "Watching...";\
	    else \
	        echo "You chose not to install air. Exiting..."; \
	        exit 1; \
	    fi; \
	fi

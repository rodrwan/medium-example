VERSION=v0.0.1
SVC=medium-example
BIN=$(PWD)/bin/$(SVC)

GO ?= go
LDFLAGS='-extldflags "static" -X main.svcVersion=$(VERSION) -X main.svcName=$(SVC)'
TAGS=netgo -installsuffix netgo

REGISTRY_URL=gotoschool

run r:
	@echo "[running] Running service..."
	@go run cmd/server/main.go

build b:
	@echo "[build] Building service..."
	@cd cmd/server && $(GO) build -o $(BIN) -ldflags=$(LDFLAGS) -tags $(TAGS)

linux l:
	@echo "[build] Building for linux..."
	@cd cmd/server && GOOS=linux $(GO) build -a -o $(BIN) --ldflags $(LDFLAGS) -tags $(TAGS)

docker d: linux
	@echo "[docker] Building image..."
	@docker build -t $(REGISTRY_URL)/$(SVC):$(VERSION) .

push: docker
	@echo "[docker] pushing $(REGISTRY_URL)/$(SVC):$(VERSION)"
	docker push $(REGISTRY_URL)/$(SVC):$(VERSION)
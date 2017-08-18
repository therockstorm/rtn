HANDLER ?= handler
PACKAGE ?= $(HANDLER)
BUILD ?= build/
GOPATH  ?= $(HOME)/go

docker:
	@docker run --rm                                                             \
	  -e HANDLER=$(HANDLER)                                                      \
	  -e PACKAGE=$(PACKAGE)                                                      \
	  -e GOPATH=$(GOPATH)                                                        \
	  -v $(CURDIR):$(CURDIR)                                                     \
	  $(foreach GP,$(subst :, ,$(GOPATH)),-v $(GP):$(GP))                        \
	  -w $(CURDIR)                                                               \
	  eawsy/aws-lambda-go-shim:latest make all

.PHONY: docker

all: clean build pack perm

.PHONY: all

build:
	@cd cmd && go build -buildmode=plugin -ldflags='-w -s' -o $(HANDLER).so && cd ..

.PHONY: build

pack:
	@pack $(HANDLER) cmd/$(HANDLER).so $(PACKAGE).zip

.PHONY: pack

perm:
	@chown $(shell stat -c '%u:%g' .) cmd/$(HANDLER).so $(PACKAGE).zip
	@mkdir $(BUILD) && mv $(PACKAGE).zip $(BUILD)$(PACKAGE).zip && cp serverless.yml $(BUILD)serverless.yml

.PHONY: perm

deploy:
	@cd $(BUILD) && serverless deploy --stage $(NODE_ENV) --verbose && cd ..

.PHONY: deploy

clean:
	@rm -rf $(BUILD) && rm -f cmd/$(HANDLER).so

.PHONY: clean

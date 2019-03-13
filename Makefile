REGISTRY=todefine
NAME=rancher-operator

.PHONY: build
build:
	operator-sdk build $(REGISTRY)/$(NAME)

push-test:
	docker push $(REGISTRY)/$(NAME)

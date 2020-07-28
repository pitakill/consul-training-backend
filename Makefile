.PHONY: create deploy update-version all

image = pitakill/consul-training-backend

all: update-version create deploy

update-version:
	@echo -e "package main\n\nconst VERSION = \"$(version)\"" > version.go

create:
	docker build -t $(image):$(version) .

deploy:
	docker push $(image):$(version)
	docker tag $(image):$(version) $(image):latest
	docker push $(image):latest

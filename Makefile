.PHONY: create deploy all

image = pitakill/consul-training-backend

all: create deploy

create:
	docker build -t $(image) .

deploy:
	docker push $(image)

.PHONY: create deploy all

image = pitakill/consul-training-backend

all: create deploy

create:
	docker build -t $(image):$(version) .

deploy:
	docker push $(image):$(version)
	docker tag $(image):$(version) $(image):latest
	docker push $(image):latest

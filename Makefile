.PHONY: create deploy all

image = pitakill/consul-training-backend

create:
	docker build -t $(image) .

deploy:
	docker push $(image)

all: create deploy

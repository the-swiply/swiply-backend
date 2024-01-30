export DOCKER_DEFAULT_PLATFORM=linux/amd64

YC_CONTAINER_REGISTRY := cr.yandex/crpf76jp63emqup99s4l

.PHONY: up
up:
	docker-compose up --build

.PHONY: tf-apply
tf-apply:
	cd ./deployments/terraform && bash tf.sh apply

.PHONY: tf-destroy
tf-destroy:
	cd ./deployments/terraform && bash tf.sh destroy

.PHONY: push-user-image
push-user-image:
	docker build -f ./user/Dockerfile --tag $(YC_CONTAINER_REGISTRY)/user:stable . && docker push cr.yandex/crpf76jp63emqup99s4l/user:stable

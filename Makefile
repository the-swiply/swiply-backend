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
	DOCKER_DEFAULT_PLATFORM=linux/amd64	docker build -f ./user/Dockerfile --tag $(YC_CONTAINER_REGISTRY)/user:stable . && docker push cr.yandex/crpf76jp63emqup99s4l/user:stable

.PHONY: push-profile-image
push-profile-image:
	DOCKER_DEFAULT_PLATFORM=linux/amd64	docker build -f ./profile/Dockerfile --tag $(YC_CONTAINER_REGISTRY)/profile:stable . && docker push cr.yandex/crpf76jp63emqup99s4l/profile:stable

.PHONY: push-randomcoffee-image
push-randomcoffee-image:
	DOCKER_DEFAULT_PLATFORM=linux/amd64	docker build -f ./randomcoffee/Dockerfile --tag $(YC_CONTAINER_REGISTRY)/randomcoffee:stable . && docker push cr.yandex/crpf76jp63emqup99s4l/randomcoffee:stable

.PHONY: push-notification-image
push-notification-image:
	DOCKER_DEFAULT_PLATFORM=linux/amd64	docker build -f ./notification/Dockerfile --tag $(YC_CONTAINER_REGISTRY)/notification:stable . && docker push cr.yandex/crpf76jp63emqup99s4l/notification:stable

.PHONY: push-backoffice-image
push-backoffice-image:
	DOCKER_DEFAULT_PLATFORM=linux/amd64	docker build -f ./backoffice/Dockerfile --tag $(YC_CONTAINER_REGISTRY)/backoffice:stable . && docker push cr.yandex/crpf76jp63emqup99s4l/backoffice:stable

.PHONY: push-event-image
push-event-image:
	DOCKER_DEFAULT_PLATFORM=linux/amd64	docker build -f ./event/Dockerfile --tag $(YC_CONTAINER_REGISTRY)/event:stable . && docker push cr.yandex/crpf76jp63emqup99s4l/event:stable

.PHONY: push-chat-image
push-chat-image:
	DOCKER_DEFAULT_PLATFORM=linux/amd64	docker build -f ./chat/Dockerfile --tag $(YC_CONTAINER_REGISTRY)/chat:stable . && docker push cr.yandex/crpf76jp63emqup99s4l/chat:stable

.PHONY: push-recommendation-image
push-recommendation-image:
	DOCKER_DEFAULT_PLATFORM=linux/amd64	docker build -f ./recommendation/Dockerfile --tag $(YC_CONTAINER_REGISTRY)/recommendation:stable . && docker push cr.yandex/crpf76jp63emqup99s4l/recommendation:stable

.PHONY: push-oracle-image
push-oracle-image:
	DOCKER_DEFAULT_PLATFORM=linux/amd64	docker build -f ./oracle/Dockerfile --tag $(YC_CONTAINER_REGISTRY)/oracle:stable . && docker push cr.yandex/crpf76jp63emqup99s4l/oracle:stable

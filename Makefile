IMAGE_NAME=di-notebook
GIT_SHA = $(shell git rev-parse HEAD)
IMAGE_REPO=jasonblanchard/${IMAGE_NAME}
LOCAL_TAG = ${IMAGE_REPO}
LATEST_TAG= ${IMAGE_REPO}:latest
SHA_TAG = ${IMAGE_REPO}:${GIT_SHA}

createdb:
	# createuser -e -d -P -E di
	createdb -U di -e -O di di_notebook

dropdb:
	dropdb di_notebook

dbmigrate:
	USER=di PASSWORD=di HOST=localhost PORT=5432 DATABASE=di_velocity go run cmd/db/main.go up

dbmigratedown:
	USER=di PASSWORD=di HOST=localhost PORT=5432 DATABASE=di_velocity go run cmd/db/main.go down

migration:
	migrate create -ext sql -dir cmd/db/migrations -seq $$SEQ

docker-build:
	docker build -t ${LOCAL_TAG} .

docker-tag: build
	docker tag ${LOCAL_TAG} ${SHA_TAG}

docker-push: tag
	docker push ${LATEST_TAG}
	docker push ${SHA_TAG}

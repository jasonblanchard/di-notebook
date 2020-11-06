BUILDER=heroku/buildpacks:18
IMAGE_NAME=di-notebook
IMAGE_REPO=jasonblanchard/${IMAGE_NAME}
LOCAL_TAG=${IMAGE_REPO}
LATEST_TAG=${IMAGE_REPO}:latest
VERSION=$(shell git rev-parse HEAD)
VERSION_TAG=${IMAGE_REPO}:${VERSION}

createdb:
	# createuser -e -d -P -E di
	createdb -U di -e -O di di_notebook

dropdb:
	dropdb di_notebook

dbmigrate:
	go run cmd/cli/main.go db migrate --config config/local.yaml

dbmigratedown:
	go run cmd/cli/main.go db migrate -d --config config/local.yaml

migration:
	migrate create -ext sql -dir cmd/db/migrations -seq $$SEQ

build:
	pack build ${IMAGE_REPO} --builder ${BUILDER}
	docker tag ${IMAGE_REPO} ${VERSION_TAG}

push: build
	docker push ${LATEST_TAG}
	docker push ${VERSION_TAG}

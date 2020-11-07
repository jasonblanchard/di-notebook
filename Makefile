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
	migrate -source file://migrations -database postgres://di:di@localhost:5432/di_notebook?sslmode=disable up

dbmigratedown:
	migrate -source file://migrations -database postgres://di:di@localhost:5432/di_notebook?sslmode=disable down

migration:
	migrate create -ext sql -dir cmd/db/migrations -seq $$SEQ

build:
	# pack build ${IMAGE_REPO} --builder ${BUILDER}
	docker build -t ${IMAGE_REPO} .
	docker tag ${IMAGE_REPO} ${VERSION_TAG}

push: build
	docker push ${LATEST_TAG}
	docker push ${VERSION_TAG}

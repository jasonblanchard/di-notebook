BUILDER=heroku/buildpacks:18
IMAGE_NAME=di-notebook
IMAGE_REPO=jasonblanchard/${IMAGE_NAME}
LOCAL_TAG=${IMAGE_REPO}
LATEST_TAG=${IMAGE_REPO}:latest
VERSION=$(shell git rev-parse HEAD)
VERSION_TAG=${IMAGE_REPO}:${VERSION}

db:
	docker rm postgres
	docker run --name postgres -p 5432:5432 -e POSTGRES_PASSWORD=testpass -d postgres:11
	docker exec -it -e PGPASSWORD=testpass -e PGUSER=postgres postgres createuser -e -d -P -E di
	docker exec -it -e PGPASSWORD=testpass -e PGUSER=postgres postgres createdb -U di -e -O di di_notebook

dbclean:
	docker kill postgres && docker rm postgres

dbconnect:
	docker run --rm -it --network=host -e PGPASSWORD=testpass -e PGUSER=postgres -e PGHOST=localhost postgres:11 /bin/bash

dbmigrate:
	docker run -v $(shell pwd):/di --network host migrate/migrate --source file:///di/migrations -database postgres://di:di@localhost:5432/di_notebook?sslmode=disable up

dbmigratedown:
	docker run -v $(shell pwd):/di --network host -it migrate/migrate -source file://di/migrations -database postgres://di:di@localhost:5432/di_notebook?sslmode=disable down

dbdrop:
	docker run -v $(shell pwd):/di --network host -it migrate/migrate -source file://di/migrations -database postgres://di:di@localhost:5432/di_notebook?sslmode=disable drop

migration:
	migrate create -ext sql -dir cmd/db/migrations -seq $$SEQ

build: swagger
	# pack build ${IMAGE_REPO} --builder ${BUILDER}
	docker build -t ${IMAGE_REPO} .
	docker tag ${IMAGE_REPO} ${VERSION_TAG}

swagger:
	wget https://raw.githubusercontent.com/jasonblanchard/di-apis/main/gen/pb-go/notebook.swagger.json -O cmd/http/notebook.swagger.json

push: build
	docker push ${LATEST_TAG}
	docker push ${VERSION_TAG}

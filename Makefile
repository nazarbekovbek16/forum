.SILENT:

build:
	go build -o forum cmd/app/main.go

run:
	go run ./cmd/app/main.go

dbuild:
	docker image build -t forum-img .

drun:
	docker container run -p 9090:9090 -d --name forum-container forum-img

dstop:
	docker stop forum-container

drm:
	docker rm forum-container

drmi:
	docker rmi forum-img

dclear:
	docker system prune -a
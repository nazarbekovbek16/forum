create:
	docker build -t forum .
run:
	docker run -p 8000:8000 forum
stop:
	docker stop forum
start:
	docker start forum
prune:
	docker system prune -a
check:
	docker image ls
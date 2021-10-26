build:
	docker build -t splitter .
start:
	docker run -d --name splitter -p 8003:3001 splitter
stop:
	docker stop splitter
cleanup:
	docker stop splitter && docker rm splitter && docker rmi splitter

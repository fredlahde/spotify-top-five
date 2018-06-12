
all: build-docker
build:
	bash build.sh

build-docker: build
	docker build -t spotify-top-five .

run-local: build
	./spotify-top-five

run: all
	docker run -it --rm spotify-top-five

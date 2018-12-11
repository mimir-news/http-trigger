NAME = $(shell appv name)
VERSION = $(shell appv version)
IMAGE = $(shell appv image)

build:
	docker build -t $(IMAGE) .

delete:
	docker rmi $(IMAGE)
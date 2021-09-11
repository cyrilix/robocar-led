.PHONY: test docker

DOCKER_IMG = docker.io/cyrilix/robocar-led

rc-led: binary-amd64

test:
	go test -race -tags no_d2xx ./cmd/rc-led ./part ./led

docker:
	docker buildx build . --platform linux/arm/7,linux/arm64,linux/amd64 -t ${DOCKER_IMG} --push


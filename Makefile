.PHONY: docker
docker:
	@rm webook-live || true
	@rm webook || true
	@GOOS=linux GOARCH=arm go build -tags=k8s -o webook .
	@docker rmi -f wureny/webook-live:v0.0.1
	@docker build -t wureny/webook-live:v0.0.1 .
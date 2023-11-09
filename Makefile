run:
	mkdir -p "reports"
	export ALLURE_OUTPUT_FOLDER=results/ && \
	echo "$ALLURE_OUTPUT_FOLDER" && \
	go test -v -mod vendor ./...

docker-allure-run: run
	docker run -d -it -p 5050:5050 -e CHECK_RESULTS_EVERY_SECONDS=3 -e KEEP_HISTORY=1 -v ${PWD}/../:/app/projects frankescobar/allure-docker-service:2.19.0
docker-allure-ui-run: docker-allure-run
	docker run -d -p 5252:5252 -e ALLURE_DOCKER_PUBLIC_API_URL=http://localhost:5050 frankescobar/allure-docker-service-ui
docker-run:
	docker run --rm -it --network=host -v allure-results:/workspace/allure-results $$(docker build -q .)
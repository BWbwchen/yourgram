test :
	docker build --file Dockerfile-ci --tag bwbwchen/yourgram_test .
	-docker-compose --file docker-compose-ci.yml up --abort-on-container-exit --exit-code-from server
	docker-compose --file docker-compose-ci.yml down

service_test:
	go test -v yourgram/authentication/service

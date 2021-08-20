test :
#	docker build --file Dockerfile-ci --tag bwbwchen/yourgram_test .
	-docker-compose up --build --abort-on-container-exit account_service 
	docker-compose down
	-docker-compose up --build --abort-on-container-exit jwt_service 
	docker-compose down

account_test:
	go test -v .

jwt_test:
	go test -v .
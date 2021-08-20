test :
#	docker build --file Dockerfile-ci --tag bwbwchen/yourgram_test .
	-docker-compose --file docker-compose-ci.yml up --build --abort-on-container-exit account_service 
	docker-compose down
	-docker-compose --file docker-compose-ci.yml up --build --abort-on-container-exit jwt_service 
	docker-compose down

account_test:
	-docker-compose --file docker-compose-ci.yml up --build --abort-on-container-exit account_service 
	docker-compose down

jwt_test:
	-docker-compose --file docker-compose-ci.yml up --build --abort-on-container-exit jwt_service 
	docker-compose down

build_run :
	docker-compose --file docker-compose.yml up --build

run :
	docker-compose --file docker-compose.yml up 
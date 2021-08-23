.PHONY: run build_run test account_test jwt_test upload_test

test : account_test jwt_test upload_test

account_test:
	-docker-compose --file account/docker-compose.yml up --build --abort-on-container-exit account_service 
	docker-compose --file account/docker-compose.yml down

jwt_test:
	-docker-compose --file jwt/docker-compose.yml up --build --abort-on-container-exit jwt_service 
	docker-compose --file jwt/docker-compose.yml down

upload_test:
	-docker-compose --file upload/docker-compose.yml up --build --abort-on-container-exit upload_service 
	docker-compose --file upload/docker-compose.yml down

build_run :
	docker-compose up --build

run :
	docker-compose up 
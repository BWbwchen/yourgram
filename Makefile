.PHONY: test account_test jwt_test build_run run

test : account_test jwt_test

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
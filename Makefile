.PHONY: run build_run test account_test jwt_test upload_test

test : account_test jwt_test upload_test

account_test:
	cd account/service && go test -v -cover .

jwt_test:
	cd jwt/service && SECRETKEY=bwbwchen go test -v -cover .

upload_test:
	cd upload/service && go test -v -cover .

build_run :
	docker-compose up --build

run :
	docker-compose up 
FROM golang:1.16 AS build

WORKDIR /app
COPY . .

RUN CGO_ENABLED=0 go build -o /account ./cmd/main.go

CMD [ "go", "test", "-v", "-cover", "yourgram/account/service" ]

FROM alpine:latest AS deploy

WORKDIR /
COPY --from=build /account /
CMD ["/account"]
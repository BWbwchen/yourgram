FROM golang:1.16 AS build

WORKDIR /app
COPY . .

RUN CGO_ENABLED=0 go build -o /jwt 

CMD [ "go", "test", "-v", "yourgram/jwt/service" ]

FROM alpine:latest AS deploy

WORKDIR /
COPY --from=build /jwt /
CMD ["/jwt"]
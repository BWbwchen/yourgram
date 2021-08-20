# Yourgram
##  environment variable
```
export DB_URL=localhost:27017
export CONSUL_URL=localhost:8500
export PORT=8080
export SECRETKEY=bwbwchen
```

## run binary on scratch docker image
```
CGO_ENABLED=0 GOOS=linux go build -a -ldflags '-extldflags "-static"' -o howie_tcp .
```

## Consul
```
docker run -d --name=sd -p 8500:8500 consul agent -server -bootstrap -ui -client 0.0.0.0
```
# Yourgram
##  environment variable
```
export REDIS_URL=localhost:6379
export DB_URL=localhost:27017
export ZOOKEEPER_URL=localhost:2181
export PORT=8080
export SECRETKEY=bwbwchen
```

## run binary on scratch docker image
```
CGO_ENABLED=0 GOOS=linux go build -a -ldflags '-extldflags "-static"' -o howie_tcp .
```

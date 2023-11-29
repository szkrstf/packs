# packs

A http server for calulating packages for orders.
You can see a running instance on [https://packs.fly.dev](https://packs.fly.dev).

## Build
```
make build
```

This will run all tests and build the application. (Alternatively you can use `make build-linux` to build a linux specific executable.)

## Test
```
make test
```

## Run
```
./packs -addr ":8080" -config "sizes"
```

The ui is runninng on [http://localhost:8080/](http://localhost:8080/)

## Build and run with Docker
A Dockerfile is provided to build and run the application in a container.
### Build the image
```
docker build -t packs .
```
### Run the container:
```
docker run -p 8080:8080 packs
```
You can also use the `-d` flag to run it in detached mode.
```
docker run -d -p 8080:8080 packs
```

## Config
The sizes can be configured with a file. The location of the config file can be specified with the `-config` flag. Example:
```
250
500
1000
2000
5000
```

## API

### Calculate
Calculate package sizes for an order.

#### Request
```
curl -i -d '{"items":501}' http://localhost:8080/api/calculate
```

#### Response
```
HTTP/1.1 200 OK
Content-Type: application/json
Date: Thu, 16 Nov 2023 13:11:30 GMT
Content-Length: 34

{"data":{"250":1,"500":1}}
```


### Get sizes

#### Request
```
curl -i http://localhost:8080/api/sizes
```

#### Response
```
HTTP/1.1 200 OK
Content-Type: application/json
Date: Thu, 16 Nov 2023 13:11:30 GMT
Content-Length: 27

{"data":[250,500,1000,2000,5000]}
```


### Save sizes

#### Request
```
curl -i -d '{"sizes":[250,500,1000,2000,5000]}' http://localhost:8080/api/sizes
```

#### Response
```
HTTP/1.1 200 OK
Content-Type: application/json
Date: Thu, 16 Nov 2023 13:11:30 GMT
Content-Length: 0

```
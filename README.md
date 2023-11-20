# packs

A http server for calulating packages for orders.

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
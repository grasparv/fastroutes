I hope the code is easy to read. I'll just note some highlights:

* I tried to avoid leaking any coordinates into error message, simulating that
  this is user-sensitive data.
* I annotated structs that are marshalled/unmarshalled mainly as a hint to the
  reader.
* The OSRM server is an old HTTP/1.0 server with keep-alive extension.
  Therefore, I manually set the `Connection` header to `keep-alive`/`close` to play
  nice with the server. It could be unnecessary to close the connection if more
  requests are about to come in, but in this limited example I did not find it
  relevant.
* There is really no scaling/validation on `GetRoutes()`, i.e. you can send in
  millions of destinations and it will just be passed to the OSRM server.
* The `GetRoutes()` call is using a hard-coded context timeout, this should perhaps
  be configurable.
* The `webservice` package calls `GetRoutes()` with the background context, this
  should probably be chained with the web service incoming context.
* There is no caching of responses from downstream calls. This could speedup
  repeated calls for the same source/destination coordinates.
* It seems like the OSRM did not support any kind of batch processing of
  source/destination pairs, so each destination will yield one more HTTP call.
* I kept the web service API extremely simple. There is no logging, no
  validation of input size, no gorilla mux HTTP routing, no port configuration,
  very sparse feedback to the client in case of error, just a HTTP 400/500, no
  check on final write, etc.

# How to run

You can build/run from source or from docker.

## Build from source

This requires at least go 1.22.1.

Start the server:
```bash
$ go run ./cmd/fastroutes
```

## Build with docker

```bash
$ docker build -t jonas/fastroutes:1.0 .
...
$ docker run -d --name fastroutes --rm -p 8080:8080 jonas/fastroutes:1.0
$
```

## Test run

Send a request: 
```bash
$ curl --verbose 'http://localhost:8080/routes?src=13.388860,52.517037&dst=13.397634,52.529407&&dst=13.428555,52.523219'
*   Trying 127.0.0.1:8080...
* Connected to localhost (127.0.0.1) port 8080 (#0)
> GET /routes?src=13.388860,52.517037&dst=13.397634,52.529407&&dst=13.428555,52.523219 HTTP/1.1
> Host: localhost:8080
> User-Agent: curl/7.88.1
> Accept: */*
>
< HTTP/1.1 200 OK
< Content-Type: application/json
< Date: Fri, 15 Mar 2024 09:52:10 GMT
< Content-Length: 189
<
* Connection #0 to host localhost left intact
{"source":"13.388860,52.517037","routes":[{"destination":"13.397634,52.529407","distance":1886.8,"duration":260.3},{"destination":"13.428555,52.523219","distance":3804.2,"duration":389.3}]}
```

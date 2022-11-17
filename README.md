# Rate limiter


## Technical choices

The rate limit strategy I'm using in here is based on Sliding Window rate limitng algorithm.

For this task I'm using in-memory storage to store the users'
id and time window as a means for caching, but this could've been 
replaced by using a cache/key-value storage solution, such as Redis or Cassandra.
That would increase latency a little, but it's a good compromise to
offer precise rate limiting per requester.
To facilitate this implementation, I'm using a library that can manage
the cache storage: github.com/patrickmn/go-cache

The tests should cover the requirements, which is to allow a requester
to send 100 requests per hour, but to get the "429" http you should start a server to handle
the http responses properly. You can also customise the amount of request and window times.
Everything is documented below.

There are many scenarios I couldn't cover given the time, but I would've liked to
have improved test coverage, performance and more time in structuring logs and responses better.
Also adding a Dockerfile with Redis + the Rate Limiter talking to each other.
Thank you for the opportunity and would love to hear some feedback.

## Running this project

Set this environment variables to customise the limits of the rate limiter if you want to
```bash
# examples
export MAX_REQUESTS=10
export WINDOW_DURATION=10m
```

There's a `Makefile` in this project, so you can basically run the following:

```bash
## Install all the dependencies
make install

## running all tests
make test 

## building and running the binary for MACOS
make build_and_run

## clean up
make clean

```

### Testing
When you run `make build` and `make run` you are going to have an HTTP server listening on 8080.
Then, just run a simple call to get a response: 
```bash
curl -v -X GET http://localhost:8080/svc
```
You can reduce the amount of request(`MAX_REQUESTS`) to very low to start getting `429: Too many requests`

Or you can simply run the tests:
```bash
go test .
```
The function responsible for testing the 100 request / hour is: `TestRefusedIsRequestAllowed()`

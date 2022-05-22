# Rate limiter for Airtasker


## Technical choices

The rate limit strategy I'm using in here is based on Sliding Window rate limitng algorithm.

For this task I'm using in-memory storage to store the users'
id and time window as a means for caching, but this could've been 
replaced by using a cache/key-value storage solution, such as Redis or Cassandra.
That would increase latency a little, but it's a good compromise to
offer precise rate limiting per requester.
To facilitate this implementation, I'm using a library that can manage
the cache storage: github.com/patrickmn/go-cache

The tests should cover your requirements, whereas is required to allow a requester
to send 100 requests per hour, but to get the "429" http you should start a server to handle
the http responses properly. Everything is documented below.

There are many scenarios I couldn't cover given the time but I would've liked to
have improved test coverage, performance and more time in structuring logs and responses better.
Also adding a Dockerfile with Redis + this task talking to each other.
Thank you for the opportunity and would love to hear some feedback.

## Running this project

There's a `Makefile` in this project, so you can basically run the following:

```bash
## running the default tests
## 

```

go get github.com/patrickmn/go-cache
go get github.com/gorilla/mux
go mod tidy
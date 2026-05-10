# Tests

This folder contains black-box tests for the distributed cache project.

Run everything:

```bash
go test ./...
```

Run only this test package:

```bash
go test ./tests
```

Run one test group:

```bash
go test ./tests -run TestCache
go test ./tests -run TestHashRing
go test ./tests -run TestNode
go test ./tests -run TestRouter
```

Run with verbose output:

```bash
go test -v ./tests
```

Docker smoke test:

```bash
docker compose up --build
```

In another terminal:

```bash
curl "http://localhost:9000/set?key=user:1&value=Yash"
curl "http://localhost:9000/get?key=user:1"
curl "http://localhost:9000/delete?key=user:1"
curl -i "http://localhost:9000/get?key=user:1"
```

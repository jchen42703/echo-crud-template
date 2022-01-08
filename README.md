# echo-crud-template

[WIP] Basic Echo CRUD template (no pagination)

# Overview

Based on https://github.com/xesina/golang-echo-realworld-example-app.

Echo CRUD Template with:

- [ ] HTTPS & HTTP 2
- [ ] Grouped Paths
- [ ] Basic Authentication
  - [ ] Postgres
  - [ ] Redis
- [ ] Metrics with Prometheus
- [ ] CORS
- [ ] Deploy with Docker
- [ ] Automatic documentation with swagger
- [ ] Testing
  - [ ] https://github.com/alicebob/miniredis

# Run

With a .env file:

```
env $(cat .env | xargs) go run main.go
```

# echo-crud-template

[WIP] Basic Echo CRUD template (no pagination)

# Overview

Based on https://github.com/xesina/golang-echo-realworld-example-app.

Echo CRUD Template with:

- [ ] HTTPS & HTTP 2
- [x] Grouped Paths
- [ ] Basic Authentication
  - [ ] Email Verification
  - [x] Remember Me (Saves Username)
  - [x] Postgres
  - [x] Redis
- [ ] Metrics with Prometheus
- [x] CORS
- [x] CSRF Protection
- [ ] Deploy with Docker
- [ ] Automatic documentation with swagger
- [ ] Testing
  - [ ] https://github.com/alicebob/miniredis

# Run

With a .env file:

```
env $(cat .env | xargs) go run main.go
```

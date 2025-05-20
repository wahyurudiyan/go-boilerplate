# Golang Boilerplate for Backend [under development]

Go Boilerplate is a project template for golang microservices project. It has sample to use this boilerplate for HTTP User Server application. The directory tree for this boilerplate will found below:

```
go-boilerplate
├── api
│   ├── grpc
│   │   └── contract.proto
│   └── rest
│       ├── controller
│       │   ├── controllers.go
│       │   ├── health_check.go
│       │   └── singnup_controller.go
│       └── routes
│           └── routes.go
├── app
│   └── rest.go
├── cmd
│   ├── rest
│   │   ├── gin.go
│   │   └── server.go
│   └── worker
│       └── README.md
├── config
│   └── config.go
├── core
│   ├── dto
│   │   └── user
│   │       ├── login.go
│   │       ├── signup.go
│   │       └── user.go
│   ├── entities
│   │   └── user
│   │       └── user.go
│   ├── repositories
│   │   └── user
│   │       ├── user.database.go
│   │       ├── user.iface.go
│   │       └── user.mongo.go
│   └── services
│       └── user
│           ├── user.iface.go
│           ├── user.impl.go
│           ├── user_signup.impl.go
│           └── user_test.go
├── docs
│   ├── docs.go
│   ├── swagger.json
│   └── swagger.yaml
├── go.mod
├── go.sum
├── LICENSE
├── main.go
├── makefile
├── migrations
│   └── 0001-create_user_table.sql
├── pkg
│   ├── common
│   │   └── http_helper.go
│   ├── configz
│   │   ├── config.go
│   │   └── config_test.go
│   ├── graceful
│   │   ├── graceful.go
│   │   └── graceful_test.go
│   ├── middleware
│   ├── redis
│   │   ├── config.go
│   │   └── connection.go
│   ├── semconv
│   ├── sql
│   │   ├── config.go
│   │   └── connection.go
│   └── telemetry
│       ├── meter.go
│       ├── propagation.go
│       ├── telemetry.go
│       └── tracer.go
└── README.md

32 directories, 45 file
```

## Clean Code Approach

TBD

## TODO - Golang Boilerplate

This is checklist for the next features. For more detail will be discussed withing roadmap.

- [X] COMMON CONFIGURATION
- [X] PACKAGES AS LIBRARY IN `pkg` FOR COMMON LIBRARY THAT USUALLY USED
- [X] GRACEFUL SHUTDOWN BLOCKING WITH CHANNEL
- [X] SWAGGER FOR BETTER DOCUMENTATION
- [X] BASIC PROJECT FOR CRUD API
- [ ] BASIC OPERATION CRUD FOR GRPC
- [ ] EXTERNAL, DIRECTORY CODES FOR ACCESS OTHER RESOURCE WITH HTTP OR GRPC
- [ ] WORKER INTERFACE WITH NATS JETSTREAM AND KAFKA
- [ ] WORKER HANDLER
- [X] OPENTELEMETRY (TRACE, METER AND RUNTIME MONITORING DATA)
- [X] OBSERVABILITY MIDDLEWARE
- [ ] AUTH MIDDLEWARE
- [ ] HTTP WITH OTELHTTP FOR EXTERNAL CALL
- [ ] GRPC WITH OTELGRPC FOR EXTERNAL CALL
- [ ] SEMANTIC CONVENTIONS FOR ERROR CODE, ERROR MESSAGE, ETC.
- [ ] HANDLING CONCURRENCY GROUP SQL-SELECT OR HTTP-GET WITH SINGLEFLIGHT 
- [ ] WIKI DOCS: [Wiki](https://github.com/wahyurudiyan/go-boilerplate/wiki)
- [ ] DOCKERFILE
- [ ] KUBERNETES MANIFEST


## Roadmap

TBD

## Contribution

For contribution please open issue within this repository.

# kelvin-rest-api

## Description

Tech Stack:

- Golang
- Echo Framework
- Go Clean Code Architecture (ref: https://github.com/bxcodec/go-clean-arch)
- Viper
- Postgres

Rule of Clean Architecture by Uncle Bob

- Independent of Frameworks. The architecture does not depend on the existence
  of some library of feature laden software. This allows you to use such
  frameworks as tools, rather than having to cram your system into their limited
  constraints.
- Testable. The business rules can be tested without the UI, Database, Web
  Server, or any other external element.
- Independent of UI. The UI can change easily, without changing the rest of the
  system. A Web UI could be replaced with a console UI, for example, without
  changing the business rules.
- Independent of Database. You can swap out Oracle or SQL Server, for Mongo,
  BigTable, CouchDB, or something else. Your business rules are not bound to the
  database.
- Independent of any external agency. In fact your business rules simply don’t
  know anything at all about the outside world.

More at
https://8thlight.com/blog/uncle-bob/2012/08/13/the-clean-architecture.html

This project has 4 Domain layer :

- Models Layer
- Repository Layer
- Usecase Layer
- Delivery Layer

#### The diagram:

![golang clean architecture](https://github.com/bxcodec/go-clean-arch/raw/master/clean-arch.png)

### How To Run This Project

> Make Sure you already have goland and postgres & run command below
```bash
$ export ENV=DEV 
$ export PORT=9090 
```

Since the project already use Go Module, I recommend to put the source code in
any folder but GOPATH.

#move to directory $ cd kelvin-rest-api

# Clone into your workspace

```bash
$ git clone https://github.com/Kelvinmijaya/kelvin-rest-api.git
```

#move to project

```bash
$ cd kelvin-rest-api
```

# Run the application

```bash
$ go run app/server.go
```

# Execute the call in another terminal

```bash
$ curl localhost:9090/
```

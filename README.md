# Project auth-demo 

One Paragraph of project description goes here

## Getting Started

These instructions will get you a copy of the project up and running on your local machine for development and testing purposes. See deployment for notes on how to deploy the project on a live system.

## Environment Variables

Create a `.env` file in the root directory with the following structure

```env
PORT=3000
APP_ENV=local
DEBUG=true

DB_PORT=5432
DB_NAME=auth-demo
DB_HOST=localhost
DB_USERNAME=your_username
DB_PASSWORD=your_password

JWT_SECRET=your_secret_key
```

## MakeFile

run all make commands with clean tests
```bash
make all build
```

build the application
```bash
make build
```

run the application
```bash
make run
```

Create DB container
```bash
make docker-run
```

Shutdown DB container
```bash
make docker-down
```

run the test suite
```bash
make test
```

clean up binary from the last build
```bash
make clean
```


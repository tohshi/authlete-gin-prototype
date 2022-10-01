# Authlete-Gin-Prototype

Prototype of an authorization server using Gin and Authlete.  
This implementation is not perfect and should be kept for reference only.  
Redis is used for session management.

## Getting Started

1. Install dependencies

```sh
$ go mod tidy
```

2. Set up environment variables

```sh
$ vi .env
```

3. Start Redis

```sh
$ docker-compose up -d
```

4. Start server on `http://localhost:8080`

```sh
$ go run *.go
```

### Endpoints

| Endpoint               | Path                  |
| :--------------------- | :-------------------- |
| Authorization Endpoint | `/auth/authorization` |
| Token Endpoint         | `/auth/token`         |

### Users

| Login ID | Password  | Consent Required |
| :------- | :-------- | :--------------- |
| `user1`  | `passwd2` | `true`           |
| `user2`  | `passwd2` | `false`          |

# JWT Server with Gin
## Getting Started
1. Start Postgres
    ```
    $ docker-compose up
    ```

1. Run server
    ```
    $ go run main.go
    ```

## API
### POST /signup
```sh
curl --location --request POST 'localhost:8888/signup' \
--form 'email=test@example.com' \
--form 'password=abcd123' \
--form 'username=test'
```

### POST /login
```sh
curl --location --request POST 'localhost:8888/login' \
--form 'email=test@example.com' \
--form 'password=abcd123'
```

### POST /v1/me
```sh
curl --location --request GET 'localhost:8888/v1/me' \
--header 'Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJlbWFpbCI6InRlc3Q2QHNhbXBsZS5jb20iLCJhY2NvdW50X2lkIjo4LCJleHAiOjE2MTUxMjA3Mjh9.mgYfZVWZ_Uec5GBtWE02n2R5v-Air_A5mw2uKW-4tVA'
```

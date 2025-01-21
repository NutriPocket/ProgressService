# ProgressService
Simple auth REST-like API in go using JWT, Gin and MySQL

Routes:

    - /auth
      - POST /register
      - POST /login
      - POST /logout
    - /users
      - GET /
      - GET /:username

Build & Run

```
docker-compose up api
```


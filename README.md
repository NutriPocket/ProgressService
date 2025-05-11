# ProgressService

Routes:

    - Anthropometric data
      - PUT /users/:id/anthropometrics
      - GET /users/:id/anthropometrics
        - Params:
          - date: "YYYY-MM-DD" string format date
    - Fixed user data
      - PUT /users/:id/fixedData
      - GET /users/:id/fixedData
        - Params:
          - base: true/false, get base fixed data (birthday instead of age)
    - Objectives
      - PUT /users/:id/objectives
      - GET /users/:id/objectives

Build & Run

```
docker-compose up --build
```

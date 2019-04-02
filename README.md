## Overview
This Project implements the backend apis using gin framework.
## Local Docker Deployment
Create a docker-compose configuration, and use `./docker-compose.yml` as the compose file.

## API Documentation
### Post /create-new-user
Parameters:
```
{
  firstName: Tom,
  lastName: Hanks,
  email: tomHanks@example.com,
  password: password
}
```
### Get /get-all-users
Expected Result:
```
{
  [
    firstName: Tom,
    lastName: Hanks,
    email: tomHanks@example.com
]
  }
```


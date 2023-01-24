# API requests for kitsu APIs

## User

### Create User

- **Method**
  - `POST`

- **Endpoint**
  - `https://kitsu.io/api/edge/users`

- **Request Header**
  - `Content-Type:application/vnd.api+json`

- **Request Body**

``` JSON
{
    "data": {
        "attributes": {
            "email": <userEmail>,
            "name": <userName>,
            "password": <password>
        },
        "type": "users"
    }
}
```

### Read User

- **Method**
  - `GET`

- **Endpoint**
  - `https://kitsu.io/api/edge/users/<id>`

- **Request Header**
  - `Content-Type:application/vnd.api+json`

### Delete User

- **Method**
  - `DELETE`

- **Endpoint**
  - `https://kitsu.io/api/edge/users/<id>`

- **Request Header**
  - `Content-Type:application/vnd.api+json`
  - `Authorization:Bearer <access_token>`

## Authentication

### Get Token

- **Method**
  - `POST`

- **Endpoint**
  - `https://kitsu.io/api/oauth/token`

- **Request Header**
  - `Content-Type:application/vnd.api+json`

- **Request Body**

``` JSON
{
  "grant_type": "password",
  "username": <email>,
  "password": <password>
}
```

## User Library

### Create Entry

- **Method**
  - `POST`

- **Endpoint**
  - `https://kitsu.io/api/edge/library-entries`

- **Request Body**

``` JSON
{
    "data": {
        "attributes": {
            "status": <sate of the entity i.e. completed, started, want to watch>
        },
        "relationships": {
            // only one of the two objects needed
            "anime":{
                "data": {
                    "id": <id of the anime>,
                    "type": "anime"
                }
            },
            "manga":{
                "data": {
                    "id": <id of the manga>,
                    "type": "manga"
                }
            },
            "user":{
                "data": {
                    "id": <id of the user>
                    "type": "users"
                }
            }
        },
        "type": "library-entries"
    }
}
```

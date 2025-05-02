# chirpy
A RESTful API that gives the ability to manage profiles and Chirps.

## Why Chirpy?
Chirpy is a guided project from boot.dev with the main goal of teaching me how HTTP servers work including how to authenticate and authorize, in general how to build a RESTful API from the ground up. 

While this is a relatively simple project, it has taught me a lot about the many things that go into HTTP servers, feel free to look at the code to see the general structure of everything, I'm sure I made mistakes here and there and likely weird naming choices, feel free to take this code and spin it as your own thing by editing it or adding onto it!

## Installation
This project expects you to have the latest version of Go, Goose, PostgreSQL.
That's about it, set up a Chirpy database with PostgreSQL, run "goose {connection_string} up" in sql/schema to let the SQL code build everything up in the database, and you're good to go.

## How to use
Since this is a RESTful API it'll follow standard conventions with GET, POST, PUT, DELETE methods being available. Below is a list of the endpoints, what you should give them, and what to expect back.

**Note, you should include a .env file that includes a DB_URL, PLATFORM, JWTSECRET, and POLKA_KEY, these will be needed to allow the code to work, authenticate, and use its webhook endpoint.**

### User endpoints
1. POST /api/users

**Give**
```
{
    "email": test@test.com,
    "password": Password123
}
```
**Receive**
```
{
    "id": 123456789,
    "created_at": 2025-05-01 12:34:56,
    "updated_at": 2025-05-01 12:34:56,
    "email": test@test.com,
    "is_chirpy_red": false
}
```

2. POST /api/login

**Give**
```
{
    "email": test@test.com,
    "password": Password123
}
```
**Receive**
```
{
    "id": 123456789,
    "created_at": 2025-05-01 12:34:56,
    "updated_at": 2025-05-01 12:34:56,
    "email": test@test.com,
    "is_chirpy_red": false,
    "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
    "refresh_token": "d1efb07e-eabc-467f-b8ee-fd93b7d88be2"
}
```

3. PUT /api/users

**Give**
*Headers:*
Authorization: Bearer ${AccessToken}
```
{
    "email": test@test.com,
    "password": Password123
}
```
**Receive**
```
{
    "email": test@test.com
}
```

### Chirp Endpoints
1. POST /api/chirps
Authorization: Bearer ${AccessToken}

**Give**
```
{
    "body": Chirpy rocks!
}
```
**Receive**
```
{
    "id": 123456789,
    "created_at": 2025-05-01 12:34:56,
    "updated_at": 2025-05-01 12:34:56,
    "body": Chirpy rocks!,
    "user_id": 123456789
}
```

2. GET /api/chirps

**Give**
*Can query by author_id and sort, defaults to ascending order by created_at eg. GET /api/chirps?sort=desc*
**Receive**
```
[
{
    "id": 123456789,
    "created_at": 2025-05-01 12:34:56,
    "updated_at": 2025-05-01 12:34:56,
    "body": Chirpy rocks!,
    "user_id": 123456789
}
]
```

3. GET /api/chirps/{chirpID}

**Give**
*Can query by chirpID to get a specific Chirp by its ID*
**Receive**
```
{
    "id": 123456789,
    "created_at": 2025-05-01 12:34:56,
    "updated_at": 2025-05-01 12:34:56,
    "body": Chirpy rocks!,
    "user_id": 123456789
}
```

4. DELETE /api/chirps/{chirpID}

**Give**
Authorization: Bearer ${AccessToken}
*Can query by chirpIP to delete the specified Chirp*
**Receive**
Just a 204 status code

### Admin Endpoints
There are also "POST /admin/reset" and "GET /admin/metrics" endpoints with one deleting everything in the database for a clean slate and the other returning how many hits the API has gotten respectively, they're pretty self explanator, just call them and it should work, since this is all local there's not much security to these.

## Conclusion
As you can see, this is a pretty simple API, I learned a ton from doing this and I hope you enjoy playing around with it. Feel free to contribute by forking the repo and opening pull requests, all pull requests should be submitted to the main branch.

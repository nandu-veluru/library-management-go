# Add Author(TASK 1)
- path: internal/author

## Implementations

## Author
- Id(Primary key), Name, Biography, Year of Birth, Nationality

## Create Author Request
- Name (required) (After trimming spaces the charaters should be in between 1 to 160)
- Biography (optional) (when present, is at most 4,000 characters.)
- Year of birth (optional) (when present, is from 1 through the current calendar year)
- Nationality (optional) (when present, is at most 80 characters)

## Service
- Register Author With Valid details

## Testing
- Tested upon successful registration, duplicate name, Invalid Author details

## Database
- insert doctor details
- check author existing by name

## API Request

**Implemented URL** 

 (POST)http://localhost:8080/api/v1/author 

**Login**

- before adding author, librarian or admin should login (A member, and anyone who is not signed in, is refused)

```bash
TOKEN=$(curl -s -X POST http://localhost:8080/api/v1/auth/login \
  -H 'Content-Type: application/json' \
  -d '{"email":"admin@library.local","password":"ChangeMe123!"}' \
  | jq -r .accessToken)
```

**Register Author**

```bash
curl -s -X POST http://localhost:8080/api/v1/author \
  -H 'Content-Type: application/json' \
  -d '{
    "name": "mitch albom",
    "biography": "I am author of book called Tuesdays with morrie",
    "birth": 1958,
    "nationality": "American"
}'
```

## API Response

**Successful Author Register**

```bash
{
    "name": "oliver",
    "biography": "I am author of book called alice in wonderland",
    "birth": 1958,
    "nationality": "American"
}
```

**Adding Author without Login**

```bash
{
    "timestamp": "2026-11-13T11:04:16.540416815Z",
    "status": 401,
    "error": "Unauthorized",
    "message": "Authentication is required",
    "path": "/api/v1/author"
}
```

**Adding Duplicate Author**

```bash
{
    "timestamp": "2026-11-13T11:06:50.747485263Z",
    "status": 409,
    "error": "Conflict",
    "message": "Author already exists",
    "path": "/api/v1/author"
}
```
# View Author(TASK 2)

- admin, librarian and member can view the author details through id 

**Implemented URL**

(GET)http://localhost:8080/api/v1/authors?id

- before Viewing author, librarian or admin should login (A member and Unauthorized Members are refused to view authors)

**View Author who is existing**

```bash
{
    "id": "",
    "name": "Toshikazu Kawaguchi",
    "biography": "I am author of book called Before coffee gets cold",
    "birth": 1958,
    "nationality": "American",
    "created_at": "0001-01-01T00:00:00Z",
    "updated_at": "0001-01-01T00:00:00Z"
}
```

**View Author without Login**

```bash
{
    "timestamp": "2026-11-13T11:04:16.540416815Z",
    "status": 401,
    "error": "Unauthorized",
    "message": "Authentication is required",
    "path": "/api/v1/author"
}
```

**View author who is not existing**

```bash
{
    "timestamp": "2026-11-18T04:25:13.454533674Z",
    "status": 500,
    "error": "Internal Server Error",
    "message": "An unexpected error occurred",
    "path": "/api/v1/authors"
}
```
# search Author(Task 3)

- Admin, librarian and member can search author through name

- each page contains limited authors in ascending order

**Implemented URL**

(GET)http://localhost:8080/api/v1/authors?name=&page=&limit=
- can provide required limit and page number 

**Search the existing Authors**
- A name filter matches authors whose name contains the given text, ignoring letter case

(GET)http://localhost:8080/api/v1/authors?name=pa&page=1&limit=2
LIMIT : 2
PAGE : 1

```bash
{
    "data": [
        {
            "id": "50c595e7-48ac-4ed0-9198-dd37182a5ef6",
            "name": "padama",
            "biography": "I am author",
            "birth": 1958,
            "nationality": "indian",
            "created_at": "2026-11-19T12:30:05.765901Z",
            "updated_at": "2026-11-19T12:30:05.765901Z"
        },
        {
            "id": "d2e6b13b-e1b9-48c9-a623-84c034a1912e",
            "name": "palle vasu",
            "biography": "I am author of book called WE WERE NEVER MEANT TO BE",
            "birth": 1958,
            "nationality": "American",
            "created_at": "2026-11-17T13:46:46.667227Z",
            "updated_at": "2026-11-17T13:46:46.667227Z"
        }
    ],
    "total": 3
}
```
(GET)http://localhost:8080/api/v1/authors?name=pa&page=2&limit=2
LIMIT : 2
PAGE : 2
```bash
{
    "data": [
        {
            "id": "170ba55c-c907-4df8-a637-36fceb0c1108",
            "name": "patrik",
            "biography": "I am author",
            "birth": 1958,
            "nationality": "American",
            "created_at": "2026-11-19T12:28:57.341086Z",
            "updated_at": "2026-11-19T12:28:57.341086Z"
        }
    ],
    "total": 3
}
```

- With no name, limit, page filter, the page lists authors in name order.

```bash
{
    "data": [
        {
            "id": "922227ac-6406-455e-bdef-a0b305a672a1",
            "name": "Mitch Albom",
            "biography": "I am author of book called Tuesdays with morrie",
            "birth": 1958,
            "nationality": "American",
            "created_at": "2026-11-17T06:54:44.475172Z",
            "updated_at": "2026-11-17T06:54:44.475172Z"
        },
        {
            "id": "45e9bc8e-d9c0-4940-bc5d-e5c0b9f8473c",
            "name": "oliver",
            "biography": "I am author of book called alice in wonderland",
            "birth": 1990,
            "nationality": "Mexican",
            "created_at": "2026-11-18T06:53:09.114332Z",
            "updated_at": "2026-11-18T06:53:09.114332Z"
        },
        {
            "id": "50c595e7-48ac-4ed0-9198-dd37182a5ef6",
            "name": "padama",
            "biography": "I am author",
            "birth": 1958,
            "nationality": "indian",
            "created_at": "2026-11-19T12:30:05.765901Z",
            "updated_at": "2026-11-19T12:30:05.765901Z"
        },
        {
            "id": "d2e6b13b-e1b9-48c9-a623-84c034a1912e",
            "name": "palle vasu",
            "biography": "I am author of book called WE WERE NEVER MEANT TO BE",
            "birth": 1958,
            "nationality": "American",
            "created_at": "2026-11-17T13:46:46.667227Z",
            "updated_at": "2026-11-17T13:46:46.667227Z"
        },
        {
            "id": "170ba55c-c907-4df8-a637-36fceb0c1108",
            "name": "patrik",
            "biography": "I am author",
            "birth": 1958,
            "nationality": "American",
            "created_at": "2026-11-19T12:28:57.341086Z",
            "updated_at": "2026-11-19T12:28:57.341086Z"
        },
        {
            "id": "918e8cea-f2ef-4a79-abca-1155768b1d6a",
            "name": "ritivik",
            "biography": "I am author of book called i cant say goodbye to you",
            "birth": 1958,
            "nationality": "American",
            "created_at": "2026-11-18T08:08:48.254054Z",
            "updated_at": "2026-11-18T08:08:48.254054Z"
        },
        {
            "id": "8ab60d2a-8402-4737-a91b-18e5c88a3e8c",
            "name": "Toshikazu Kawaguchi",
            "biography": "I am author of book called Before coffee gets cold",
            "birth": 1958,
            "nationality": "American",
            "created_at": "2026-11-18T04:08:30.906566Z",
            "updated_at": "2026-11-18T04:08:30.906566Z"
        }
    ],
    "total": 7
}
```
- When nothing matches, the result is an empty page.

(GET)http://localhost:8080/api/v1/authors?name=xy&page=1&limit=2

```bash
    {
    "data": [],
    "total": 0
}
```
# Remove Author(Task 5)

- Admin and librarian can remove the author
- while removal Admin and librarian should know who is removed

**Implemented URL**
(DELETE)http://localhost:8080/api/v1/authors/+id

- Author successfully deleted response
```bash
{
    "Author": "ritivik",
    "message": "author deleted sucessfully"
}
```




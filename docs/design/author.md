# Add Author
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





# Task Tracker API

A simple REST API for managing tasks, including user registration, login, and task management.

## Table of Contents
- [Setup Instructions](#setup-instructions)
- [API Documentation](#api-documentation)
  - [POST /register](#post-register)
  - [POST /login](#post-login)
  - [POST /logout](#post-logout)
  - [POST /task](#post-task)
  - [GET /tasks](#get-tasks)
  - [GET /tasks/{user_id}](#get-tasks-user-id)
  - [PUT /task/{task_id}/status](#put-task-task-id-status)
  - [DELETE /task/{task_id}](#delete-task-task-id)

---

## Setup Instructions

### Prerequisites:
- Go 1.18+ installed
- Docker (for PostgreSQL)
- A Postgres database running (via Docker or local installation)

### Steps to run the app locally:

1. Clone the repository:
   ```bash
   git clone <repository-url>
   cd task-tracker-api

2. Build and run the application:

```bash
go mod tidy
go run main.go
```
3. The server should now be running on http://localhost:8080.

4. (Optional) To run PostgreSQL in Docker, use the following command:

```bash
docker run --name postgres -e POSTGRES_PASSWORD=postgres -e POSTGRES_DB=task_tracker -p 5432:5432 -d postgres
```



# API Documentation

## POST /register

**Description:**  
This endpoint is used to register a new user. The user must provide a `username` and `password`.

**Request Body:**
```json
{
  "username": "testuser",
  "password": "testpassword"
}
```
**Response:**

```json

//Success Response
{
    "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3MzgxNjkwNzUsInVzZXJfaWQiOjN9.nRRhGAk_dRhNS9-4ingJmVDVFIOeiWACAEy0rMNzhNI"
}

//error Response
{
    "error": "Key: 'User.Password' Error:Field validation for 'Password' failed on the 'required' tag"
}
```

## POST /login
This endpoint is used to login user. The user must provide a `username` and `password`.

**Request Body:**
```json
{
    "username":"testuser",
    "password":"testuser"
}
```
**Response:**

```json
//Success Response
{
    "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3MzgxNjkyMDcsInVzZXJfaWQiOjB9.JFUhEkouvVv39a7Q5F5hTHYGDTdKtlkcI61JSOKaFCs"
}

//Error Response
{
    "error": "Invalid username or password"
}
```

## POST /logout
Log out the user by invalidating the JWT token.

**Request**
Header: token= "jwt-token

**Response**
```json
//Success Response
{
    "message": "Successfully logged out"
}


//Error Response
{
    "error": "Invalid token"
}
```

Tokens are required for all the Task API calls (passed in headers in all APIs).
```json
// If the token is not provided

{
    "error": "Token is required"
}

// If the token is expired
{
    "error": "Please Login Again..."
}
```
## POST /task



Create a new task.

**Request**
```json
{
  "title": "Write blog post",
  "description": "Draft a blog post about Golang REST APIs",
  "assignee_id":3
}
```

**Response**
```json
//Success Response
{
    "id": 2,
    "title": "Write blog post",
    "description": "Draft a blog post about Golang REST APIs",
    "status": "todo",
    "assignee_id": 3
}

//error Responses
{
    "error": "Task title is required"
}
{
    "error": "Invalid assignee ID"
}

```

## GET /tasks
Retrieve all tasks.

**Response**
```json
[
    {
        "assigned_user": "testuser",
        "assignee_id": 3,
        "created_at": "2025-01-26T22:24:13.630955Z",
        "description": "Draft a blog post about Golang REST APIs",
        "id": 2,
        "status": "todo",
        "title": "Write blog post",
        "updated_at": null
    }
]

```

## GET /tasks/{user_id}

This endpoint retrieves all tasks assigned to a specific user, identified by user_id.


**Response**
```json
[
    {
        "created_at": "2025-01-26T22:24:13.630955Z",
        "description": "Draft a blog post about Golang REST APIs",
        "id": 2,
        "status": "todo",
        "title": "Write blog post",
        "updated_at": null
    }
]
```
## PUT /task/{task_id}/status
This endpoint updates the status of a specific task. The status can be one of the following: todo, in_progress, completed.


**Request**
```json
{
  "status": "completed"
}
```

**Response**

```json
{
    "id": 2,
    "status": "completed",
    "title": "Write blog post"
}

//error response
{
    "error": "Invalid status"
}
```

## DELETE /task/{task_id}
This endpoint deletes a specific task by its task_id.

**Response**
```json
{
    "message": "Task deleted successfully"
}
```
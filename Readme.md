# GShortify

## 📋 <a name="table">Table of Contents</a>

1. ⭐ Introduction
2. 🔨 Tech Stack
3. 📜 Features

## <a name="introduction">⭐ Introduction</a>

GShortify, a URL Shortener application built to demonstrate intermediate CRUD operations that can be performed using Golang. It involves the creation and consumption of APIs, allowing to shorten and manage URLs efficiently, while showcasing key concepts such as data handling, routing, and database interaction.

## <a name="tech-stack">🔨 Tech Stack</a>

- HTML/CSS (Frontend)
- Go (Backend)
- Supabase (Database)

## 📜 Features

👉 **API Creation**: using `gorilla/mux` package to serve different endpoints for URL shortener

👉 **API consumption**: by the UI application to shorten a long URL

## 🔧 API Endpoints

| Method | Endpoint                           | Description                                      |
| ------ | ---------------------------------- | ------------------------------------------------ |
| GET    | `/api/v1/routes/redirect/:shortid` | Redirect to original URL                         |
| GET    | `/api/v1/routes/urls`              | List all URLs and corresponding shorten id       |
| POST   | `/api/v1/routes/shorten`           | Create shorten id and shorten URL                |
| PUT    | `/api/v1/routes/update/:shortid`   | Update original URL                              |
| DELETE | `/api/v1/routes/update/:shortid`   | Update original URL and corresponding shorten id |

<details>

<summary style="font-size: 18px;">List original URL and corresponding shortID</summary>

### `GET` /api/v1/urls

`Request`

- Client's IP Address

`Response`

```go
{
  "code": 200,
  "message": "Data found",
  "data": [
    {
      "tw9apb98": "https://google.com"
    },
    {
      "ngMI98wQ": "https://yahoo.com"
    },
    {
      "O8dYvXGH": "https://netflix.com"
    },
    {
      "d9jYcHDU": "https://zee5.com"
    }
  ]
}
```

</details>

<details>

<summary style="font-size: 18px;">Shorten long URL</summary>

### `POST` /api/v1/shorten

`Request`

```go
{
  "url": "https://linkedin.com"
}
```

`Response`

```go
{
  "code": 201,
  "message": "Shorten url generated successfully",
  "data": {
    "shortened-url": [
      {
        "localhost:8080/z5VAZ6bN"
      }
    ]
  }
}
```

</details>

<details>

<summary style="font-size: 18px;">Update original URL</summary>

### `PUT` /api/v1/update/:shortid

`Request`

```go
{
  "url": "https://google.co.in"
}
```

`Response`

```go
{
  "code": 200,
  "message": "Data updated successfully",
  "data": [
    {
    "previous-url": "https://google.com",
    "original-url": "https://google.co.in",
    "short-id": "tw9apb98"
    }
  ]
}
```

</details>

<details>

<summary style="font-size: 18px;">Delete original URL and corresponding shortID</summary>

### `DELETE` /api/v1/delete/:shortid

`Response`

`204 No Content`

</details>

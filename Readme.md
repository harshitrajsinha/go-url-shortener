<div align="center">
 
  <br />

  <h2 align="center">GShortify (GO Crud Operations)</h2>
</div>

## 📋 <a name="table">Table of Contents</a>

1. ⭐ [Introduction](#introduction)
2. 🔨 [Tech Stack](#tech-stack)
3. 📜 [Features](#features)

## <a name="introduction">⭐ Introduction</a>

GShortify, a URL Shortener application built to demonstrate intermediate CRUD operations that can be performed using Golang. It involves the creation and consumption of APIs, allowing to shorten and manage URLs efficiently, while showcasing key concepts such as data handling, routing, and database interaction.

## <a name="tech-stack">🔨 Tech Stack</a>

- HTML/CSS (Frontend)
- Go (Backend)
- Supabase (Database)

## <a name="features">📜 Features</a>

👉 **API Creation**: using `gorilla/mux` package

- GET : Redirect to original URL based on shortened ID
- POST : Generate a shorten URL
- PUT : Update the original URL based on shortened ID
- DELETE : Delete an existing URL and corresponding shortened ID

👉 **API consumption**: by the UI application

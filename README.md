# 🚀 Enterprise POS Engine (Backend)

![Go](https://img.shields.io/badge/go-%2300ADD8.svg?style=for-the-badge&logo=go&logoColor=white)
![Postgres](https://img.shields.io/badge/postgres-%23316192.svg?style=for-the-badge&logo=postgresql&logoColor=white)
![JWT](https://img.shields.io/badge/JWT-black?style=for-the-badge&logo=JSON%20web%20tokens)

A high-performance, concurrent Point of Sale (POS) backend engine built with **Golang 1.22+**. Designed with Clean Architecture principles, ensuring scalability, maintainability, and enterprise-grade security.

## 🏗️ Architectural Highlights
* **Clean Architecture:** Strict separation of concerns (Models, Handlers, Middleware, Config).
* **High-Performance Routing:** Utilizes `go-chi/chi` for lightweight and blazing-fast HTTP routing.
* **Stateless Security:** Implements JSON Web Tokens (JWT) for secure session management and Bcrypt for irreversible password hashing.
* **Enterprise Database Driver:** Uses `jackc/pgx/v5` for optimized, connection-pooled PostgreSQL transactions.

## ⚙️ Local Development Setup
1. Clone the repository.
2. Create a `.env` file in the root directory:
   ```env
   DATABASE_URL=postgres://user:password@localhost:5432/next_pos_db?sslmode=disable
   PORT=8080
   JWT_SECRET=your_super_secret_key
# ShrimPG 🦐
> **Shrimp**-powered **P**assword **G**ate — secure, fast, and elegant.

<p align="center">
  <img src="assets/photo_2026-03-08_23-04-36.jpg" width="200" alt="ShrimPG Logo">
</p>

[![License: GPL v3](https://img.shields.io/badge/License-GPLv3-blue.svg)](https://www.gnu.org/licenses/gpl-3.0)
[![Go Version](https://img.shields.io/badge/go-1.22+-blue.svg)](https://golang.org)
...

ShrimPG is a secure secrets management system designed with a focus on cryptographic integrity, modularity, and clean architecture.

## Architecture Overview
The project follows a layered architecture, decoupling business logic from storage and security concerns:
* **API Layer**: RESTful handlers with session-based middleware.
* **Logic Layer**: Password validation and user context management.
* **Security Layer**: High-level cryptographic primitives (AES-GCM, Argon2/Scrypt).
* **Storage Layer**: PostgreSQL 16 with volume persistence.

## Security Workflow
- **Zero-Password Storage**: The Master Password is never stored in the database. 
- **Cryptographic Auth**: Authentication is verified by attempting to decrypt a "master_check" record. If decryption fails, the key is invalid.
- **Unique Salting**: Every password entry uses a unique 12-byte salt for key derivation, protecting against rainbow table attacks.
- **Graceful Shutdown**: The server ensures all database transactions are completed and connections are closed properly before exiting.

## Tech Stack
- **Core**: Go (Golang) 1.26+
- **Database**: PostgreSQL 16
- **Infrastructure**: Docker & Docker Compose
- **Auth**: Token-based Session Management

## Getting Started

### Prerequisites
- Docker & Docker Compose

### Installation
1. Clone the repository:
   ``` Bash
   git clone [https://github.com/Krev3tka/ShrimPG.git](https://github.com/Krev3tka/ShrimPG.git)
   cd ShrimPG

2. Start the infrastructure:
    ``` Bash
    docker-compose up --build -d

## API Reference

The server enforces TLS (ListenAndServeTLS), so all endpoints must be accessed via https://.
For local testing with self-signed certs, use curl -k to bypass verification.
1. User Registration (POST /api/v1/register)
    ```Bash
    
    curl -k -X POST [https://127.0.0.1:8080/api/v1/register](https://127.0.0.1:8080/api/v1/register) \
      -H "Content-Type: application/json" \
      -d '{"username": "your_username", "password": "your_master_password"}'

2. User Login (POST /api/v1/login)
    ```Bash
    
    curl -k -X POST [https://127.0.0.1:8080/api/v1/login](https://127.0.0.1:8080/api/v1/login) \
      -H "Content-Type: application/json" \
      -d '{"username": "your_username", "password": "your_master_password"}'
    ```

    Expected Response:
    JSON
    
    {"token": "YOUR_TOKEN"}

3. User Logout (POST /api/v1/logout)
    ```Bash
    
    curl -k -X POST [https://127.0.0.1:8080/api/v1/logout](https://127.0.0.1:8080/api/v1/logout) \
      -H "Authorization: Bearer <YOUR_SESSION_TOKEN>"

4. Create Password Entry (POST /api/v1/passwords/create)
    ```Bash
    
    curl -k -X POST [https://127.0.0.1:8080/api/v1/passwords/create](https://127.0.0.1:8080/api/v1/passwords/create) \
      -H "Content-Type: application/json" \
      -H "Authorization: Bearer <YOUR_TOKEN>" \
      -d '{"service": "github.com", "login": "your_username", "password": "secret_password"}'

5. Get Specific Password (POST or GET /api/v1/passwords/get)

   Note: Service name must be passed in the JSON body (max 40 chars).
   ```Bash
    
   curl -k -X GET [https://127.0.0.1:8080/api/v1/passwords/get](https://127.0.0.1:8080/api/v1/passwords/get) \
     -H "Content-Type: application/json" \
     -H "Authorization: Bearer <YOUR_SESSION_TOKEN>" \
     -d '{"service": "github.com"}'

6. Delete Password Entry (DELETE /api/v1/passwords/delete)
    Note: Service name must be passed in the JSON body.
    ```Bash
    
    curl -k -X DELETE [https://127.0.0.1:8080/api/v1/passwords/delete](https://127.0.0.1:8080/api/v1/passwords/delete) \
      -H "Content-Type: application/json" \
      -H "Authorization: Bearer <YOUR_SESSION_TOKEN>" \
      -d '{"service": "github.com"}'

7. List All Passwords (GET /api/v1/passwords/list)
    ```Bash
    
    curl -k -X GET [https://127.0.0.1:8080/api/v1/passwords/list](https://127.0.0.1:8080/api/v1/passwords/list) \
      -H "Authorization: Bearer <YOUR_TOKEN>"

Roadmap

- [x] PostgreSQL Integration: Docker-ready with volume persistence.

- [x] Session-based Auth: Secure middleware with master-key validation.

- [x] CRUD Core: Fully functional REST API for password management.

- [ ] TUI client, based on BubbleTea framework.

License

Distributed under the GNU GPL v3 License. See LICENSE for more information.

Built with 🦐 passion by Krev3tka

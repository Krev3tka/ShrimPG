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
    docker-compose up -d
    
3. Run the application:
    ``` Bash
    go run cmd/passwordManager/main.go

🗺 Roadmap

    [x] PostgreSQL Integration: Docker-ready with volume persistence.

    [x] Session-based Auth: Secure middleware with master-key validation.

    [x] CRUD Core: Fully functional REST API for password management.

    [x] Rust Integration: Moving encryption logic to a Rust module via gRPC.

    [ ] Desktop Client: Cross-platform GUI (Tauri).

📄 License

Distributed under the GNU GPL v3 License. See LICENSE for more information.

Built with 🦐 passion by Krev3tka

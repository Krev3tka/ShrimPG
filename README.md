# ShrimPG 🦐
> **Shrimp**-powered **P**assword **G**ate — secure, fast, and elegant.

<p align="center">
  <img src="assets/photo_2026-03-08_23-04-36.jpg" width="200" alt="ShrimPG Logo">
</p>

[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)
[![Go Version](https://img.shields.io/badge/go-1.22+-blue.svg)](https://golang.org)
...

ShrimPG is a high-performance secrets management system designed with a strong focus on security, modularity, and horizontal scalability.

## Architecture Overview
The project is built on microservices architecture principles, decoupling business logic between a high-performance Go backend and a memory-safe Rust cryptographic module.

```mermaid
graph LR
    User -->|gRPC/REST| GoApp[Go Backend]
    GoApp -->|gRPC| RustCrypto[Rust Crypto Service]
    GoApp -->|SQL| Postgres[(PostgreSQL)]
```

Tech Stack

    Core: Go (Golang)

    Security Module: Rust

    Communication: gRPC / Protocol Buffers

    Database: PostgreSQL 16

    Infrastructure: Docker & Docker Compose
 
Getting Started

Prerequisites

    Docker & Docker Compose

Installation

Clone the repository:
```Bash
git clone [https://github.com/Krev3tka/ShrimPG.git](https://github.com/Krev3tka/ShrimPG.git)
cd ShrimPG
```

Start the infrastructure:
```Bash
docker-compose up -d
```

Roadmap
- [x] **PostgreSQL Integration:** Docker-ready with volume persistence.
- [x] **Session-based Auth:** Secure middleware with master-key validation.
- [x] **CRUD Core:** Fully functional REST API for password management.
- [ ] **Rust Integration:** Moving encryption logic to the Rust module via gRPC.
- [ ] **Desktop Client:** Cross-platform GUI (Tauri + Rust).

License

Distributed under the MIT License. See LICENSE for more information.

Built with passion by Krev3tka

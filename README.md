# Gochain <!-- omit in toc -->

A full-stack blockchain interoperability demo showcasing a modular **Go backend** (with OpenAPI validation and multi-chain adapters) and a **Next.js frontend** with **Metamask integration**.

## Table of Contents<!-- omit in toc -->

- [Overview](#overview)
- [Architecture](#architecture)
- [Usage](#usage)
  - [Requirements](#requirements)
  - [Setup and Running](#setup-and-running)
    - [1. Clone the repository](#1-clone-the-repository)
    - [2. Configure environment variables](#2-configure-environment-variables)
    - [3. Run the backend](#3-run-the-backend)
    - [4. Run the frontend](#4-run-the-frontend)
- [API Endpoints](#api-endpoints)
- [Development Notes](#development-notes)
- [Contribution \& Development Workflow](#contribution--development-workflow)
  - [Regenerating OpenAPI Types](#regenerating-openapi-types)
  - [Adding a New Endpoint](#adding-a-new-endpoint)
  - [Testing](#testing)
  - [Contributing Guidelines](#contributing-guidelines)

---

## Overview

**Gochain** demonstrates a unified bridge layer that connects multiple blockchains through a common API.
It includes:

- A **Go-based backend** exposing REST endpoints defined via **OpenAPI 3.0**, generated using [`oapi-codegen`](https://github.com/oapi-codegen/oapi-codegen).
- A modular **adapter system** supporting GoChain, Ethereum, and Bitcoin demo integrations.
- A **Next.js frontend** built with **React Query**, **TailwindCSS**, and **Material UI**, integrating directly with **Metamask**.

The goal is to showcase how a single service can abstract wallet operations, transactions, and chain-specific logic across networks, while providing a type-safe, validated API surface.

---

## Architecture

```plaintext
.
├── api/                  # Go backend (OpenAPI + modular adapters)
│   ├── cmd/server/       # main.go entrypoint
│   ├── internal/
│   │   ├── adapters/     # chain-specific adapters (gochain, ethereum, bitcoin)
│   │   ├── core/         # core types & registry
│   │   ├── http/         # OpenAPI handlers and middleware
│   │   ├── platform/     # config, logging, middleware
│   │   └── service/      # business logic layer (wallet service)
│   ├── openapi.yaml      # API definition
│   └── oapi-codegen.yaml # codegen config
└── web/                  # Next.js frontend (App Router)
    └── src/app/          # dashboard UI, wallet & keypair components
```

---

## Usage

### Requirements

Install:

- **Go** ([install](https://go.dev/dl))
- **pnpm** ([installation guide](https://pnpm.io/installation))
- **git** ([download](https://git-scm.com/downloads))
- An **Infura API key** (for Ethereum testnet access) — [get one](https://infura.io/)
- Optional: **Metamask** wallet ([download](https://metamask.io/download/))
- Optional: **Docker and Docker Compose** ([install](https://docs.docker.com/get-docker/))

---

### Setup and Running

#### 1. Clone the repository

```bash
git clone https://github.com/afrodynamic/gochain.git
cd gochain
```

---

#### 2. Configure environment variables

**`web/.env.local`**

```env
NEXT_PUBLIC_API_BASE_URL=http://localhost:8080
NEXT_PUBLIC_INFURA_API_KEY=your_infura_api_key_here
```

---

#### 3. Run the backend

```bash
cd api
make init
make run
```

Backend available at: [http://127.0.0.1:8080](http://127.0.0.1:8080)

> Note: You can also run the backend using Docker from the project root:

```bash
docker-compose up
```

---

#### 4. Run the frontend

```bash
cd ../web
pnpm install
pnpm run dev
```

Frontend available at:

> [http://localhost:3000](http://localhost:3000)

---

## API Endpoints

| Method | Path                    | Description                                      |
| ------ | ----------------------- | ------------------------------------------------ |
| `GET`  | `/v1/health`            | Health check (returns status + timestamp)        |
| `GET`  | `/v1/chains`            | List supported chains                            |
| `POST` | `/v1/keys/new`          | Generate a new keypair (random or deterministic) |
| `GET`  | `/v1/balance/{address}` | Query balance for an address                     |
| `POST` | `/v1/tx/build`          | Build a new transaction                          |
| `POST` | `/v1/tx/sign`           | Sign a transaction                               |
| `POST` | `/v1/tx/broadcast`      | Broadcast a signed transaction                   |
| `GET`  | `/v1/tx/{id}`           | Get transaction status                           |

---

## Development Notes

- **OpenAPI-first design:**
  All backend routes are defined in `api/openapi.yaml` and validated at runtime via `oapi-codegen/nethttp-middleware`.

- **Key generation:**

  - `mode=random`: generates a new cryptographically secure random keypair.
  - `mode=deterministic`: generates reproducible keys from a provided `seed` or `passphrase`.
  - For demo purposes, private keys are included in responses.

- **Frontend integration:**

  - Metamask handles client-side EVM wallet actions (sign, send, chain switching).
  - Backend adapters simulate blockchain functionality and bridge operations.
  - Together, they demonstrate full-stack interoperability.

---

## Contribution & Development Workflow

### Regenerating OpenAPI Types

Whenever you update `api/openapi.yaml`, regenerate your Go types and server interface:

```bash
make generate
```

This updates:

- `api/internal/http/openapi/*.gen.go` (server + models)
- `api/internal/http/handlers` (implementation interface)

---

### Adding a New Endpoint

1. **Edit the OpenAPI spec**
   Add your new route in `api/openapi.yaml` under `paths`.

2. **Regenerate the code**

   ```bash
   make generate
   ```

3. **Implement the handler**
   Add a method to `api/internal/http/handlers/handlers.go` implementing the generated interface.

4. **Write tests**
   Add tests in `api/internal/http/handlers/handlers_test.go`.

5. **Run everything locally**

   ```bash
   make test
   make run
   ```

6. **Frontend integration (optional)**
   Add a hook under `web/src/app/(features)` using `React Query` to call your new endpoint.

---

### Testing

```bash
make test
```

- Uses Go’s built-in testing with race detection and coverage.
- Validates backend logic, adapters, and handler responses.

---

### Contributing Guidelines

- Keep Go code **idiomatic, modular, and testable**.
- Maintain **OpenAPI compliance** — all endpoints must be defined in `openapi.yaml`.
- Follow the **adapter pattern** for any new blockchain integration.
- Write concise, focused commits with meaningful messages.
- Before pushing:

  ```bash
  make tidy && make test
  ```

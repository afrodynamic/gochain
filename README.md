# Gochain <!-- omit in toc -->

A full-stack blockchain interoperability demo built in **Go** with **gRPC** + **REST (grpc-gateway)** and a **Next.js** frontend powered by **ConnectRPC**.

---

## Table of Contents<!-- omit in toc -->

- [Overview](#overview)
- [Usage](#usage)
  - [Requirements](#requirements)
  - [Setup and Running](#setup-and-running)
    - [1. Clone the repository](#1-clone-the-repository)
    - [2. Configure environment variables](#2-configure-environment-variables)
    - [3. Run the backend](#3-run-the-backend)
    - [4. Run the frontend](#4-run-the-frontend)
- [API Access](#api-access)
  - [gRPC](#grpc)
  - [REST (grpc-gateway)](#rest-grpc-gateway)
  - [Project Capabilities Summary](#project-capabilities-summary)

---

## Overview

**Gochain** demonstrates a unified blockchain abstraction layer with gRPC and REST endpoints.
It supports multiple blockchain “adapters” — GoChain, Ethereum, and Bitcoin — under a single service.

Core features:

- **gRPC** core API for fast, type-safe inter-service communication
- **grpc-gateway** automatic REST/JSON translation
- **gRPC-Web** for browser clients (Next.js)
- Modular **adapter pattern** for multi-chain support
- Example **GoChain blockchain implementation** (PoW/PoS-swappable)
- Fully deployable as a single-port Cloud Run container

---

## Usage

### Requirements

Install:

- **Go** ([install](https://go.dev/dl))
- **pnpm** ([installation guide](https://pnpm.io/installation))
- **git** ([download](https://git-scm.com/downloads))
- **buf** ([installation guide](https://docs.buf.build/installation))
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

#### 2. Configure environment variables

**`web/.env.local`**

```env
NEXT_PUBLIC_RPC_URL=http://localhost:8080
CHAIN=gochain
```

#### 3. Run the backend

```bash
cd api
go run ./cmd/gochaind
```

Server listens on **:8080** for gRPC, gRPC-Web, and REST.

> Or via Docker:

```bash
docker compose up --build
```

#### 4. Run the frontend

```bash
cd ../web
pnpm install
pnpm dev
```

Frontend: [http://localhost:3000](http://localhost:3000)

---

## API Access

### gRPC

```bash
grpcurl -plaintext localhost:8080 list
grpcurl -plaintext -d '{"height":"0"}' localhost:8080 chain.v1.Chain/GetBlock
grpcurl -plaintext -d '{}' localhost:8080 wallet.v1.Wallet/NewKey
grpcurl -plaintext -d '{"service":""}' localhost:8080 grpc.health.v1.Health/Check
```

> **Note:**
> The unified server supports raw gRPC on **:8080** via **h2c** (HTTP/2 Cleartext).
> If you deploy behind HTTPS or a proxy (e.g. Cloud Run, NGINX), add
> `--authority <service-name>` to `grpcurl`.

---

### REST (grpc-gateway)

```bash
curl http://localhost:8080/health
curl -X POST http://localhost:8080/v1/wallet:key -H 'content-type: application/json' -d '{}'
curl http://localhost:8080/v1/wallet/0xabc/balance
```

- All REST endpoints are automatically exposed from gRPC services through `grpc-gateway`.
- `/health` is always available for probes and container health checks.
- JSON routes mirror your gRPC definitions.

---

### Project Capabilities Summary

| Layer                        | Description                                                                 |
| ---------------------------- | --------------------------------------------------------------------------- |
| **Core blockchain**          | Minimal GoChain reference blockchain (PoW/PoS-swappable), state in memory.  |
| **Adapters**                 | Pluggable architecture for multiple chains — GoChain, Ethereum, Bitcoin.    |
| **gRPC API**                 | Type-safe RPC layer with reflection, health, and extensible services.       |
| **REST Gateway**             | JSON/HTTP access automatically generated via `grpc-gateway`.                |
| **gRPC-Web**                 | Browser clients (Next.js frontend) communicate natively through ConnectRPC. |
| **Health & CORS**            | `/health` endpoint and full CORS support for local dev & Cloud Run.         |
| **Single-Port Deployment**   | gRPC, gRPC-Web, REST, and health served together on `:8080`.                |
| **Docker & Cloud Run Ready** | Distroless image, minimal runtime footprint, deploys cleanly.               |
| **Frontend Integration**     | Next.js 14 + React Query + ConnectRPC + Tailwind + MUI dashboard.           |

You can now:

- Generate and sign wallets or transactions via **gRPC**, **REST**, or **frontend**.
- Switch between adapters (GoChain, Ethereum, Bitcoin) instantly.
- Build and serve from a single binary for both local and cloud environments.
- Access all APIs interchangeably (curl, grpcurl, or browser).

# Products API

A production-ready REST API built with **Go**, **Fiber**, and **gormx**.

## Features

- ✅ **JWT Authentication** — Register & Login
- ✅ **Products CRUD** — Create, Read, Update, Delete
- ✅ **Items (Nested)** — Items belong to products
- ✅ **Pagination** — Page-based with metadata
- ✅ **Filtering** — By category, price range
- ✅ **Search** — By name or SKU
- ✅ **Validation** — On all inputs
- ✅ **Error Handling** — Centralized
- ✅ **Graceful Shutdown** — Production-ready
- ✅ **API Discovery** — Info, Endpoints, Schema

## Quick Start

```bash
# Install dependencies
go mod tidy

# Run
go run ./cmd/server
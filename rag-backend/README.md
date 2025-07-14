# RAG Backend API

A Go-based REST API for managing documents with vector embeddings using PostgreSQL and pgvector.

## Features

- 📄 **Document Management**: Upload, retrieve, update, and delete documents
- 🔍 **Vector Search**: Store and search document embeddings
- 🗂️ **File Operations**: Filename search and date-range filtering
- 🗑️ **Soft Delete**: Recycle bin functionality with restore capability
- 📊 **Metadata**: Lightweight file metadata retrieval
- 🔒 **Validation**: Input validation and error handling
- 📖 **Documentation**: Swagger/OpenAPI documentation

## Quick Start

### Using Docker

Pull and run the latest image:

```bash
# Pull the Docker image
docker pull ghcr.io/fain17/rag-backend:latest

# Run with database connection
docker run -d \
  --name rag-backend \
  -p 8080:8080 \
  -e DATABASE_URL="postgres://username:password@localhost:5432/rag_db?sslmode=disable" \
  -e GIN_MODE=release \
  ghcr.io/fain17/rag-backend:latest
```

### Environment Variables

| Variable | Required | Description | Example |
|----------|----------|-------------|---------|
| `DATABASE_URL` | Yes | PostgreSQL connection string | `postgres://user:pass@localhost:5432/db?sslmode=disable` |
| `GIN_MODE` | No | Gin framework mode | `release` (default: `debug`) |

### Database Connection Examples

```bash
# Local PostgreSQL
DATABASE_URL="postgres://rag_user:rag_password@localhost:5432/rag_database?sslmode=disable"

# PostgreSQL with SSL
DATABASE_URL="postgres://rag_user:rag_password@localhost:5432/rag_database?sslmode=require"

# Docker Compose PostgreSQL
DATABASE_URL="postgres://rag_user:rag_password@postgres:5432/rag_database?sslmode=disable"
```

## API Endpoints

### Files
- `GET /files/{id}` - Get file by ID
- `GET /files/getall` - Get all files
- `GET /files/search?query={query}` - Search files by filename
- `GET /files/date-range?start={date}&end={date}` - Get files by date range
- `GET /files/metadata` - Get file metadata
- `POST /files/upload` - Upload new file
- `PUT /files/{id}` - Update file
- `DELETE /files/{id}` - Delete file permanently

### Recycle Bin
- `PATCH /files/{id}/soft-delete` - Soft delete file
- `PATCH /files/{id}/restore` - Restore soft-deleted file
- `GET /files/recycle-bin` - Get all soft-deleted files

### Documentation
- `GET /docs/swagger/index.html` - Swagger UI
- `GET /swagger/doc.json` - OpenAPI JSON spec

## API Documentation

Once running, visit: `http://localhost:8080/docs/swagger/index.html`

## Health Check

```bash
curl http://localhost:8080/files/metadata
```

## Development

### Local Development

```bash
# Install dependencies
go mod download

# Run tests
./test/run_tests.sh

# Run linting
./linting/lint.sh

# Start server
go run main.go
```

### Database Requirements

Requires PostgreSQL with pgvector extension:

```sql
CREATE EXTENSION IF NOT EXISTS vector;
```

## License

See [LICENSE](LICENSE) file.

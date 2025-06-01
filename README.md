# Kinko

Kinko is a simple and self-hosted secret management web application. It provides a web UI for securely storing, retrieving, and managing secrets (such as API keys, passwords, or configuration values) for your projects or personal use.

## Features

- Manage secrets through a web interface
- Store secrets in a lightweight SQLite database
- Authentication support (Basic Auth)
- RESTful API for secret operations
- SvelteKit-powered frontend for a modern user experience
- Easy to run locally or deploy anywhere

## Getting Started

### Prerequisites

- Go (for backend)
- Node.js and npm (for frontend development)
- SQLite (included by default)

### Running the Application

#### 1. Backend

```bash
# Run the backend (Windows example)
CGO_ENABLED=1 go run cmd/main/main.go
```

You must provide the following flags or environment variables:
- `--auth-token` or `API_TOKEN`: API authentication token (required)
- `--enc-key` or `ENC_KEY`: Encryption key for secrets (required)
- `--auth-user` and `--auth-passwd`: Username and password for Basic Auth (optional but recommended)
- `--port` and `--host`: Port and host to listen on (optional)

#### 2. Frontend

```bash
cd frontend
npm install
npm run dev
```

The frontend is a SvelteKit app served via the backend by default.

## Usage

- Access the web interface in your browser.
- Add, view, or delete secrets.
- Secrets are stored securely in the backend and shown in a table view.

## Project Structure

- `cmd/main/`: Go backend server
- `frontend/`: SvelteKit frontend app
- `api/`: API routes and middleware
- `internal/`: Internal Go packages (logging, config)

## License

This project is licensed under the terms of the GNU General Public License (GPL).

# Password Generator REST API

A production-style REST API for generating cryptographically secure passwords.

This is Project 2 of my software engineering and DevOps portfolio. It expands
the password-generation logic from Project 1 into a documented HTTP service
built with FastAPI.

## Features

- Generate one or more secure passwords over HTTP
- Configurable password length
- Optional punctuation symbols
- Password strength ratings
- Structured JSON responses
- Automatic request validation
- Swagger and OpenAPI documentation
- Health-check endpoint
- Application logging
- Unit and integration tests
- Ruff formatting and linting
- GitHub Actions continuous integration
- Docker support

## Technologies

- Python
- FastAPI
- Pydantic
- pytest
- HTTPX
- Ruff
- Docker
- GitHub Actions

## Project structure

```text
02-password-generator-api/
│
├── app/
│   ├── __init__.py
│   ├── main.py
│   ├── models.py
│   └── services.py
│
├── tests/
│   ├── __init__.py
│   └── test_api.py
│
├── .github/workflows/python.yml
├── .dockerignore
├── .gitignore
├── Dockerfile
├── LICENSE
├── README.md
├── pyproject.toml
└── requirements.txt
```

## Installation

Open PowerShell in this project directory.

Create a virtual environment:

```powershell
python -m venv .venv
```

Activate it:

```powershell
.venv\Scripts\Activate.ps1
```

Install dependencies:

```powershell
python -m pip install --upgrade pip
python -m pip install -r requirements.txt
python -m pip install ruff
```

## Run the API

Start the development server:

```powershell
fastapi dev app/main.py
```

The API will normally be available at:

```text
http://127.0.0.1:8000
```

Interactive Swagger documentation:

```text
http://127.0.0.1:8000/docs
```

Alternative ReDoc documentation:

```text
http://127.0.0.1:8000/redoc
```

## API endpoints

### Root

```http
GET /
```

### Health check

```http
GET /health
```

### Generate passwords

```http
GET /api/v1/passwords
```

Supported query parameters:

| Parameter | Type | Default | Limits |
|---|---:|---:|---:|
| `length` | integer | `16` | `4–256` |
| `count` | integer | `1` | `1–50` |
| `symbols` | boolean | `false` | — |

Example request:

```text
http://127.0.0.1:8000/api/v1/passwords?length=24&count=3&symbols=true
```

Example response:

```json
{
  "length": 24,
  "count": 3,
  "symbols_included": true,
  "passwords": [
    {
      "password": "r7#Kq2!Bv8@Xm4Lp9$Nz1AaQ",
      "strength": "Excellent"
    }
  ]
}
```

## Run tests

```powershell
python -m pytest
```

## Check formatting and linting

```powershell
python -m ruff format --check .
python -m ruff check .
```

Automatically format the project:

```powershell
python -m ruff format .
```

## Docker

Build the image from the project directory:

```powershell
docker build -t password-generator-api .
```

Run the container:

```powershell
docker run --rm -p 8000:8000 password-generator-api
```

Then open:

```text
http://127.0.0.1:8000/docs
```

## Architecture

The project separates responsibilities into three main layers:

- `main.py` defines HTTP routes, middleware, logging, and API behavior.
- `models.py` defines validated response schemas.
- `services.py` contains reusable password-generation business logic.

This separation keeps the generation algorithm independent from FastAPI and
makes it easier to test, reuse, and replace.

## Security considerations

- Passwords are generated with Python's `secrets` module.
- Generated passwords are never written to disk.
- Generated passwords are not stored in logs.
- Query values are constrained to prevent extreme resource usage.
- The API does not provide password history or persistent storage.

## Future improvements

- Add API-key authentication
- Add rate limiting
- Add structured JSON logging
- Add Prometheus metrics
- Add a React frontend
- Deploy to Azure
- Provision infrastructure with Terraform
- Deploy through Kubernetes

## License

MIT

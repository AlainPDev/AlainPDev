"""
FastAPI application for the Password Generator API

Run locally with:

    fastapi dev app/main.py

Interactive documentation:

    http://127.0.0.1:8000/docs

"""

import logging
from contextlib import asynccontextmanager
from typing import Annotated

from fastapi import FastAPI, Query, Request
from fastapi.middleware.cors import CORSMiddleware
from fastapi.responses import JSONResponse

from app.models import PasswordItem, PasswordResponse
from app.services import (
    Max_Password_Count,
    Max_Password_Length,
    Min_Password_Length,
    calculate_strength,
    generate_password,
)

logging.basicConfig(
    level=logging.INFO,
    format=("%(asctime)s | %(levelname)s | %(name)s | %(message)s"),
)

logger = logging.getLogger("password_generator_api")


@asynccontextmanager
async def lifespan(_: FastAPI):
    """Log application startup and shutdown events."""

    logger.info("Password Generator API starting.")
    yield
    logger.info("Passowrd Generator API shutting down.")


app = FastAPI(
    title="Password Generator API",
    summary="Generate cryptographically secure passwords over HTTP.",
    description=(
        "A portfolio REST API built with FastAPI. It generates secure "
        "passwords, validates query parameters, calculates strength ratings, "
        "and returns structured JSON responses."
    ),
    version="1.0.0",
    lifespan=lifespan,
)

# This allows a future React frontend to call the API during local development.

app.add_middleware(
    CORSMiddleware,
    allow_origins=[
        "http://localhost:3000",
        "http://localhost:5173",
    ],
    allow_credentials=True,
    allow_methods=["GET"],
    allow_headers=["*"],
)


@app.exception_handler(ValueError)
async def value_error_handler(
    request: Request,
    exception: ValueError,
) -> JSONResponse:
    """Convert service-level ValueError exceptions into HTTP 400 responses."""

    logger.warning(
        "Invalid request for %s: %s",
        request.url.path,
        exception,
    )

    return JSONResponse(
        status_code=400,
        content={
            "error": "invalid_request",
            "detail": str(exception),
        },
    )


@app.get(
    "/",
    tags=["System"],
    summary="API welcome route",
)
async def root() -> dict[str, str]:
    """Return basic API information."""

    return {
        "name": "Password Generator API",
        "version": "1.0.0",
        "documentation": "/docs",
        "health": "/health",
    }


@app.get(
    "/health",
    tags=["System"],
    summary="Check API health",
)
async def health_check() -> dict[str, str]:
    """Return a simple health-check response."""

    return {
        "status": "healthy",
    }


@app.get(
    "/api/v1/passwords",
    response_model=PasswordResponse,
    tags=["Passwords"],
    summary="Generate secure passwords",
)
async def create_passwords(
    length: Annotated[
        int,
        Query(
            ge=Min_Password_Length,
            le=Max_Password_Length,
            description="Length of each generated password.",
        ),
    ] = 16,
    count: Annotated[
        int,
        Query(
            ge=1,
            le=Max_Password_Count,
            description="Number of passwords to generate.",
        ),
    ] = 1,
    symbols: Annotated[
        bool,
        Query(
            description="Include punctuation symbols.",
        ),
    ] = False,
) -> PasswordResponse:
    """
    Generate one or more cryptographically secure passwords.

    Query parameters are validated automatically by FastAPI before this
    function runs.
    """

    logger.info(
        "Generating %s password(s), length=%s, symbols=%s.",
        count,
        length,
        symbols,
    )

    generated_passwords = []

    for _ in range(count):
        password = generate_password(
            length=length,
            include_symbols=symbols,
        )

        generated_passwords.append(
            PasswordItem(
                password=password,
                strength=calculate_strength(password),
            )
        )

    return PasswordResponse(
        length=length,
        count=count,
        symbols_included=symbols,
        passwords=generated_passwords,
    )

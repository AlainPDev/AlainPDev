"""
Tests for the Password Generator API.
"""

import string

from fastapi.testclient import TestClient

from app.main import app

client = TestClient(app)


def test_root_endpoint() -> None:
    """The root endpoint should describe the API."""

    response = client.get("/")

    assert response.status_code == 200
    assert response.json()["name"] == "Password Generator API"
    assert response.json()["documentation"] == "/docs"


def test_health_endpoint() -> None:
    """The health endpoint should report that the API is healthy."""

    response = client.get("/health")

    assert response.status_code == 200
    assert response.json() == {"status": "healthy"}


def test_generate_default_password() -> None:
    """The default request should return one 16-character password."""

    response = client.get("/api/v1/passwords")

    assert response.status_code == 200

    data = response.json()

    assert data["length"] == 16
    assert data["count"] == 1
    assert data["symbols_included"] is False
    assert len(data["passwords"]) == 1
    assert len(data["passwords"][0]["password"]) == 16


def test_generate_multiple_passwords() -> None:
    """The count query parameter should control response size."""

    response = client.get(
        "/api/v1/passwords",
        params={
            "length": 24,
            "count": 3,
            "symbols": "true",
        },
    )

    assert response.status_code == 200

    data = response.json()

    assert data["length"] == 24
    assert data["count"] == 3
    assert data["symbols_included"] is True
    assert len(data["passwords"]) == 3

    for item in data["passwords"]:
        password = item["password"]

        assert len(password) == 24
        assert any(character.islower() for character in password)
        assert any(character.isupper() for character in password)
        assert any(character.isdigit() for character in password)
        assert any(character in string.punctuation for character in password)


def test_reject_password_that_is_too_short() -> None:
    """Lengths below four should fail validation."""

    response = client.get(
        "/api/v1/passwords",
        params={"length": 3},
    )

    assert response.status_code == 422


def test_reject_password_that_is_too_long() -> None:
    """Lengths above the configured maximum should fail validation."""

    response = client.get(
        "/api/v1/passwords",
        params={"length": 257},
    )

    assert response.status_code == 422


def test_reject_invalid_count() -> None:
    """Counts below one should fail validation."""

    response = client.get("/api/v1/passwords", params={"count": 0})

    assert response.status_code == 422


def test_openapi_document_is_available() -> None:
    """FastApi should publish an Open API schema."""

    response = client.get("/openapi.json")

    assert response.status_code == 200
    assert response.json()["info"]["title"] == "Password Generator API"

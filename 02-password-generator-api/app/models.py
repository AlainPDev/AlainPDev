"""
Pydantic response models for the Password Generator API.

FastAPI uses these models to validate responses and generate
the OpenAPI documentation displayed in Swagger UI.
"""

from pydantic import BaseModel, Field


class PasswordItem(BaseModel):
    """Represent one generated password and its strength rating."""

    password: str = Field(description="Cryptographically secure generated password.")

    strength: str = Field(description="Human-readable password strength rating.")


class PasswordResponse(BaseModel):
    """Represent the complete password-generation response."""

    length: int = Field(description="Length of each generated password.")

    count: int = Field(description="Number of passwords generated.")

    symbols_included: bool = Field(
        description="Whether punctuation symbols were included."
    )

    passwords: list[PasswordItem] = Field(
        description="Generated passwords and their strength ratings."
    )

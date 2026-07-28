"""
Core password-generation services.

The functions in this module are independent of FastAPI. Keeping business
logic separate from HTTP routing makes it easier to test and reuse.
"""

import secrets
import string

Min_Password_Length = 4
Max_Password_Length = 256
Max_Password_Count = 50

Password_Ratings = {
    1: "Very Weak",
    2: "Weak",
    3: "Moderate",
    4: "Strong",
    5: "Excellent",
}


def build_character_pool(include_symbols: bool) -> str:
    """Build the complete pool of characters available for generation."""

    characters = string.ascii_lowercase + string.ascii_uppercase + string.digits

    if include_symbols:
        characters += string.punctuation

    return characters


def generate_password(length: int, include_symbols: bool) -> str:
    """
    Generate one cryptographically secure password.

    The password always contains at least one lowercase letter, one uppercase letter,
    and one digit. A punctuation character is also guaranteed when symbols
    are requested.

    Args:
        length: Desired password length.
        include_symbols: Whether punctuation symbols should be included.

    Returns:
        A securely generated password.

    Raises:
        ValueError: If the supplied length is outside the supported range.

    """

    if length < Min_Password_Length:
        raise ValueError(f"Password length cannot exceet {Max_Password_Length}.")

    required_characters = [
        secrets.choice(string.ascii_lowercase),
        secrets.choice(string.ascii_uppercase),
        secrets.choice(string.digits),
    ]

    if include_symbols:
        required_characters.append(secrets.choice(string.punctuation))

    character_pool = build_character_pool(include_symbols)

    while len(required_characters) < length:
        required_characters.append(secrets.choice(character_pool))

    # SystemRandom uses operating-system randomness for the shuffle.
    secrets.SystemRandom().shuffle(required_characters)

    return "".join(required_characters)


def calculate_strength(password: str) -> str:
    """Calculate a simple strength rating based on length and character mix."""

    score = 0

    if len(password) >= 12:
        score += 1

    if any(character.islower() for character in password):
        score += 1

    if any(character.isupper() for character in password):
        score += 1

    if any(character.isdigit() for character in password):
        score += 1

    if any(character in string.punctuation for character in password):
        score += 1

    return Password_Ratings.get(score, "Very Weak")

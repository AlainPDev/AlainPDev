"""

Password Generator CLI
Author: Alain P.
Project 1

Description: A secure command-line password generator using Python's
cryptographically secure 'secrets' module.

Features: Configurable password length, Multiple password generation,
Optional symbols, Clipboard support, Password strength analysis

"""

import argparse
import secrets
import string

import pyperclip

# ==================================================================
# Configuration Constants
# ==================================================================

Default_Password_Length = 16
Default_Password_Count = 1
Min_Password_Length = 4

Password_Ratings = {
    1: "Very Weak",
    2: "Weak",
    3: "Moderate",
    4: "Strong",
    5: "Excellent",
}

# ==================================================================
# Command Line Arguments
# ==================================================================

def parse_arguments() -> argparse.Namespace:
    """
    Parse command-line arguments.

    Returns:
        argparse.Namespace
    """

    parser = argparse.ArgumentParser(
    description="Generate secure passwords."
    )

    parser.add_argument(
    "--length",
    type=int,
    default=Default_Password_Length,
    help="Length of each password.",
    )

    parser.add_argument(
    "--count",
    type=int,
    default=Default_Password_Count,
    help="Number of passwords to generate.",
    )

    parser.add_argument(
    "--symbols",
    action="store_true",
    help="Include punctuation symbols.",
    )

    parser.add_argument(
    "--copy",
    action="store_true",
    help="Copy the first generated password to the clipboard.",
    )

    return parser.parse_args()

# ======================================================================
# Character Pool
# ======================================================================

def build_character_pool(include_symbols: bool) -> str:
    """
    Build the available chracter pool
    """

    characters = (
        string.ascii_lowercase
        + string.ascii_uppercase
        + string.digits
    )

    if include_symbols:
        characters += string.punctuation

    return characters

# =====================================================================
# Password Generation
# =====================================================================

def generate_password(length: int, include_symbols: bool) -> str:
    """
    Generates a cryptographically secure password.
    The algorithm guarantees at least one lowercase, one uppercase, digit,
    and one symbol.
    """

    if length < Min_Password_Length:
        raise ValueError("Password length is too short.")

    password = [
        secrets.choice(string.ascii_lowercase),
        secrets.choice(string.ascii_uppercase),
        secrets.choice(string.digits),
    ]

    pool = build_character_pool(include_symbols)

    if include_symbols:
        password.append(secrets.choice(string.punctuation))

    while len(password) < length:
        password.append(secrets.choice(pool))

    # Shuffle securely so required characters are random positions.

    secrets.SystemRandom().shuffle(password)

    return "".join(password[:length])

# ========================================================================
# Password Strength
# ========================================================================

def password_strength(password: str) -> str:
    """
    Return a human-readable strength rating.
    """

    score = 0

    if len(password) >= 12:
        score += 1

    if any(c.islower() for c in password):
        score += 1

    if any(c.isupper() for c in password):
        score += 1

    if any(c.isdigit() for c in password):
        score += 1

    if any(c in string.punctuation for c in password):
        score += 1

    return Password_Ratings.get(score, "Very Weak")

# ============================================================
# Display
# ============================================================

def display_passwords(passwords: list[str]) -> None:
    """
    Print generated passwords.
    """

    print()

    for index, password in enumerate(passwords, start=1):

        print(f"Password {index}")

        print("-" * 40)

        print(password)

        print(f"Strength : {password_strength(password)}")

        print()

# =======================================================================
# Main Application
# =======================================================================

def main() -> None:

    args = parse_arguments()

    if args.length < Min_Password_Length:

        print(
            f"Password length must be at least "
            f"{Min_Password_Length}."
        )

        return

    if args.count < 1:

        print("Count must be greater than zero.")

        return

    passwords = []

    for _ in range(args.count):

        password = generate_password(
            args.length,
            args.symbols,
        )

        passwords.append(password)

    display_passwords(passwords)

    if args.copy:

        pyperclip.copy(passwords[0])

        print("First password copied to clipboard.")

if __name__ == "__main__":
    main()

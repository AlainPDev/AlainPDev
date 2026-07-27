"""
Unit tests for Password Generator CLI.
"""

import string

from password_generator import generate_password

def test_password_length():

    password = generate_password(20, False)

    assert len(password) == 20

def test_contains_uppercase():

    password = generate_password(20, False)

    assert any(c.isupper() for c in password)

def test_contains_lowercase():

    password = generate_password(20, False)

    assert any(c.islower() for c in password)

def test_contains_digit():

    password = generate_password(20, False)

    assert any(c.isdigit() for c in password)

def test_contains_symbol():

    password = generate_password(20, True)

    assert any(c in string.punctuation for c in password)

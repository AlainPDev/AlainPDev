# Password Generator CLI

A secure command-line password generator written in Python.

This project was built as part of my Software Engineering Portfolio. 

---

## Features

- Cryptographically secure password generation
- Configurable password length
- Generate multiple passwords
- Optional symbols
- Clipboard support
- Password strength analysis
- Type hints
- Clean architecture
- Unit tests
- Github Actions ready

---

## Technologies

- Python 3
- argparse
- secrets
- pyperclip
- pytest
- Ruff

---

## Installation

Clone the repository. 

```bash
git clone https://github.com/AlainPDev/AlainPDev.git
```

Install dependencies. 

```bash
pip install -r requirements.txt
```

---

## Usage

Generate a default password.

```bash
python password_generator.py
```

Generate a 24-character password.

```bash
python password_generator.py --length 24
```

Generate five passwords.

```bash
python password_generator.py --count 5
```

Generate passwords with symbols.

```bash
python password_generator.py --symbols
```

Copy the first password to the clipboard.

```bash
python password_generator.py --copy
```

Combine options.

```bash
python password_generator.py --length 32 --symbols --count 3 --copy
```

---

## Example

```
Password 1
----------------------------------------
L4v!Q8#bRz7@Xn9$TwPm
Strength : Excellent
```

---

## Project Structure

```
password-generator-cli/
│
├── password_generator.py
├── README.md
├── requirements.txt
├── pyproject.toml
├── tests/
└── .github/
```

---

## Future Improvements

- Export passwords to a file
- Password history
- Configuration file
- GUI version
- Web version using FastAPI
- Docker container

---

## License

MIT

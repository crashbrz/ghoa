# ghoa - GitHub OAuth Token Validator

ghoa is a powerful tool designed to validate and retrieve information about GitHub OAuth tokens. It supports concurrency and allows validation of multiple tokens from a file, with additional features to retrieve user details, token scopes, and private repository information.

---

## Features

- Validate a single token or multiple tokens from a file.
- Retrieve user details, including name, email, and associated metadata.
- Display token scopes for valid tokens.
- List private repositories (names and URLs) for valid tokens.
- Support for concurrency using goroutines for efficient file token validation.
- Flexible endpoint support with a default of `https://api.github.com/user`.

---

## Installation

Clone the repository and build the tool using Go:

```bash
git clone <repository-url>
cd <repository-folder>
go build -o ghoa
```

---

## Usage

### Command-Line Flags

| Flag             | Description                                                                                           | Example                                                                                 |
|-------------------|-------------------------------------------------------------------------------------------------------|-----------------------------------------------------------------------------------------|
| `-k <token>`     | Validate a single GitHub OAuth token.                                                                 | `ghoa -k <your_token>`                                                                  |
| `-f <file>`      | File containing GitHub tokens (one per line).                                                         | `ghoa -f tokens.txt`                                                                    |
| `-t <number>`    | Number of goroutines to use when processing multiple tokens (default: 1).                             | `ghoa -f tokens.txt -t 10`                                                              |
| `-i`             | Retrieve and display detailed information about valid tokens, including scopes.                       | `ghoa -k <your_token> -i`                                                               |
| `-p`             | Retrieve and display private repositories for valid tokens.                                           | `ghoa -k <your_token> -p`                                                               |
| `-d`             | Display invalid tokens when validating multiple tokens.                                               | `ghoa -f tokens.txt -d`                                                                 |
| `-e <endpoint>`  | Specify a custom GitHub API endpoint (default: `https://api.github.com/user`).                        | `ghoa -k <your_token> -e https://api.github.com/user`                                   |
| `--remove-color` | Remove color formatting from the output.                                                              | `ghoa -k <your_token> --remove-color`                                                   |

---

### Examples

#### Validate a Single Token
```bash
ghoa -k <your_token>
```

#### Validate Tokens from a File with Concurrency
```bash
ghoa -f tokens.txt -t 5
```

#### Validate Tokens and Retrieve User Info and Scopes
```bash
ghoa -f tokens.txt -i
```

#### Validate Tokens and Display Private Repositories
```bash
ghoa -f tokens.txt -p
```

#### Validate Tokens with Invalid Token Output
```bash
ghoa -f tokens.txt -d
```

---

### Output

- **Valid Tokens**: Displayed in green with associated user details (if `-i` is set) and private repositories (if `-p` is set).
- **Invalid Tokens**: Displayed in red (only if `-d` is set).

---

### License

ghoa is licensed under the SushiWare license. For more information, check [docs/license.txt](docs/license.txt).

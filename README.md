## Open Trivia Online API 🎮

The `oto-api` is the backend service for Open Trivia Online app, written in Go. It provides APIs to manage trivia games, questions, and user interactions.

## Project Setup ⚙️

1. Clone the repository:

   ```bash
   git clone https://github.com/your-repo/oto-api.git
   cd oto-api
   ```

2. Install dependencies 📦:

   ```bash
   go mod tidy
   ```

3. Copy the `.env` file 📝:
   ```bash
   cp .env.example .env
   ```
   Update the `.env` file with your local configuration.

## Running the Project ▶️

To start the project locally, use the CLI:

```bash
go run .
```

## Running Migrations 🗄️

To run database migrations, use the following command:

```bash
go run . migrate
```

Ensure your database is configured correctly in the `.env` file before running migrations.

## Running Tests

To run all unit tests, use the following command:

```bash
go test ./... -v
```

# Mobo Ticketing Service

Mobo is a modern, high-performance backend API for a movie theatre ticketing and reservation system. It provides a robust set of features for managing users, movies, showtimes, and venues, as well as aggregated analytics. The system supports secure authentication (including OIDC built on OAuth 2.0), role-based access control (Admin vs Customer), and transactionally safe bookings.

## System Architecture

Mobo is designed for high availability and is deployed in a production environment using a reverse proxy and multiple backend instances.

```mermaid
flowchart TD
    Client["React client"] -->|HTTPS| Caddy["Caddy reverse proxy"]

    subgraph VPS["Production VPS"]
        Caddy --> API1["Mobo API :3000"]
        Caddy --> API2["Mobo API :3001"]
        Caddy --> API3["Mobo API :3002"]
        Worker["mobo-email-worker<br>systemd service"]
    end

    subgraph Cloud["Managed services"]
        DB[("PostgreSQL")]
        MQ["RabbitMQ<br>environment vhost"]
        Mail["Email provider<br>Resend"]
    end

    API1 --> DB
    API2 --> DB
    API3 --> DB
    Worker -->|"claim outbox rows"| DB
    Worker -->|"publish / consume"| MQ
    Worker -->|"send"| Mail
```

- **Frontend:** A React web interface served to the client.
- **Reverse Proxy:** Caddy handles SSL termination and load balances traffic across the backend instances.
- **Backend:** 3 independent instances of the Go API running as `systemd` services on a VPS, ensuring redundancy and fault tolerance.
- **Database:** NeonDB, a serverless cloud PostgreSQL provider, handles all data persistence.
- **Async worker:** a separate `systemd` service relays transactional-email intents from PostgreSQL, consumes RabbitMQ jobs, and sends email. It is not an HTTP service and does not sit behind Caddy.
- **Broker:** RabbitMQ buffers email work and provides delayed retries and a dead-letter queue. A managed broker is recommended for production; a private, self-hosted broker is supported for a small deployment.

### Software Architecture (Domain-Driven Design)

I'm trying to follow a strict **Domain-Driven Design (DDD)** pattern. The codebase is organised by business capabilities (domains) rather than technical layers. Technology adapters (like the HTTP server and Postgres repositories) depend downward on pure domain packages, ensuring business logic is isolated and testable.

### Project Structure

```text
ticketing-service/
├── cmd/server/                 # Application entry point. Initializes dependencies and starts the server.
├── cmd/email-worker/           # Outbox relay and asynchronous email consumer entry point.
├── config/                     # Configuration management (Viper) reading from .env.
├── deploy/                     # Deployment configurations (e.g., Caddyfile, systemd service files).
├── internal/
│   ├── auth/                   # Core authentication logic (JWT, OIDC).
│   ├── api/                    # HTTP Transport layer. Contains Chi router, handlers, and middleware.
│   ├── postgres/               # Database adapter layer. Manages pgx connection pooling.
│   ├── broker/rabbitmq/        # RabbitMQ topology, publisher confirms, and consumer adapter.
│   ├── email/                  # Provider-neutral jobs, templates, and delivery policy.
│   ├── mailer/resend/          # Resend implementation of the email Sender contract.
│   ├── outbox/                 # Transactional outbox relay contract and coordination.
│   ├── dbgen/                  # Auto-generated SQL code via sqlc.
│   │
│   │ # DOMAIN PACKAGES (Pure Business Logic)
│   ├── analytics/              # Aggregated dashboard metrics and revenue data.
│   ├── movie/                  # Movie catalog management.
│   ├── showtime/               # Scheduling and availability tracking for movies.
│   ├── user/                   # User identity, roles, and authentication workflows.
│   └── venue/                  # Physical theater location management.
│
├── pkg/                        # Reusable, domain-agnostic utilities (e.g., zap logger).
├── sql/                        # Raw SQL schema migrations and sqlc query definitions.
└── web/                        # React frontend / Web UI application.
```

## Tech Stack & Tools

- **Language:** Go (Golang)
- **Database:** PostgreSQL
- **Database Driver / Pool:** [pgx/v5](https://github.com/jackc/pgx)
- **Data Access:** [sqlc](https://sqlc.dev/) (Type-safe SQL compiler)
- **HTTP Router:** [go-chi/chi](https://github.com/go-chi/chi)
- **Validation:** [go-playground/validator](https://github.com/go-playground/validator)
- **Authentication:** JWT (JSON Web Tokens) with rotating HttpOnly cookies
- **OAuth:** [markbates/goth](https://github.com/markbates/goth) (Google Auth integration)
- **Logging:** [uber-go/zap](https://github.com/uber-go/zap) (Structured JSON logging)
- **Configuration:** [spf13/viper](https://github.com/spf13/viper)
- **Async messaging:** [RabbitMQ](https://www.rabbitmq.com/) via `amqp091-go`
- **Transactional email:** [Resend](https://resend.com/), behind a provider-neutral sender interface

## Key Features

- **Decoupled Architecture:** Business logic operates independently of how data is stored or served.
- **Secure Authentication:** Combines local credentials and OIDC (built on top of OAuth 2). Uses secure, HttpOnly, SameSite cookie-based JWTs with short-lived access tokens and longer-lived refresh tokens.
- **Role-Based Access:** Distinct `Admin` and `User` roles enforced via middleware.
- **Transaction Safety:** Repository methods safely wrap complex multi-step operations (like linking an OAuth identity to a new user profile) in atomic PostgreSQL transactions.
- **Performance Optimized:** Uses `pgxpool` for connection lifecycle management and `httprate` for IP-based rate limiting.
- **Durable async email:** New-account email intents are inserted into PostgreSQL's `email_outbox` in the same transaction as the user record. A background worker publishes them to RabbitMQ using publisher confirms, retries transient failures, and sends terminal failures to a DLQ.

## Getting Started

### Prerequisites

- Go 1.21+
- PostgreSQL database
- [sqlc](https://docs.sqlc.dev/en/latest/overview/install.html) (for modifying database queries)
- RabbitMQ only when running the asynchronous email worker. A local Docker broker is suitable for development.

### Setup

1. **Clone the repository:**

   ```bash
   git clone <repo-url>
   cd ticketing-service
   ```

2. **Configure environment:**
   Create a `.env` file in the root directory mirroring the necessary configuration (Database URI, JWT secret, OAuth credentials, etc.).

3. **Generate database code (if modifying queries):**

   ```bash
   sqlc generate
   ```

4. **Run the API:**
   ```bash
   go run cmd/server/main.go
   ```

### Running Tests

```bash
go test ./...
```

## Async transactional email

When a new local or OAuth account is created, Mobo writes both the user and a welcome-email intent to PostgreSQL in one transaction. This is the transactional outbox pattern: a successful user creation cannot lose its email intent because RabbitMQ or the email provider is temporarily unavailable.

The separate worker later claims the outbox row, publishes the job to RabbitMQ with publisher confirmation, renders the provider-neutral template, and delivers through Resend. The provider request uses the outbox ID as its idempotency key to reduce duplicate sends after crashes or network ambiguity.

```mermaid
flowchart LR
    API["Mobo API"] -->|"one PostgreSQL transaction"| User["users"]
    API -->|"same transaction"| Outbox["email_outbox"]
    Outbox --> Relay["email-worker relay"]
    Relay -->|"publisher confirm"| Queue["RabbitMQ: email.send"]
    Queue --> Worker["email-worker consumer"]
    Worker --> Resend["Resend"]
    Worker -->|"transient failure"| Retry["1m → 5m → 15m → 1h → 6h"]
    Retry --> Queue
    Worker -->|"permanent / exhausted"| DLQ["email.failed"]
```

### Run locally

Apply the migrations, start PostgreSQL and RabbitMQ, then configure the API and worker. `RABBITMQ_URL` is optional for the API but required for the worker.

```text
# Required by API and worker
DATABASE_URI=postgres://...
SYMMETRIC_KEY=a-minimum-32-character-secret

# Required by the worker
RABBITMQ_URL=amqp://mobo:password@localhost:5672/mobo-development
EMAIL_PROVIDER=resend
EMAIL_FROM=Mobo <bookings@example.com>
RESEND_API_KEY=re_...

# Optional worker tuning defaults shown
EMAIL_WORKER_COUNT=2
EMAIL_WORKER_PREFETCH=10
EMAIL_SEND_TIMEOUT=10s
OUTBOX_POLL_INTERVAL=5s
OUTBOX_LEASE_DURATION=1m
OUTBOX_BATCH_SIZE=25
```

Run the processes in separate terminals:

```bash
make run
make email-worker
```

The worker creates the RabbitMQ topology at startup. It requires the database migration that creates `email_outbox`; apply that migration before deploying an API version that writes welcome-email outbox records.

### Retry and failure policy

- **Success:** the worker acknowledges the RabbitMQ delivery after the provider accepts the message.
- **Transient/provider rate-limit error:** the job is delayed for 1 minute, 5 minutes, 15 minutes, 1 hour, then 6 hours.
- **Malformed job or permanent provider error:** the job is sent directly to `email.failed`.
- **Exhausted retries:** the job is sent to `email.failed` with its failure reason and retry metadata.

`email.failed` is intentionally not auto-consumed. Inspect and replay jobs only after correcting their root cause.

### Production deployment

The existing Caddy and three API systemd services remain HTTP-only. Add one new service, `mobo-email-worker`, built from the same commit but started as a separate binary. It does not need a Caddy route or public port.

1. Provision RabbitMQ. Prefer a managed, TLS-enabled broker with a dedicated production vhost and least-privilege user. For a small initial deployment, RabbitMQ can run on the VPS and be reached privately at `amqp://...@127.0.0.1:5672/<vhost>`; do not expose AMQP publicly.
2. Apply the `009_email_outbox.sql` migration before deploying the API change.
3. Build and release both binaries: `cmd/server/main.go` and `cmd/email-worker/main.go`.
4. Configure the worker with database, RabbitMQ, and Resend secrets in a systemd `EnvironmentFile` or secret manager.
5. Start or restart the worker, then deploy the API instances using the existing rolling process.
6. Register a test account and verify an outbox record is published, the queue drains, and the email arrives.

Example systemd unit:

```ini
[Unit]
Description=Mobo asynchronous email worker
After=network-online.target
Wants=network-online.target

[Service]
User=ubuntu
WorkingDirectory=/home/ubuntu/production
EnvironmentFile=/etc/mobo/email-worker.env
ExecStart=/home/ubuntu/production/mobo-email-worker
Restart=on-failure
RestartSec=5

[Install]
WantedBy=multi-user.target
```

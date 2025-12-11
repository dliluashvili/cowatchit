# CoWatchIt

A real-time collaborative watching room application built with Go, HTMX, and WebSockets.

## Features

- Real-time video synchronization
- WebSocket-based communication
- User authentication and session management
- Room creation and management
- Video streaming with Backblaze B2 integration
- PostgreSQL database with GORM
- Redis for caching and session storage
- Rate limiting and middleware protection

## Tech Stack

### Backend
- **Go 1.24.1** - Server-side logic
- **Chi Router** - HTTP routing
- **GORM** - ORM for PostgreSQL
- **Redis** - Caching and session management
- **WebSocket** - Real-time bidirectional communication
- **Templ** - Type-safe Go templating
- **Validator** - Request validation

### Frontend
- **TypeScript** - Type-safe JavaScript
- **HTMX** - Modern HTML interactions
- **Tailwind CSS** - Utility-first styling
- **DaisyUI** - Component library
- **Video.js** - HTML5 video player
- **Vite** - Build tool and dev server
- **Slim Select** - Custom select dropdowns

### Infrastructure
- **PostgreSQL** - Primary database
- **Redis** - Cache and session store
- **Backblaze B2** - Object storage for videos

## Project Structure

```
.
├── cmd/
│   ├── server/      # Main application server
│   └── migrator/    # Database migration tool
├── internal/
│   ├── handlers/    # HTTP request handlers
│   ├── services/    # Business logic
│   ├── repositories/# Data access layer
│   ├── models/      # Database models
│   ├── dtos/        # Data transfer objects
│   ├── params/      # Request parameters
│   ├── templates/   # Templ templates
│   ├── middlewares/ # HTTP middlewares
│   ├── interceptors/# Request interceptors
│   ├── helpers/     # Utility functions
│   └── shared/      # Shared constants and validators
├── db/
│   ├── connection.go # Database connection
│   └── migrations/   # SQL migrations
└── static/
    ├── src/         # Frontend TypeScript/CSS
    └── dist/        # Compiled frontend assets
```

## Prerequisites

- Go 1.24.1 or higher
- PostgreSQL 14 or higher
- Redis 6 or higher
- Node.js 18 or higher (for frontend development)
- Backblaze B2 account (for video storage)

## Environment Variables

Create a `.env.dev` file in the root directory:

```env
APP_ENV=dev
REDIS_HOST=0.0.0.0
REDIS_PORT=6379
POSTGRES_HOST=0.0.0.0
POSTGRES_PORT=5432
POSTGRES_USER=postgres
POSTGRES_PASSWORD=your_password
POSTGRES_DB=cowatchit
POSTGRES_TEST_DB=cowatchittest
B2_APPLICATION_KEY_ID=your_key_id
B2_APPLICATION_KEY=your_key
B2_BUCKET_NAME=your_bucket_name
```

## Installation

### 1. Clone the repository

```bash
git clone https://github.com/dliluashvili/cowatchit.git
cd cowatchit
```

### 2. Install Go dependencies

```bash
go mod download
```

### 3. Install frontend dependencies

```bash
cd static
npm install
cd ..
```

### 4. Set up the database

Start PostgreSQL and Redis:

```bash
# Using Docker
docker run -d -p 5432:5432 -e POSTGRES_PASSWORD=your_password postgres:14
docker run -d -p 6379:6379 redis:6
```

Run migrations:

```bash
go run cmd/migrator/main.go
```

### 5. Build frontend assets

```bash
cd static
npm run build
cd ..
```

## Development

### Run the backend server

```bash
go run cmd/server/main.go
```

### Run frontend in development mode

In a separate terminal:

```bash
cd static
npm run dev
```

### With Air (live reload)

Install Air:

```bash
go install github.com/air-verse/air@latest
```

Run with Air:

```bash
air
```

## Building for Production

### Build backend

```bash
go build -o bin/server cmd/server/main.go
```

### Build frontend

```bash
cd static
npm run build
```

### Run production server

```bash
./bin/server
```

## Testing

Run tests:

```bash
go test ./...
```

Run tests with coverage:

```bash
go test -cover ./...
```

## API Endpoints

The application uses HTMX for most interactions. Key endpoints include:

- `/` - Home page
- `/rooms` - Room listing and creation
- `/rooms/:id` - Individual room view
- `/ws` - WebSocket connection for real-time features
- `/auth/login` - User login
- `/auth/register` - User registration
- `/auth/logout` - User logout

## Features in Detail

### Real-time Synchronization
WebSocket-based video synchronization ensures all viewers in a room see the same content at the same time.

### Rate Limiting
Built-in rate limiting (5 requests per second, burst of 20) protects against abuse.

### Session Management
Redis-backed session storage for fast, scalable authentication.

### Image Processing
Automatic image conversion to WebP format for optimal performance.

### Validation
Request validation using go-playground/validator ensures data integrity.

## Contributing

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## License

This project is licensed under the ISC License.

## Acknowledgments

- Built with [Templ](https://templ.guide/) for type-safe Go templates
- Uses [HTMX](https://htmx.org/) for modern HTML interactions
- Styled with [Tailwind CSS](https://tailwindcss.com/) and [DaisyUI](https://daisyui.com/)

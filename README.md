# House Marketplace Server

A robust backend server for a real estate marketplace application, built with Go. This server provides APIs for managing property listings, user authentication, file uploads, and real-time messaging between users.

## Features

- User authentication and authorization
- Property listing management
- Real-time messaging system
- File upload and storage
- RESTful API endpoints
- CORS support
- Database integration with PostgreSQL

## Technologies Used

- **Backend Framework**: [Gin](https://github.com/gin-gonic/gin) - High-performance HTTP web framework
- **Database**: [PostgreSQL](https://www.postgresql.org/) with [pgx](https://github.com/jackc/pgx) driver
- **Authentication**: JWT (JSON Web Tokens)
- **File Storage**: Google Cloud Storage
- **WebSocket**: [Gorilla WebSocket](https://github.com/gorilla/websocket) for real-time communication
- **Environment Management**: [godotenv](https://github.com/joho/godotenv)
- **Validation**: [go-playground/validator](https://github.com/go-playground/validator)
- **UUID Generation**: [google/uuid](https://github.com/google/uuid)

## Prerequisites

- Go 1.24.1 or higher
- PostgreSQL database
- Google Cloud Storage account (for file storage)
- Environment variables set up (see Configuration section)

## Project Structure

```
.
├── internal/          # Internal application code
│   ├── controller/   # HTTP handlers and routing
│   ├── repository/   # Data access layer
│   └── usecases/     # Business logic
├── pkg/              # Public packages
├── db/               # Database migrations and scripts
├── main.go           # Application entry point
└── .env              # Environment variables
```

## Setup and Installation

1. Clone the repository:
   ```bash
   git clone https://github.com/groggy7/house-marketplace-server
   cd house-marketplace-server
   ```

2. Install dependencies:
   ```bash
   go mod download
   ```

3. Set up environment variables:
   Create a `.env` file in the root directory with the following variables:
   ```
    DB_URL=
    FRONTEND_URL=
    JWT_SECRET=
    FIREBASE_CREDENTIALS=
    FIREBASE_BUCKET=
   ```

4. Set up the database:
   - Create a PostgreSQL database
   - Run the database migrations from the `db/` directory

5. Run the application:
   ```bash
   go run main.go
   ```

The server will start on the specified port (default: 8080).


## Contributing

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add some amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request
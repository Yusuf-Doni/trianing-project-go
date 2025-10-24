## Features

- ✅ Hello World web application
- ✅ PostgreSQL database connection
- ✅ RESTful API endpoints
- ✅ Beautiful responsive web interface
- ✅ User management (CRUD operations)
- ✅ Database table creation
- ✅ Sample data insertion
- ✅ Environment variable configuration
- ✅ Clean Architecture (Controller, Service, Model, Routes)
- ✅ Dependency Injection
- ✅ Business logic validation
- ✅ API documentation
- ✅ Monolith architecture
- ✅ **User Authentication & Authorization**
- ✅ **Login/Logout functionality**
- ✅ **User Registration**
- ✅ **Session Management**
- ✅ **Protected Routes**

## Prerequisites

- Go 1.21 or higher
- PostgreSQL database server
- Git (optional)

## Installation & Setup

### 1. Clone the repository
```bash
git clone <your-repo-url>
cd trianing-project-go
```

### 2. Install dependencies
```bash
go mod tidy
```

### 3. Setup PostgreSQL Database

#### Option A: Using Docker (Recommended)
```bash
# Run PostgreSQL in Docker
docker run --name postgres-db -e POSTGRES_PASSWORD=password -e POSTGRES_DB=testdb -p 5432:5432 -d postgres:15
```

#### Option B: Local PostgreSQL Installation
1. Install PostgreSQL on your system
2. Create a database:
```sql
CREATE DATABASE testdb;
CREATE USER postgres WITH PASSWORD 'password';
GRANT ALL PRIVILEGES ON DATABASE testdb TO postgres;
```

### 4. Configure Environment Variables

Copy the example configuration file:
```bash
cp config.env.example .env
```

Edit the `.env` file with your database credentials:
```
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=your_password_here
DB_NAME=testdb
```

### 5. Run the Application

#### Option A: Using environment variables
```bash
# Windows PowerShell
$env:DB_HOST="localhost"
$env:DB_PORT="5432"
$env:DB_USER="postgres"
$env:DB_PASSWORD="password"
$env:DB_NAME="testdb"
go run main.go

# Windows Command Prompt
set DB_HOST=localhost
set DB_PORT=5432
set DB_USER=postgres
set DB_PASSWORD=password
set DB_NAME=testdb
go run main.go

# Linux/Mac
export DB_HOST=localhost
export DB_PORT=5432
export DB_USER=postgres
export DB_PASSWORD=password
export DB_NAME=testdb
go run main.go
```

#### Option B: Direct execution (uses default values)
```bash
go run main.go
```

## Web Application Usage

### 1. Access the Web Interface
Open your browser and go to: `http://localhost:8080`

### 2. API Endpoints

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/` | Web interface |
| GET | `/api/users` | Get all users |
| POST | `/api/users` | Create new user |
| GET | `/api/users/{id}` | Get user by ID |
| PUT | `/api/users/{id}` | Update user |
| DELETE | `/api/users/{id}` | Delete user |
| GET | `/health` | Health check |

### 3. Example API Usage

**Get all users:**
```bash
curl http://localhost:8080/api/users
```

**Create a new user:**
```bash
curl -X POST http://localhost:8080/api/users \
  -H "Content-Type: application/json" \
  -d '{"name":"Jane Doe","email":"jane@example.com"}'
```

**Get user by ID:**
```bash
curl http://localhost:8080/api/users/1
```

## Expected Output

### Console Output:
```
🚀 Hello World! Starting Go Web Application with PostgreSQL!
✅ Successfully connected to PostgreSQL database!
✅ Table 'users' created or already exists
✅ Sample data inserted successfully
🌐 Server starting on port 8080
📱 Web Interface: http://localhost:8080
🔗 API Endpoint: http://localhost:8080/api/users
❤️  Health Check: http://localhost:8080/health
📚 API Documentation: http://localhost:8080/api/docs
🎉 Hello World Web Application with Clean Architecture is ready!
```

### Web Interface Features:
- Beautiful responsive design
- Add new users
- View all users in a table
- Delete users
- Real-time updates
- Error handling and success messages

## Authentication System

### Default Login Credentials
- **Admin User**: username=`admin`, password=`admin123`
- **Test Users**: You can register new users through the registration page

### Authentication Features
- **Login Page**: `/login` - Secure user authentication
- **Registration Page**: `/register` - New user registration
- **Logout**: Automatic session termination
- **Session Management**: Secure session handling with expiration
- **Protected Routes**: All inventory management routes require authentication

### Authentication Flow
1. **Access Control**: All inventory routes (`/dashboard`, `/addproduct`, etc.) are protected
2. **Login Required**: Unauthenticated users are redirected to `/login`
3. **Session Validation**: Sessions are validated on each request
4. **Automatic Logout**: Sessions expire after 24 hours

### User Roles
- **Admin**: Full access to all features
- **User**: Standard user access
- **Manager**: Enhanced permissions (can be extended)

## Project Structure

```
trianing-project-go/
├── main.go                    # Main application file
├── go.mod                    # Go module file
├── config.env.example        # Environment configuration example
├── controllers/              # Controllers layer
│   └── user_controller.go    # HTTP controllers
├── services/                 # Services layer (Business Logic)
│   └── user_service.go       # Business logic and validation
├── models/                   # Models layer (Data Access)
│   └── user.go              # Database operations
├── routes/                   # Routes layer
│   └── routes.go            # Route definitions
├── templates/               # Templates layer
│   ├── index.html           # Web interface template
│   └── api-docs.html        # API documentation
├── static/                  # Static files
│   ├── css/
│   │   └── style.css        # Additional CSS styles
│   └── js/
│       └── app.js           # Additional JavaScript
└── README.md                # This file
```

## Database Schema

The application creates a `users` table with the following structure:

```sql
CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    email VARCHAR(100) UNIQUE NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
```

## Troubleshooting

### Common Issues

1. **Database Connection Failed**
   - Ensure PostgreSQL is running
   - Check database credentials
   - Verify database exists

2. **Port Already in Use**
   - Change the port in your configuration
   - Stop other PostgreSQL instances

3. **Permission Denied**
   - Ensure the database user has proper permissions
   - Check firewall settings

### Getting Help

If you encounter issues:
1. Check the error messages in the console
2. Verify your database connection settings
3. Ensure PostgreSQL is running and accessible

## Architecture Overview

### Clean Architecture Layers:

1. **Controllers** (`controllers/`)
   - Handle HTTP requests and responses
   - Input validation and error handling
   - Call services for business logic

2. **Services** (`services/`)
   - Business logic and validation
   - Data transformation
   - Business rules enforcement

3. **Models** (`models/`)
   - Database operations
   - Data access layer
   - Database queries

4. **Routes** (`routes/`)
   - Route definitions
   - URL mapping
   - Static file serving

## Next Steps

- Add authentication and authorization
- Implement unit tests for each layer
- Add database migrations
- Implement caching layer
- Add monitoring and metrics
- Create Docker containerization
- Add API rate limiting
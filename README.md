# Superbank App

This project is composed of two main components:
- **Frontend**: Built with [NextJS](https://nextjs.org/)
- **Backend**: Built with [Golang Fiber](https://gofiber.io/)

You can run the project using one of the two methods:
1. ✅ **Manual Setup** (run frontend & backend separately)
2. 🐳 **Docker Setup** (via `docker-compose`)

---

## Project Structure
```
project-root/
├── superbank-backend/            # Golang backend service
├── superbank-frontend/           # NextJS frontend app
└── docker-compose.yml            # Docker Compose setup
```

---

## Manual Setup

### ***Backend (Golang)***

#### Prerequisites:
- Go 1.24.1 installed
- Go gotest.tools/gotestsum@latest installed
- PostgreSQL or your preferred database running

#### Steps:

```bash
- cd superbank-backend

# Install dependencies
- go mod tidy

# Create environment
- create .env and customize the contents like .env.example

# Run Migration
- go run main.go migrate

# Run Seeder
- go run main.go seed

# Run go server
- go run main.go
```

### Default user for login:
```json
	{
		"email":    "admin@example.com",
		"password": "password",
	},
	{
		"email":    "employee@example.com",
		"password": "password",
	},
```

### Unit Test
- MacOS or Linux:
```bash
#Running unit test
./test.sh
```
- Windows:
```bash
# Create a directory to store test results if it doesn't exist
mkdir -p test

# Run Go tests with coverage for all packages under ./module
go test ./module/... -coverprofile test/coverage-full.out -covermode atomic -coverpkg ./...

# Copy the header line from the coverage file to the filtered file
head -n 1 test/coverage-full.out > test/coverage.out

# Exclude specific files (DTOs, interfaces, containers) from the coverage report
grep 'module/' test/coverage-full.out | \
  grep -vE '\.dto\.go|\.interface\.go|\.container\.go' >> test/coverage.out

# Generate an HTML coverage report
go tool cover -html=test/coverage.out -o test/index.html

```
- Open HTML Report
> open test folder
>> copy path index.html
>>> paste on browser

---

### ***Frontend (NextJS)***

### Prerequisites:
- Node 23.6.1 installed
- npm or yarn

### Steps:
```bash
cd frontend

# Install dependencies
npm install

# Run the frontend app
npm run dev
```

---

## Docker Setup (Recommended)

### Prerequisites:
- Docker & Docker Compose installed

### Steps:
```bash
docker-compose up --build -d
```

---

### Notes
- Feel free to modify the docker-compose.yml file to adjust ports, volumes, or environment variables.
- If you make code changes and want to rebuild the containers: 
```bash
docker-compose down --volumes --remove-orphans
```
- When docker is being built in it will run a unit test which if the unit test is below threshold 95% then the build will be declared failed.
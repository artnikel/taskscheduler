# Task Scheduler

A lightweight HTTP-based task scheduler written in Go. It allows scheduling and monitoring of asynchronous tasks such as pinging hosts or checking HTTP status codes.

## Features

- Schedule ping tasks (`tcp` to port 80).
- Schedule HTTP status check tasks.
- Monitor task status and results via REST API.
- Configurable concurrency via YAML config.
- Basic logging to file.

## Installation

1. **Clone the repository:**

```bash
git clone https://github.com/artnikel/taskscheduler.git
cd taskscheduler
```

2. **Build the project:**
   
```bash
go build -o taskscheduler ./cmd/taskscheduler
```

3. **Prepare the configuration:**

Create a `config.yaml` file in the root directory.

Example `config.yaml`:

```yaml
server:
  port: 8080

logging:
  path: "logs"

scheduler:
  max_concurrent_tasks: 3

worker:
  ping_sites:
    - "google.com"
    - "yahoo.com"
```

4. **Run the server:**

```bash
./taskscheduler
```

You should see the server running at: `http://localhost:8080`


## API Endpoints

### 1. Create Ping Task
- **URL:** `/tasks/ping`
- **Method:** `POST`
- **Description:** Schedules a TCP ping task to check if a host is reachable on port 80.
- **Request Body:**
  ```json
  {
    "address": "example.com"
  }
  ```
- **Response:**
  ```json
  {
    "task_id": "your-generated-task-id"
  }
  ```

### 2. Create HTTP Status Task
- **URL:** `/tasks/http/status`
- **Method:** `POST`
- **Description:** Schedules an HTTP GET request to the provided URL.
- **Request Body:**
  ```json
  {
    "url": "https://example.com"
  }
  ```
- **Response:**
  ```json
  {
    "task_id": "your-generated-task-id"
  }
  ```

### 3. Get Task Status
- **URL:** `/tasks/{id}`
- **Method:** `GET`
- **Description:** Returns the status and result/error of a specific task.
- **Response (example):**
  ```json
  {
    "id": "task-id",
    "status": "done",
    "result": "ping example.com success, time: 200ms"
  }
  ```

### 4. Get Statistics
- **URL:** `/tasks/stats`
- **Method:** `GET`
- **Description:** Returns a summary of tasks grouped by their status.
- **Response:**
  ```json
  {
    "pending": 1,
    "running": 0,
    "done": 3,
    "failed": 1
  }
  ```
## Logging

The application uses structured logging to track important events and errors. Logs are written to a file configured in the application settings. There are two main loggers:

- **Info logger**: Records informational messages about normal operations.
- **Error logger**: Records errors and warnings to help diagnose issues.

The log file `app.log` is created immediately after the server starts, located in the `logs` directory by default.

## Testing


```bash
make test
```

## Linter


```bash
make lint
```


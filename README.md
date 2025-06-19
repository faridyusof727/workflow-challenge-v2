# ⚡ Workflow Editor

A modern workflow editor app for designing and executing custom automation workflows (e.g., weather notifications). Users can visually build workflows, configure parameters, and view real-time execution results.

## 🛠️ Tech Stack

- **Frontend:** React + TypeScript, @xyflow/react (drag-and-drop), Radix UI, Tailwind CSS, Vite
- **Backend:** Go API, PostgreSQL database
- **DevOps:** Docker Compose for orchestration, hot reloading for rapid development

## 🚀 Quick Start

### Prerequisites

- Docker & Docker Compose (recommended for development)
- Node.js v18+ (for local frontend development)
- Go v1.23+ (for local backend development)

> **Tip:** Node.js and Go are only required if you want to run frontend or backend outside Docker.

### 1. Start All Services

```bash
docker-compose up --build
```

- This launches frontend, backend, and database with hot reloading enabled for code changes.
- To stop and clean up:

  ```bash
  docker-compose down
  ```

### 2. Access Applications

- **Frontend (Workflow Editor):** [http://localhost:3003](http://localhost:3003)
- **Backend API:** [http://localhost:8086](http://localhost:8086)
- **Database:** PostgreSQL on `localhost:5876`

### 3. Verify Setup

1. Open [http://localhost:3003](http://localhost:3003) in your browser.
2. You should see the workflow editor with sample nodes.

## 🏗️ Project Architecture

```text
workflow-challenge-v2/
├── api/                    # Go Backend (Port 8086)
│   ├── main.go             # Entry point for the backend application
│   ├── internal/           # Internal logic and services
│   │   └── workflow/       # Workflow related logic
│   ├── pkg/                # Packages for configuration, DI, helpers, mailer, nodes
│   │   ├── config/         # Configuration files
│   │   ├── di/             # Dependency Injection setup
│   │   └── nodes/          # Implementations for different node types
│   ├── Dockerfile          # Dockerfile for building the backend image
│   └── README.md           # API documentation
├── web/                    # React Frontend (Port 3003)
│   ├── src/                # Source code for the frontend
│   │   ├── components/     # React components
│   │   ├── hooks/          # Custom React hooks
│   │   ├── App.tsx         # Main application component
│   │   ├── index.css       # Global styles
│   │   └── main.tsx        # Entry point for the frontend application
│   ├── public/             # Public assets
│   ├── package.json        # NPM dependencies and scripts
│   ├── vite.config.ts      # Vite configuration
│   ├── Dockerfile          # Dockerfile for building the frontend image
│   └── README.md           # Frontend documentation
├── docker-compose.yml      # Docker Compose file for defining and running multi-container Docker applications
├── DESIGN_RATIONALE.md     # Design rationale and architectural decisions
└── README.md               # Project documentation
```

## 🔧 Development Workflow

### 🌐 Frontend

- Edit files in `web/src/` and see changes instantly at [http://localhost:3003](http://localhost:3003) (hot reloading via Vite).
- **IMPORTANT** - Workflow ID is hardcoded. You'll need to replace it to correct ID from the database table `workflow`.`id`.

### 🖥️ Backend

- Edit files in `api/` and changes are reflected automatically (hot reloading in Docker).
- If you add new dependencies or make significant changes, rebuild the API container:

  ```bash
  docker-compose up --build api
  ```

### 🗄️ Database

- Schema/configuration details: see [API README](api/README.md#database)
- After schema changes or migrations, restart the database:

  ```bash
  docker-compose restart postgres
  ```

- To apply schema changes to the API after updating the database:

  ```bash
  docker-compose restart api
  ```

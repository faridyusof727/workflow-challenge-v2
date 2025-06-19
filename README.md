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
│   ├── main.go
│   ├── cmd/
│   │   ├── api.go
│   │   └── root.go
│   ├── api/
│   │   ├── middlewares.go
│   │   ├── route.go
│   │   └── server.go
│   ├── internal/
│   │   ├── edge/
│   │   │   └── types.go
│   │   ├── node/
│   │   │   └── types.go
│   │   └── workflow/
│   │       ├── handler.go
│   │       ├── port.go
│   │       ├── repository.go
│   │       ├── service.go
│   │       └── types.go
│   ├── pkg/
│   │   ├── config/
│   │   │   ├── config.go
│   │   │   ├── cors.go
│   │   │   └── db.go
│   │   ├── di/
│   │   │   ├── db.go
│   │   │   ├── interfaces.go
│   │   │   ├── logger.go
│   │   │   ├── nodes.go
│   │   │   └── service.go
│   │   ├── helper/
│   │   │   └── filter.go
│   │   ├── mailer/
│   │   │   ├── interfaces.go
│   │   │   └── noop.go
│   │   ├── nodes/
│   │   │   ├── README.md
│   │   │   ├── condition/
│   │   │   │   ├── executor.go
│   │   │   │   ├── executor_test.go
│   │   │   │   ├── replacer.go
│   │   │   │   └── types.go
│   │   │   ├── email/
│   │   │   │   ├── executor.go
│   │   │   │   ├── executor_test.go
│   │   │   │   └── replacer.go
│   │   │   ├── form/
│   │   │   │   ├── executor.go
│   │   │   │   └── executor_test.go
│   │   │   ├── service.go
│   │   │   ├── types/
│   │   │   │   └── executor.go
│   │   │   └── weatherapi/
│   │   │       ├── executor.go
│   │   │       └── executor_test.go
│   │   ├── openstreetmap/
│   │   │   ├── client.go
│   │   │   ├── client_test.go
│   │   │   ├── interfaces.go
│   │   │   └── types.go
│   │   ├── openweather/
│   │   │   ├── client.go
│   │   │   ├── client_test.go
│   │   │   ├── interfaces.go
│   │   │   └── types.go
│   │   ├── postgres/
│   │   │   ├── migrations/
│   │   │   │   ├── 20250618035144_create_workflow_table.sql
│   │   │   │   ├── 20250618093549_create_nodes_table.sql
│   │   │   │   ├── 20250618094125_create_workflow_nodes_table.sql
│   │   │   │   ├── 20250618100526_create_edges_table.sql
│   │   │   │   └── 20250618101127_seed_workflows.sql
│   │   │   └── service.go
│   │   └── render/
│   │       ├── errors.go
│   │       └── response.go
│   ├── go.mod
│   ├── go.sum
│   ├── Dockerfile
│   ├── Dockerfile.migrator
│   └── README.md
├── web/                    # React Frontend (Port 3003)
│   ├── src/
│   │   ├── components/
│   │   │   ├── ExecutionResults.tsx
│   │   │   ├── UserInputForm.tsx
│   │   │   ├── WorkflowDiagram.tsx
│   │   │   └── WorkflowNode.tsx
│   │   ├── hooks/
│   │   │   ├── useExecuteWorkflow.ts
│   │   │   └── useWorkflow.ts
│   │   ├── App.tsx
│   │   ├── constants.ts
│   │   ├── index.css
│   │   ├── main.tsx
│   │   ├── types.ts
│   │   └── vite-env.d.ts
│   ├── public/
│   │   └── checkbox.ico
│   ├── package.json
│   ├── package-lock.json
│   ├── vite.config.ts
│   ├── tsconfig.json
│   ├── tsconfig.app.json
│   ├── tsconfig.node.json
│   ├── tailwind.config.js
│   ├── postcss.config.js
│   ├── eslint.config.mjs
│   ├── nginx.conf
│   ├── Dockerfile
│   ├── README.md
│   └── index.html
├── docker-compose.yml
├── DESIGN_RATIONALE.md
└── README.md
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

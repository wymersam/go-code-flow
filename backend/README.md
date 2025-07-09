# 🧠 Code Flow Visualiser & Function Summariser (React + D3 + Go)

This project lets you upload Go source files, parses them to extract function call relationships, and visualises the resulting call graph interactively in a React frontend using D3.js.

---

## 🚀 Features

- Upload .go source files via React UI
- Backend Go service parses code to build function call graph
- Interactive force-directed graph visualisation with D3.js
- Optional function summaries generated via OpenAI API (TODO)

---

## 🛠️ Setup & Installation

### Backend (Go)

#### 1. Clone the repo

```bash
git clone https://github.com/your-username/go-code-graph
cd go-code-graph
```

#### 2. Install dependencies

```bash
go mod tidy
```

#### 3. Create a .env file in the root directory

```bash
OPENAI_API_KEY=sk-proj-your-key-here
```

#### Step 4: Run the backend server

```bash
cd backend
go run main.go
# By default it will listen on http://localhost:8080
```

with summaries:

```bash
go run main.go -summaries ./example main
```

### Frontend (React + D3)

#### Step 1: Navigate to the frotnend directory

```bash
cd frontend
```

#### Step 2: Install dependencies

```bash
npm install
```

#### Step 3: Start the React development server

```bash
npm run dev
```

Then open your browser to `http://localhost:5173/`

## Usage

- Upload a `.go` file using the React interface.
- The frontend sends the file content to the Go backend `/parse` endpoint.
- The backend parses the Go file to identify functions and call relationships.
- Backend responds with JSON describing the graph nodes and edges.
- React + D3 renders an interactive force-directed graph where you can explore the call structure.

## API details

### POST /parse

Request: JSON body containing:

```json
{
  "fileContent": "<Go source code as string>"
}
```

Response: JSON graph data:

```json
{
  "nodes": [{ "id": "FunctionA" }, { "id": "FunctionB" }],
  "links": [{ "source": "FunctionA", "target": "FunctionB" }]
}
```

## Dependencies

[OpenAI API](https://platform.openai.com/)

Go 1.18+

React 18+

D3.js

## Author

Sammy-Jo Wymer :)

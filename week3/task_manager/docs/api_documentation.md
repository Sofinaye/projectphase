# Task Management REST API (Go + Gin + MongoDB)

Base URL: `http://localhost:8080`

This API manages tasks with persistent storage using **MongoDB** and the **Mongo Go Driver**.

---

## MongoDB Configuration

- Default connection URI: `mongodb://localhost:27017`
- Override via environment variable: `MONGO_URI`
- Database name: `task_manager_db`
- Collections:
  - `tasks` — stores task documents.
  - `counters` — stores sequence counters, e.g.:
    ```json
    { "_id": "task_id", "seq": 1 }
    ```

### Running MongoDB

- Locally: via Docker  
  ```bash
  docker run -d --name mongo -p 27017:27017 mongo

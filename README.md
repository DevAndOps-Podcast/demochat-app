# Chat Application

This is a simple chat application built with Go and Echo framework for the backend, and plain HTML, CSS, and JavaScript for the frontend.

## Configuration

The application can be configured using a `config.yaml` file in the root directory. An example configuration is shown below:

```yaml
address: ":8081"
jwt_secret: "supersecretjwtkey"
log_level: "error"
debug: true
insights_service:
  base_url: "http://localhost:8082"
  api_key: "insights-api-key"
postgres:
  host: "localhost"
  port: 5432
  user: "postgres"
  password: "mysecretpassword"
  dbname: "postgres"
  sslmode: "disable"
  schema: "demochat"
```

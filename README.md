# Chat Application

This is a simple chat application built with Go and Echo framework for the backend, and plain HTML, CSS, and JavaScript for the frontend.

## Configuration

The application can be configured using a `config.yaml` file in the root directory. An example configuration is shown below:

```yaml
address: ":8080"
jwt_secret: "your_super_secret_jwt_key"
debug: true
insights_service:
  base_url: "http://localhost:8081"
  api_key: "your_insights_api_key"
```

- `address`: The address and port the server will listen on (e.g., `:8080`).
- `jwt_secret`: A secret key used for signing JWT tokens. **Change this to a strong, unique secret in production.**
- `debug`: A boolean indicating whether debug mode is enabled.
- `insights_service`: Configuration for the insights service.
  - `base_url`: The base URL of the insights service.
  - `api_key`: The API key for authenticating with the insights service.

## Getting Started

To run the application:

1.  **Build the application:**
    ```bash
    go build -o demochat-app .
    ```
2.  **Run the application:**
    ```bash
    ./demochat-app
    ```

The application will start on the address specified in `config.yaml` (default: `:8080`).

## Configuration

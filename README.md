# Chat Application

This is a simple chat application built with Go and Echo framework for the backend, and plain HTML, CSS, and JavaScript for the frontend.

## Configuration

The application can be configured using a `config.yaml` file in the root directory. An example configuration is shown below:

```yaml
address: ":8080"
jwt_secret: "your_super_secret_jwt_key"
debug: true
```

- `address`: The address and port the server will listen on (e.g., `:8080`).
- `jwt_secret`: A secret key used for signing JWT tokens. **Change this to a strong, unique secret in production.**
- `debug`: A boolean indicating whether debug mode is enabled.

## Getting Started

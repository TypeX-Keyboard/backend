# TypeX Keyboard Backend

This is the backend service for **TypeX Keyboard**, built with [GoFrame](https://goframe.org/) in Golang and deployed via Docker.

## Project Structure

```
├── manifest/
│   ├── config.yaml           # Main configuration file (required before deployment)
│   └── config.yaml.example   # Example configuration
├── Dockerfile                # Docker build file
├── go.mod                    # Go modules
├── main.go                   # Entry point
└── ...
```

## Getting Started

### 1. Clone the Repository

```bash
git clone https://github.com/your-username/typex-keyboard-backend.git
cd typex-keyboard-backend
```

### 2. Configure the Application

Before running the service, you need to set up the configuration:

1. Navigate to the `manifest` directory:
   ```bash
   cd manifest
   ```

2. Copy the example config file:
   ```bash
   cp config.yaml.example config.yaml
   ```

3. Edit `config.yaml` and update the required settings (e.g., database, server port, etc.).

> ✅ **Make sure the server listens on port 80**. Example:
> ```yaml
> server:
>   address: ":80"
> ```

### 3. Build and Run with Docker

Make sure Docker is installed on your system.

Build the Docker image:

```bash
docker build -t typex-backend .
```

Run the container and expose port 80:

```bash
docker run -d -p 80:80 --name typex-backend typex-backend
```

### 4. Access the Service

Once running, the service will be available at:

```
http://localhost
```

(No port number needed since it runs on port 80)

## Tech Stack

- [Golang](https://golang.org/)
- [GoFrame](https://goframe.org/)
- [Docker](https://www.docker.com/)

## Contributing

Feel free to open issues or submit pull requests to help improve this project.

## License

[MIT](LICENSE)

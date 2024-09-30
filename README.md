# UniCache

**UniCache** is an in-memory caching service that acts as an intermediary between the client and a backend server. It improves system efficiency by temporarily storing responses, reducing the load on the backend and speeding up responses to clients.

## Overview

UniCache functions as a middle layer between the client and the backend. When the client sends a request, UniCache checks if the response is already cached. If so, it returns the cached response. Otherwise, it forwards the request to the backend, stores the response, and returns it to the client.

```bash
Client -> UniCache -> Backend
```

This process reduces the need for repeated queries to the backend, improving overall system performance.

## Example Workflow

### Before using UniCache

The client would make a direct request to the backend:

```bash
GET http://192.168.0.100:8080/v1/product
```

### After adding UniCache

With UniCache acting as an intermediary, the client sends the request to the cache, which decides whether to forward it to the backend:

```bash
GET http://192.168.0.50:3030/v1/product
```

If the response is already cached, UniCache returns it directly. Otherwise, the request is forwarded to the backend at `http://192.168.0.100:8080/v1/product`, and the response is stored for future use.

## Configuration

### Step 1: Build the Project

First, compile the project by running the following command in the root directory:

```bash
go build
```

This will generate the executable `unicache`.

### Step 2: Configure UniCache

After compiling, you need to configure UniCache to point to the correct backend. This is done through environment variables:

- `POINT_ADDRESS`: The IP address of the backend.
- `POINT_PORT`: The port on which the backend is listening.
- `POINT_PROTOCOL`: The protocol (e.g., `http` or `https`) used by the backend.

Example configuration for a backend located at `192.168.0.100:8080`:

```bash
POINT_ADDRESS=192.168.0.100 POINT_PORT=8080 POINT_PROTOCOL=http ./unicache
```

This sets up UniCache to forward requests to the backend running at `192.168.0.100` on port `8080`.

### Step 3: Running the UniCache

Once configured, you can start UniCache. Now, clients can send requests to `192.168.0.50:3030`, where UniCache is running, and it will decide whether to respond from the cache or retrieve the response from the backend.

## Benefits

- **Improved Performance**: By caching responses in memory, UniCache reduces the backend's workload.
- **Reduced Latency**: Response time is decreased for cached requests.
- **Scalability**: Fewer requests to the backend allow the system to scale more effectively.

## Contribution

Contributions and suggestions are always welcome! If you have any ideas for improvement or encounter issues, feel free to open an issue or submit a pull request.

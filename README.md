# Go Gate

Gateway service for authenticating and routing microservice applications written in Go.

## Configuration

Configuration of the gateway service is done through the `config/app.yaml` file.
This file is the default configuration file for the service and can be overridden
when starting the gateway service with the `--config` flag on the application.

### General

- `port`: The port number to start and run the service on.

### Route

The list of routes that are avaliable and which service to send them to. The path
property allows for regular expressions to match routes against.

- `path`: The incoming request path to match a user against.
- `service`: The backend service that the request should be sent to for processing.
- `optionalAuth`: (default: `false`) Should a user be required to authenticate or should they be allowed through even without authentication.
- `stripPrefix`: (default: `false`) A string containing the prefix to strip off of requests as they are directed to the backend.

Example:

```yaml
routes:
  - path: /recipes/.*
    service: recipes
    optionalAuth: true
  - path: /images/.*
    service: images
    stripPrefix: /images
```

### Service

The list of backend services that routes will call when a request is made.

- `name`: The name of the servcie. This is used to match the `service` property of the route.
- `url`: The base URL that the service can be reached at. Must be a fully qualified URL (i.e `http://...`)

Example:

```yaml
services:
  - name: images
    url: http://localhost:3100
  - name: recipes
    url: http://localhost:3101
```

### Custom Errors

There are three different errors that the gateway service will serve itself.

These different errors are:

- **Not Found**: A matching route or service could not be found.
- **Unauthorized**: A user was not authenticated because the user module returned an error of some kind.
- **Service Unavaliable**: The backend service could not be contacted to send the request. (i.e. the HTTP client returned an error when trying to send the request)

Example (using the default values):

```yaml
errors:
  notFound:
    status: 404
    short: Not Found
    long: Could not find route
  unauthorized:
    status: 401
    short: Authentication Failed
    long: Could not find or process the authentication
  serviceUnavaliable:
    status: 502
    short: Could Not Process Request
    long: The server was unable to process your request
```

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

- **Not Found:** A matching route or service could not be found.
- **Service Unavaliable:** The backend service could not be contacted to send the request. (i.e. the HTTP client returned an error when trying to send the request)

Example (using the default values):

```yaml
errors:
  notFound:
    status: 404
    short: Not Found
    long: Could not find route
  serviceUnavaliable:
    status: 502
    short: Could Not Process Request
    long: The server was unable to process your request
```

### Plugins

Plugins allow intercepting requests to backend services and responses back to
clients. They are able to access the request and response during the Pre-Request
stage and the Post-Request stage to allow for blocking requests and/or sending
custom statuses back to users.

Plugins can hook into 3 lifecycle events:

- **Setup:** This function is called during initialization of the server. It will be called before any requests are able to be taken from clients. General setup tasks and heavy, non-request based logic should be placed here.
- **PreRequest:** This function is called before the request from a client gets sent to the backend service responsible for handling that request. Logic such as authentication should be placed here.
- **PostRequest:** This function is called after the backend service has responded but before the client has recieved the response. This can be used to intercept a response before it goes out.

Example configuration:

```yaml
plugins:
  # Sample Plugin
  - path: lib/sample.so
  # Block Headers During Request Lifecycle
  - path: lib/head_block.so
    settings:
      inbound:
        - X-Some
      outbound:
        - X-UserId
  # Handle JWT authentication
  - path: lib/jwt.so
    settings:
      domain: https://my.auth0.com/
      jwksUrl: https://my.auth0.com/.well-known/jwks.json
```

## Creating a Plugin

Plugins get built as Go plugins which are compiled C libraries. There is an
example of a plugin at [plugins/sample/sample.go](plugins/sample/sample.go)
in this project.

To build a plugin run:

```
go build -buildmode=plugin -o lib/my_plugin.so my_plugin/*
```

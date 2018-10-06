# Go Gate

Gateway service for authenticating and routing microservice applications written in Go.

<!-- START doctoc generated TOC please keep comment here to allow auto update -->
<!-- DON'T EDIT THIS SECTION, INSTEAD RE-RUN doctoc TO UPDATE -->
**Table of Contents**

- [Configuration](#configuration)
  - [General](#general)
  - [Route](#route)
  - [Service](#service)
  - [Custom Errors](#custom-errors)
  - [Plugins](#plugins)
- [Built-In Plugins](#built-in-plugins)
  - [Head Block](#head-block)
  - [JWT Plugin](#jwt-plugin)
  - [Sample Plugin](#sample-plugin)
- [Creating a Plugin](#creating-a-plugin)
  - [Core](#core)
  - [Setup Function](#setup-function)
  - [Pre-Request Function](#pre-request-function)
  - [Post-Request Function](#post-request-function)
  - [Building a Plugin](#building-a-plugin)

<!-- END doctoc generated TOC please keep comment here to allow auto update -->

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

## Built-In Plugins

There are several built-in plugins that serve as a starting point for extending
the gateway service.

### Head Block

The Head Blocker plugin allows headers for the request going to the backend
service they were destined for to be removed. It also allows for removing
headers returning from the backend service going back to the client.

The configuration of this plugin gets done in the plugin's settings section of
the configuration file:

```yaml
plugins:
  - path: lib/head_block.so
    settings:
      inbound:
        - First-Header
        - Second-Header
      outbound:
        - Third-Header
```

Where `inbound` are the headers coming from the client and going to the backend
service and `outbound` are the headers coming from the backend and heading back
to the client.

### JWT Plugin

The JWT plugin is responsible for handling user authentication using JWT and
a JWKS endpoint that can serve up the x509 public certificate for the JWTs that
were generated. This is built off of how [Auth0](https://auth0.com) handles
their JWT validation processes.

It allows for adding the `optionalAuth` paramter to the routes which will allow
them to not reject unauthenticated users. The default behavior is to **require**
authentication on all routes. They must explicitly set `optinalAuth` to `true`
to allow unauthenticated users through.

### Sample Plugin

This plugin is designed to demonstrate the functionality of plugins and give
some quick guidance on how to create them.

## Creating a Plugin

Plugins get built as Go plugins which are compiled C libraries. There is an
example of a plugin at [plugins/sample/sample.go](plugins/sample/sample.go)
in this project.

### Core

Go requires that all plugins be in their own `main` package and have the
exported members that they want accessed by the application that consumes them.

There is an internal context that is used on all plugins that allows for data
to be shared between the functions of a plugin. All lifecycle functions will
need to be on this struct in order to be read by the application.

The context that all of the functions are built on will need to be exported
with the name `Plugin`.

```go
package main

type samplePlugin struct{}

var Plugin struct{}
```

### Setup Function

The setup function takes in the settings for the plugin from the application
configuration file. It is primarily used to set up the plugin and do any heavy
lifting before it starts to take on any requests.

This function will need to take a single parameter of type `interface{}` and be
on the plugin context or the plugin manager will not load the `Setup` function.

```go
func (*samplePlugin) Setup(settings interface{}) {
	log.Println("Do some set up work here...")
}
```

**Note:** This function is not required for a plugin to work.

### Pre-Request Function

This function will need to take three parameters of types `http.ResponseWriter`,
`*http.Request`, and `*settings.RouteSettings` and have a return of `error`
as well as be located on the plugin context or the plugin manager will not
load the `PreRequest` function.

If there is an error returned the request will be stopped from reaching the
backend server that the request was headed to.

```go
func (*samplePlugin) PreRequest(w http.ResponseWriter, r *http.Request, route *settings.RouteSettings) error {
	log.Println("In Pre Request...")

	return nil
}
```

**Note:** This function is not required for a plugin to work.

### Post-Request Function

This function will need to take three parameters of types `http.ResponseWriter`,
`*http.Request`, and `*settings.RouteSettings` and have a return of `error`
as well as be located on the plugin context or the plugin manager will not
load the `PostRequest` function.

If there is an error returned the request will be stopped returning to the client.

```go
func (*samplePlugin) PostRequest(w http.ResponseWriter, r *http.Request, route *settings.RouteSettings) error {
	log.Printf("Running sample post request for %s\n", r.URL.Path)

	return nil
}
```

**Note:** This function is not required for a plugin to work.

### Building a Plugin

To build a plugin run:

```
go build -buildmode=plugin -o lib/my_plugin.so my_plugin/*
```

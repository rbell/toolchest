# Server Package

The `server` package is a starting point for setting up and managing HTTP and/or gRPC services.

## Getting Started

The example folder contains a simple example of an HTTP server that listens on port 8080 and a GRPC server that listens on port 8888. Configuring the server is done by building up a configuration and passing it to the function that creates the server.  See the full example for more details.

## Features
### Current Features
- Serve both HTTP and gRPC services from single server.
- Ability to configuration via configuration builder with chaining.

### Todo
- TLS support
- Support for configuring middleware
- Support for configuring health checks (with default implementation)
- Support for configuring logging

## License

This Source Code Form is subject to the terms of the Apache Public License, version 2.0. If a copy of the APL was not distributed with this file, you can obtain one at https://www.apache.org/licenses/LICENSE-2.0.txt.
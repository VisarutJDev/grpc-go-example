# gRPC Go Example

This project is a simple example to demonstrate how to use gRPC with Golang. It includes setup instructions, a brief explanation of the project, and how to run it.

## Table of Contents
1. [Setup Instructions](#setup-instructions)
2. [Project Overview](#project-overview)

## Setup Instructions

To set up the project locally, follow these steps:

1. **Install Golang**: Ensure that you have Golang installed on your machine. You can download it from the [official Go website](https://golang.org/dl/).

2. **Install Protocol Buffers**: Download and install the Protocol Buffers compiler (`protoc`) for your operating system from the [official Protocol Buffers releases page](https://github.com/protocolbuffers/protobuf/releases).

3. **Install gRPC and Protocol Buffers plugins for Go**:
    ```sh
    go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
    go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
    ```

    Ensure that `$GOPATH/bin` is added to your `PATH` so that the `protoc-gen-go` and `protoc-gen-go-grpc` plugins are accessible.

4. **Clone the Repository**: Clone the repository to your local machine using:
    ```sh
    git clone https://github.com/VisarutJDev/grpc-go-example.git
    cd grpc-go-example
    ```

5. **Generate gRPC Code**: Generate the gRPC code from the `.proto` file using:
    ```sh
    protoc --go_out=. --go-grpc_out=. proto/*.proto
    ```

6. **Install Dependencies**: Install the required dependencies using:
    ```sh
    go mod tidy
    ```

7. **Run the Server**: Start the gRPC server using:
    ```sh
    go run server/main.go
    ```

8. **Run the Client**: In a separate terminal, run the gRPC client using:
    ```sh
    go run client/main.go
    ```

## Project Overview

This project demonstrates a simple gRPC service written in Golang. The main components of the project are:

- **Protocol Buffers Definition**: The `.proto` file located in the `proto/` directory defines the gRPC service, its methods, and the message types used by those methods. Protocol Buffers (Protobuf) is a language-neutral, platform-neutral, extensible way of serializing structured data.

- **gRPC Server**: The server implementation, located in `server/main.go`, includes the logic to handle incoming gRPC requests. It listens on a specified port and serves the defined gRPC service.

- **gRPC Client**: The client implementation, located in `client/main.go`, demonstrates how to call the gRPC service methods from a client application. It establishes a connection to the gRPC server and makes requests to it.

### Understanding gRPC with Golang

gRPC (gRPC Remote Procedure Calls) is a high-performance, open-source, universal RPC framework initially developed by Google. It leverages HTTP/2 for transport, Protocol Buffers as the interface description language, and it provides features such as authentication, load balancing, and more.

In this example, you'll learn:

- **How to define a gRPC service**: Using a `.proto` file, you define the service and its methods, as well as the request and response message types.
- **How to generate Go code from `.proto` files**: Using `protoc` with the `protoc-gen-go` and `protoc-gen-go-grpc` plugins, you generate the Go types and gRPC service code.
- **How to implement a gRPC server in Go**: You write the logic for the server-side handling of gRPC requests.
- **How to implement a gRPC client in Go**: You write the client-side code to connect to the server and call its methods.

This project serves as a foundation for building more complex gRPC-based microservices in Go.
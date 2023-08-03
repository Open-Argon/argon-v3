# Use Alpine as the base image
FROM golang:alpine

# Install isotope
FROM ugric/isotope:latest

# Set the Current Working Directory inside the container
WORKDIR /argon

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download all dependencies. Dependencies will be cached if the go.mod and go.sum files are not changed
RUN go mod download

# Copy the source from the current directory to the Working Directory inside the container
COPY ./src ./src

# Build the Go app
RUN go build -trimpath -ldflags="-s -w" -o bin/argon ./src

# make the binary executable
RUN chmod +x bin/argon

# add the binary to the path
ENV PATH="/argon/bin:${PATH}"
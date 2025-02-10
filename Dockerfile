# Use the official Golang image to build the Go app
FROM golang:1.20 AS build

# Set the current working directory in the container
WORKDIR /app

# Copy the Go module and download dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy the rest of the application code
COPY . .

# Build the Go application
RUN go build -o main .

# Expose the port the app will run on
EXPOSE 8081

# Run the Go application
CMD ["./main"]

# Use the official Golang image as the base image
FROM golang:1.20-alpine

# Set the working directory inside the container
WORKDIR /app

# Copy the Go module files and download dependencies
COPY go.mod go.sum ./
RUN go mod tidy

# Copy the source code into the container
COPY . .

# Build the Go application
RUN go build -o main .

# Expose the port the application will run on
EXPOSE 8080

# Command to run the Go application
CMD ["./main"]

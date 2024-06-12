# Start from the official Golang base image
FROM golang:1.22

# Set the Current Working Directory inside the container
WORKDIR /app

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download all dependencies. Dependencies will be cached if the go.mod and go.sum files are not changed
RUN go mod download

# Copy the source code into the container
COPY . .

# Set the working directory to the location of the main.go file
WORKDIR /app/cmd

# Build the Go app
RUN go build -o /app/main .

# Set the working directory back to /app
WORKDIR /app

# Expose port 8080 to the outside world
EXPOSE 8080

# Command to run the executable
CMD ["./main"]

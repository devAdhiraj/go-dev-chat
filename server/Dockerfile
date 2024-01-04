# Use an official Golang runtime as a parent image
FROM golang:1.21.0

# Set the working directory inside the container
WORKDIR /app

# Copy the local package files to the container's workspace
COPY go.mod go.sum /app/
RUN go mod download

COPY ./models/ /app/models
COPY ./api/ /app/api
COPY main.go .env /app/

# Build the Go app
RUN go build -o chatserver .

# Expose port 8080 to the outside world
EXPOSE 8085

# Command to run the executable
CMD ["go", "run", "main.go"]

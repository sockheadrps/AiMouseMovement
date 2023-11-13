FROM golang:1.21

# Set the working directory inside the container
WORKDIR /app

# Copy the local packages into the Docker image
COPY ./server /app

# Get dependencies
RUN go get

# Build the Go application inside the container
RUN go build -o bin .

# Set the entry point for the Docker container
CMD ["./bin"]
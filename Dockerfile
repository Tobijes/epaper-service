# Use the official Go image as a parent image
FROM golang:1.22

# Set the working directory in the container
WORKDIR /app

# Copy go.mod and go.sum to the working directory
COPY go.mod go.sum ./

# Download all dependencies
RUN go mod download

# Copy the rest of the source code to the working directory
COPY main.go .
COPY fonts fonts
COPY electricity electricity

# Build the Go app
RUN CGO_ENABLED=0 GOOS=linux go build -o /my-go-app

# Define the command to run your app using CMD which defines your runtime
CMD ["/my-go-app"]

# use official Golang image
FROM golang:1.25.0-alpine3.22

# set working directory
WORKDIR /app

# Copy the source code
COPY . .

# Download and install dependencies
RUN go get -d -v ./...

# Build the application
RUN go build -o main .

# Expose the port the app runs on
EXPOSE 8080 8081

# Command to run the executable
CMD ["./main"]
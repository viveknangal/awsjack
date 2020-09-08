FROM golang:latest


# Set the Current Working Directory inside the container
WORKDIR /app


# Copy go mod and other files
COPY . ./



# Build the Go app
RUN go build -o main .

# This container exposes port 8080 to the outside world
EXPOSE 8080


# Run the binary program produced by `go install`
ENTRYPOINT  ["./main"]

CMD  ["eu-west-1"]


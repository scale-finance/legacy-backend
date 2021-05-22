# Telling Docker to use this image
FROM golang:1.16.4-alpine3.13

# Make app directoy and copy all src files in it
RUN mkdir /app
ADD . /app
WORKDIR /app

# Install all dependencies
RUN go mod download

# Run build
RUN go build -o main ./cmd/api/main.go

# Expose intended port
EXPOSE 5000

# run the binary file
CMD ["/app/main"]

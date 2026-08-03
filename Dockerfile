FROM golang:1.21-alpine

# Install nodejs for the differential fuzzing
RUN apk add --no-cache nodejs npm

WORKDIR /app

# Copy everything
COPY . .

# Download dependencies
RUN go mod download

# Run the test suite as the default command
CMD ["go", "test", "-v", "./..."]

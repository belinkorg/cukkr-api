FROM golang:1.23-alpine

WORKDIR /app

# Install dependencies
RUN apk add --no-cache git make

# Copy go mod files
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Default command
CMD ["go", "test", "-v", "-cover", "./..."]
FROM golang:1.23

WORKDIR /app

# Install CompileDaemon
RUN go install github.com/githubnemo/CompileDaemon@latest

# Create tmp directory
RUN mkdir -p /app/tmp

# Copy go.mod and go.sum first
COPY go.mod go.sum ./
RUN go mod download

# Copy the rest of the code
COPY . .

EXPOSE 8080

# Run CompileDaemon
CMD CompileDaemon --build="go build -buildvcs=false -o ./tmp/main ." --command="./tmp/main"

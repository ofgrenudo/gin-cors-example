FROM golang:1.24
LABEL authors="Joshua Winters-Bronw"

WORKDIR /opt/gin-cors-example

# Install dependencies
COPY go.mod go.sum ./
RUN go mod download
COPY . .

# Build Executable
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags "-s -w" -o /opt/gin-cors-example/dist/gin-cors-example.x86.bin /opt/gin-cors-example/cmd/gin-example

# Run
CMD ["/opt/gin-cors-example/dist/gin-cors-example.x86.bin"]

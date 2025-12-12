FROM golang:1.24
LABEL authors="Joshua Winters-Bronw"

WORKDIR /opt/gin-cors-example

# Install dependencies
COPY go.mod go.sum ./
RUN go mod download
COPY *.go ./

# Build Executable
RUN CGO_ENABLED=0 GOOS=linux GOARCH="amd64" go build -o /opt/gin-cors-example/dist/gin-cors-example.x86.bin

# Run
CMD ["/opt/gin-cors-example/dist/gin-cors-example.x86.bin"]

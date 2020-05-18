# Dockerfile References: https://docs.docker.com/engine/reference/builder/

# Start from the latest golang base image
FROM golang:latest as builder

# Add Maintainer Info
LABEL maintainer="Curder <q.curder@gmail.com>"

# Set the Current Working Directory inside the container
WORKDIR /app

# Copy go mod and sum files
# COPY go.mod go.sum ./
COPY go.mod ./

# proxy setting
RUN go env -w GOPROXY=https://goproxy.cn,direct

# Download all dependencies. Dependencies will be cached if the go.mod and go.sum files are not changed
RUN go mod download

# Copy the source from the current directory to the Working Directory inside the container
COPY . .

# Build the Go app
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main .


######## Start a new stage from scratch #######
FROM alpine:latest

RUN apk --no-cache add ca-certificates

# add tzdata options, Fixed "unknown time zone Asia/Shanghai" Error
RUN apk add -U tzdata
RUN cp /usr/share/zoneinfo/Asia/Shanghai /etc/localtime

WORKDIR /root/

# Copy the Pre-built binary file from the previous stage
COPY --from=builder /app/main .

# Expose port 80 to the outside world
EXPOSE 80

# Command to run the executable
CMD ["./main"]

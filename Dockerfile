FROM golang:1.23-alpine as builder

# Set working directory to /app
WORKDIR /app

# Copy environment variables
COPY config.cold.json /app/config.cold.json
COPY config.hot.json /app/config.hot.json

# Copy source code to the container
COPY . .

# Install dependencies and build the Go binary
RUN go mod tidy
RUN go build -o user-management-service .

# Stage 2: Create a minimal runtime image
FROM alpine:latest

# Install necessary dependencies
RUN apk update && apk upgrade && \
    apk --update --no-cache add tzdata && \
    apk --no-cache add curl ca-certificates && \
    mkdir /app

# Set working directory to /app
WORKDIR /app

# Expose the port the application will run on
EXPOSE 8093

# Copy the built Go binary from the builder image
COPY --from=builder /app/user-management-service /app/

# Copy the config files from the builder image to /app
COPY --from=builder /app/config.cold.json /app/config.cold.json
COPY --from=builder /app/config.hot.json /app/config.hot.json
# TODO update "dbMysqlHost": from "localhost" to ip mysql container "172.18.0.2"
# Entry point to run the application
ENTRYPOINT ["/app/user-management-service"]

CMD []
FROM golang:1.23-alpine AS build

WORKDIR /src

# cache dependencies
COPY ./go.mod ./go.sum ./
RUN go mod download

# copy source code
COPY ./ ./

RUN CGO_ENABLED=0 GOOS=linux GOARCH=arm go build -o /app ./cmd/server

FROM --platform=linux/arm64/v8 scratch

COPY --from=build /app /app

# Copy the certs from the builder stage
COPY --from=build /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

ENTRYPOINT ["/app"]
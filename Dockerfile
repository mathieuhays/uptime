FROM golang AS build

WORKDIR /src

# cache dependencies
COPY ./go.mod ./go.sum ./
RUN go mod download

# copy source code
COPY ./ ./

RUN CGO_ENABLED=0 GOOS=linux GOARCH=arm go build -o /app ./cmd/server

FROM scratch

COPY --from=build /app /app

# Copy the certs from the builder stage
COPY --from=build /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

ENTRYPOINT ["/app"]
# Build Stage
FROM golang:1.21-alpine3.19 AS build-stage
WORKDIR /app
COPY ./ /app
RUN go mod download
RUN go build -o /app/build ./cmd/main.go

# Final Stage
FROM gcr.io/distroless/static-debian11
COPY --from=build-stage /app/build /
COPY --from=build-stage /app/.env /
EXPOSE 8000
CMD ["/build"]

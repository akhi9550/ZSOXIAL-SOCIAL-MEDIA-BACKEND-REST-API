FROM golang:1.21.2-alpine3.18 AS build-stage
WORKDIR /auth_svc
COPY ./ /auth_svc
RUN mkdir -p /auth_svc/build
RUN go mod download
RUN go build -v -o /auth_svc/build/api ./cmd


FROM gcr.io/distroless/static-debian11
COPY --from=build-stage /auth_svc/build/api /
COPY --from=build-stage /auth_svc/.env /
EXPOSE 50051
CMD ["/api"]
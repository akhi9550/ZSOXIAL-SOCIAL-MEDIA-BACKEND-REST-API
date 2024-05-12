FROM golang:1.21.2-alpine3.18 AS build-stage
WORKDIR /notification_svc
COPY ./ /notification_svc
RUN mkdir -p /notification_svc/build
RUN go mod download
RUN go build -v -o /notification_svc/build/api ./cmd


FROM gcr.io/distroless/static-debian11
COPY --from=build-stage /notification_svc/build/api /
COPY --from=build-stage /notification_svc/.env /
EXPOSE 50054
CMD ["/api"]
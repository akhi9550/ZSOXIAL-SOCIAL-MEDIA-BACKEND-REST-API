FROM golang:1.21.2-alpine3.18 AS build-stage
WORKDIR /post_svc
COPY ./ /post_svc
RUN mkdir -p /post_svc/build
RUN go mod download
RUN go build -v -o /post_svc/build/api ./cmd


FROM gcr.io/distroless/static-debian11
COPY --from=build-stage /post_svc/build/api /
COPY --from=build-stage /post_svc/.env /
EXPOSE 50052
CMD ["/api"]
FROM golang:1.21.2-alpine3.18 AS build-stage
WORKDIR /zsoxial
COPY ./ /zsoxial
RUN mkdir -p /zsoxial/build
RUN go mod download
RUN go build -v -o /zsoxial/build/api ./cmd


FROM gcr.io/distroless/static-debian11
COPY --from=build-stage /zsoxial/build/api /
COPY --from=build-stage /zsoxial/template/ /template/
COPY --from=build-stage /zsoxial/static /static/
COPY --from=build-stage /zsoxial/.env /
EXPOSE 8080
CMD ["/api"]


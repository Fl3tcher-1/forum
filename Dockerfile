FROM golang:1.18 as build
WORKDIR /app
COPY . /app
RUN CGO_ENABLED=1 go build -o /forum -a -ldflags '-linkmode external -extldflags "-static"' .

FROM scratch
WORKDIR /
COPY . ./
COPY --from=build /forum /forum
COPY --from=build /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
EXPOSE 8080
ENTRYPOINT ["/forum"]
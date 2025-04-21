FROM golang:1.24 AS build

WORKDIR /app

ENV PATH="/go/bin:${PATH}"

COPY docker/entrypoint.sh /docker/entrypoint.sh
RUN chmod +x /docker/entrypoint.sh

COPY go.mod go.sum ./
RUN go mod download
RUN go install gotest.tools/gotestsum@latest

COPY . .

RUN ./test.sh || (echo "Tests failed or coverage too low!" && exit 1)

ARG TARGETARCH
RUN CGO_ENABLED=0 GOARCH=$TARGETARCH go build -o /app/server ./main.go || (echo "Build failed!" && exit 1)

FROM alpine:latest AS final

COPY --from=build /app/server /app/server
COPY --from=build /docker/entrypoint.sh /docker/entrypoint.sh
COPY .env.development .env

EXPOSE 8000

ENTRYPOINT ["/docker/entrypoint.sh"]
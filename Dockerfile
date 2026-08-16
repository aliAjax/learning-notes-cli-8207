FROM golang:1.23-alpine AS builder

WORKDIR /src
COPY go.mod ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -trimpath -ldflags="-s -w" -o /out/learning-notes ./cmd/learning-notes

FROM alpine:3.20

RUN apk add --no-cache ca-certificates
ENV LEARNING_NOTES_DATA_DIR=/notes
VOLUME ["/notes"]
WORKDIR /notes

COPY --from=builder /out/learning-notes /usr/local/bin/learning-notes

ENTRYPOINT ["learning-notes"]

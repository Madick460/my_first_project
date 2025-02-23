FROM golang:1.23-bookworm AS builder

WORKDIR /app

RUN apt-get update && apt-get install -y \
    make \
    gcc \
    g++ \
    tesseract-ocr \
    libleptonica-dev \
    libtesseract-dev \
    wget \
    && rm -rf /var/lib/apt/lists/*

RUN mkdir -p /usr/share/tessdata && \
    wget -qO /usr/share/tessdata/eng.traineddata https://github.com/tesseract-ocr/tessdata_best/raw/main/eng.traineddata && \
    wget -qO /usr/share/tessdata/rus.traineddata https://github.com/tesseract-ocr/tessdata_best/raw/main/rus.traineddata

ENV TESSDATA_PREFIX=/usr/share/tessdata

COPY go.mod .
COPY go.sum .
RUN go mod download

COPY . .
RUN make build

FROM debian:bookworm-slim

RUN apt-get update && apt-get install -y \
    tesseract-ocr \
    libtesseract5 \
    && rm -rf /var/lib/apt/lists/*

COPY --from=builder /app/configs /app/configs
COPY --from=builder /app/.env /app/.env
COPY --from=builder /app/bin/main /app/main

WORKDIR /app

CMD ["/app/main"]

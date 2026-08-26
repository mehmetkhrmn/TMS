#derleyicimiz
FROM golang:1.26 AS builder
#containder icinde app dosyasi  olusturuyoruz
WORKDIR /app
#app a depen. leri kopyaliyoruz
COPY go.mod go.sum ./
#indiriyoruz
RUN go mod download
#source kodu kopyaliyoruz
COPY . .
#buildliyoruz
RUN go build -o ticket-service .

#artik calismasi icin go diline gerek yok zaten buildlendi baska kucuk bir imageye geciyoruz
FROM debian:bookworm-slim

RUN apt-get update \
    && apt-get install -y --no-install-recommends curl \
    && rm -rf /var/lib/apt/lists/* \

WORKDIR /app
#onceki buildi buraya kopyaliyoruz
COPY --from=builder /app/ticket-service .

EXPOSE 8080
#container calistirildiginda calisacak komut
CMD ["./ticket-service"]
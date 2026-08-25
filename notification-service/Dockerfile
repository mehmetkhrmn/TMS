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
RUN go build -o notification-service .

#artik calismasi icin go diline gerek yok zaten buildlendi baska kucuk bir imageye geciyoruz
FROM debian:bookworm-slim

WORKDIR /app
#onceki buildi buraya kopyaliyoruz
COPY --from=builder /app/notification-service .

EXPOSE 8081
#container calistirildiginda calisacak komut
CMD ["./notification-service"]
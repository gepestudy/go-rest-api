# Go Rest API
belajar rest api menggunakan standard library golang

### Membuat Certificate dan Key PEM

Untuk membuat certificate dan key PEM, jalankan perintah berikut di terminal:
```sh
openssl req -x509 -nodes -days 3650 -newkey rsa:2048 -keyout key.pem -out cert.pem -config openssl.cnf -sha256
```
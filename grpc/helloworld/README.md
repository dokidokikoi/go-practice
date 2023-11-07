## 生成 ssl 证书
```sh
$ openssl req -newkey rsa:2048 -nodes -keyout server.key -x509 -days 365 -out server.crt -subj "/CN=localhost" -addext "subjectAltName = DNS:localhost,DNS:foo/bar/my-service"
```
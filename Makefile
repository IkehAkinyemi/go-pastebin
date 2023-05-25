tls:
	mkdir tls; cd tls; go run /usr/local/go/src/crypto/tls/generate_cert.go --rsa-bits=2048 --host=localhost; cd ..

server:
	go run ./cmd/web

.PHONY: tls server
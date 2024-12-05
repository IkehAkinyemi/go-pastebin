tls:
	mkdir tls; cd tls; go run /usr/local/go/src/crypto/tls/generate_cert.go --rsa-bits=2048 --host=localhost; cd ..

server:
	go run ./cmd/web

MYSQL_ROOT_PASSWORD := root
MYSQL_DATABASE := snippetbox
MYSQL_USER := dev
MYSQL_PASSWORD := Akinyemi1234@
MYSQL_PORT := 3306
CONTAINER_NAME := snippetbox-mysql

mysql-up:
	@echo "Starting MySQL container..."
	@docker run --name $(CONTAINER_NAME) \
		-e MYSQL_ROOT_PASSWORD=$(MYSQL_ROOT_PASSWORD) \
		-e MYSQL_DATABASE=$(MYSQL_DATABASE) \
		-e MYSQL_USER=$(MYSQL_USER) \
		-e MYSQL_PASSWORD=$(MYSQL_PASSWORD) \
		-p $(MYSQL_PORT):3306 \
		-d mysql:8.0

mysql-init:
	@echo "Waiting for MySQL to be ready..."
	@sleep 20
	@docker exec $(CONTAINER_NAME) mysql -uroot -p$(MYSQL_ROOT_PASSWORD) \
		-e "GRANT ALL PRIVILEGES ON $(MYSQL_DATABASE).* TO '$(MYSQL_USER)'@'%';" \
		-e "FLUSH PRIVILEGES;"
	@echo "MySQL initialized successfully!"

mysql-down:
	@echo "Stopping MySQL container..."
	@docker stop $(CONTAINER_NAME) || true
	@docker rm $(CONTAINER_NAME) || true

setup-mysql: mysql-down mysql-up mysql-init

dev: setup-mysql tls server

.PHONY: dev tls server mysql-up mysql-down mysql-init

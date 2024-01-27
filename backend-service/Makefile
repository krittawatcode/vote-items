.PHONY: create-keypair migrate-create migrate-up migrate-down migrate-force

PWD = $(shell pwd)
PATH = $(PWD)
CPATH = $(PWD)/cert
MPATH = $(PWD)/migration

PORT = 5432

# Default number of migrations to execute up or down
N = 1

create-keypair:
	@echo "Creating an rsa 256 key pair"
	openssl genpkey -algorithm RSA -out $(CPATH)/rsa_private_$(ENV).pem -pkeyopt rsa_keygen_bits:2048
	openssl rsa -in $(CPATH)/rsa_private_$(ENV).pem -pubout -out $(CPATH)/rsa_public_$(ENV).pem

migrate-create:
	@echo "---Creating migration files---"
	migrate create -ext sql -dir $(MPATH) -seq -digits 5 $(NAME)

migrate-up:
	migrate -source file://$(MPATH) -database 'postgres://postgres:password@localhost:$(PORT)/postgres?sslmode=disable' up $(N)

migrate-down:
	migrate -source file://$(MPATH) -database 'postgres://postgres:password@localhost:$(PORT)/postgres?sslmode=disable' down $(N)

migrate-force:
	migrate -source file://$(MPATH) -database 'postgres://postgres:password@localhost:$(PORT)/postgres?sslmode=disable' force $(VERSION)
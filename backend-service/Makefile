.PHONY: create-keypair swagger

PWD = $(shell pwd)
PATH = $(PWD)
CPATH = $(PWD)/cert

create-keypair:
	@echo "Creating an rsa 256 key pair"
	openssl genpkey -algorithm RSA -out $(CPATH)/rsa_private_$(ENV).pem -pkeyopt rsa_keygen_bits:2048
	openssl rsa -in $(CPATH)/rsa_private_$(ENV).pem -pubout -out $(CPATH)/rsa_public_$(ENV).pem

# still have bug with swag init, need to copy command to terminal
swagger:
	@echo "Generating swagger"
    PATH=$(go env GOPATH)/bin:$$PATH swag init --parseDependency --parseInternal
.PHONY: create-keypair

PWD = $(shell pwd)
PATH = $(PWD)

create-keypair:
	@echo "Creating an rsa 256 key pair"
	openssl genpkey -algorithm RSA -out $(PATH)/rsa_private_$(ENV).pem -pkeyopt rsa_keygen_bits:2048
	openssl rsa -in $(PATH)/rsa_private_$(ENV).pem -pubout -out $(PATH)/rsa_public_$(ENV).pem
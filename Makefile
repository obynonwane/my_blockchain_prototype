SHELL := /bin/bash

private_key:
	go run -mod=vendor ./cmd/utilities/scripts/private_key.go
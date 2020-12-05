.PHONY: gen mod

gen:
	go mod vendor
	bash ./gen.sh crd
	rm -rf vendor

mod:
	go mod download
	go mod tidy
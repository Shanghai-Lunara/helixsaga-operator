.PHONY: gen mod

gen:
	bash ./gen.sh crd

mod:
	go mod download
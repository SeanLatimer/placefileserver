BUILD := build
EXECUTABLE := placefileserver

ifeq ($(OS), Windows_NT)
	EXECUTABLE := $(EXECUTABLE).exe
endif

.PHONY: generate
generate: ; go generate


.PHONY: build
build: generate
	go build -o ./$(BUILD)/$(EXECUTABLE)

.PHONY: clean
clean: ; rm -rf ./$(BUILD)

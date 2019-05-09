BUILD := build
EXECUTABLE := placefileserver

ifeq ($(OS), Windows_NT)
	EXECUTABLE := $(EXECUTABLE).exe
endif

.PHONY: build
build: ; go build -o ./$(BUILD)/$(EXECUTABLE)

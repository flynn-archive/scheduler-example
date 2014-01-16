
build:
	mkdir -p build
	go build -o build/scheduler-example

clean:
	rm -rf build

.PHONY: build

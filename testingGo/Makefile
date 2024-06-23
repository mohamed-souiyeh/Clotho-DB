.DEFAULT_GOAL := run 
MAIN := main


fmt: 
	@echo "formating..."
	@go fmt ./... 

vet: fmt 
	@echo "veting..."
	@go vet ./... 

build: vet 
	@echo "building..."
	@go build -o ${MAIN}

run: vet
	@echo "running..."
	@go run main.go

clean:
	@echo "cleaning..."
	@rm -rf ${MAIN}

.PHONY:fmt vet build run


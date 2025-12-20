# process_code: use 'go fmt' and 'go vet' process go code
# build: build linux exec of back-end services
# start: local start back-end services，with '-race' flag
# protoc: generate go and ts code from proto file
.PHONY: process_code build start

process_code: process_code.ps1
	powershell -executionpolicy bypass -File "./scripts/process_code.ps1"

build: build_linux.ps1
	powershell -executionpolicy bypass -File "./scripts/build_linux.ps1"

start: start_local_windows.ps1
	powershell -executionpolicy bypass -File "./scripts/start_local_windows.ps1"

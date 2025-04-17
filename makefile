run-all-tests:
	go test ./... -coverprofile cover.out
	go tool cover -html='cover.out'

install-tools-cover:
	go get golang.org/x/tools/cmd/cover


default:
	go vet ./...
	staticcheck ./...
	go test
	go build ./cmd/content2html

clean:
	-rm -f content2html
	go clean

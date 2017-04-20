COVERFILE = coverprofile.out

test:
	go test ./...

cover:
	go test -coverprofile $(COVERFILE)
	go tool cover -html $(COVERFILE)

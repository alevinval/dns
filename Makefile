COVERFILE = coverprofile.out

test:
	go test -cover -v

cover:
	go test -coverprofile $(COVERFILE)
	go tool cover -html $(COVERFILE)

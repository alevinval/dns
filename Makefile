COVERFILE = coverprofile.out

test:
	go test -cover

cover:
	go test -coverprofile $(COVERFILE)
	go tool cover -html $(COVERFILE)

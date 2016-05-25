GO_FILES=*/*.go


############################
#          BUILD           #
############################

install : go-build

go-build : $(GO_FILES)
		env GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build ./...


############################
#          SETUP           #
############################

setup: submodule go-deps

submodule:
		git submodule init && \
		git submodule update

go-deps :
		go get -t ./...


############################
#           TEST           #
############################

test :
		# in test
		go test -v -p 1 $(shell go list ./... | grep -v /vendor/)


############################
#          CLEAN           #
############################

clean :

.PHONY : install go-build setup go-deps test clean

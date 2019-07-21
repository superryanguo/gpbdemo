all: compile

.PHONY: compile
PROTOC_GEN_GO := $(GOPATH)/bin/protoc-gen-go

# If $GOPATH/bin/protoc-gen-go does not exist, we'll run this command to install
# it.
$(PROTOC_GEN_GO):
	go get -u github.com/golang/protobuf/protoc-gen-go

myobject.pb.go: proto/myobject.proto | $(PROTOC_GEN_GO)
	protoc --go_out=. proto/myobject.proto

compile: myobject.pb.go

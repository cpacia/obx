sample-config:
	cd repo && go-bindata -pkg=repo sample-obxd.conf

protos:
	cd models && PATH=$(PATH):$(GOPATH)/bin protoc --go_out=./ *.proto
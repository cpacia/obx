sample-config:
	cd repo && go-bindata -pkg=repo sample-obx.conf

protos:
	cd models && PATH=$(PATH):$(GOPATH)/bin protoc --go_out=./ *.proto
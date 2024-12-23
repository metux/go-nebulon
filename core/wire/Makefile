include ../../make.conf

compile: gen-proto doc

gen-proto:
	@protoc -I=../spec --go_out=. --go-grpc_out=. ../spec/core/*.proto ../spec/rpc/*.proto \
		--go_opt=M../spec/core/blockref.proto=github.com/metux/go-nebulon/wire

doc:
	@mkdir -p ../../doc
	@protoc -I=../spec --doc_out=../../doc --doc_opt=html,wire-format.html ../spec/core/*.proto
	@protoc -I=../spec --doc_out=../../doc --doc_opt=html,rpc-protocol.html ../spec/rpc/*.proto
#	@protoc -I=../spec --doc_out=../../doc --doc_opt=html,index.html,source_relative ../spec/core/*.proto ../spec/rpc/*.proto

clean:
	@rm -f *.pb.go ../doc/wire-format.html ../doc/rpc-protocol.html

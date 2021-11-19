PATH=$PATH:$GOPATH/bin
protodir=../../pb

protoc $protodir/shortener.proto --go_out=. --go-grpc_out=. --proto_path=$protodir

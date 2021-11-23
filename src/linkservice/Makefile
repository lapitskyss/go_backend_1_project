proto:
	protoc ../../pb/shortener.proto --go_out=. --go-grpc_out=. --proto_path=../../pb

mock:
	mockgen -destination=./internal/store/mock/mock_link.go -package=mock github.com/lapitskyss/go_backend_1_project/src/linkservice/internal/services/linksrv LinkStore
	mockgen -destination=./internal/store/mock/mock_redirect_log.go -package=mock github.com/lapitskyss/go_backend_1_project/src/linkservice/internal/services/redirectsrv RedirectStore

.PHONY: proto mock

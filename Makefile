run :
	go run ./cmd/api


wire: ## Generate wire_gen.go
	cd pkg/di && wire

swag: 
	swag init -g cmd/api/main.go -o ./cmd/docs 

test:
	go test ./...

mock: ##make mock files using mockgen
	mockgen -source pkg/repository/interfaces/user.go -destination pkg/repository/mock/user_mock.go -package mock
	mockgen -source pkg/usecase/interface/user.go -destination pkg/usecase/mock/user_mock.go -package mock


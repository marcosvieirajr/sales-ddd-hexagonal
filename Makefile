.PHONY: test test-verbose help

help:
	@echo "Targets disponíveis:"
	@echo "  make test         - Executar todos os testes"
	@echo "  make test-verbose - Executar todos os testes com saída verbosa"

test:
	go test ./catalog/... ./customer/... ./order/... ./kernel/...

test-verbose:
	go test -v ./catalog/... ./customer/... ./order/... ./kernel/...

test-coverage:
	go test -coverprofile=coverage.out ./catalog/... ./customer/... ./order/... ./kernel/...
	go tool cover -html=coverage.out -o coverage.html


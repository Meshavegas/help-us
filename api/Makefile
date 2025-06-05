.PHONY: run build test clean install dev swagger

# Variables
BINARY_NAME=api
MAIN_PATH=./main.go

# Installer les dépendances
install:
	go mod tidy
	go mod download

# Lancer l'application en mode développement
dev:
	go run $(MAIN_PATH)

# Lancer l'application
run:
	go run $(MAIN_PATH)

# Compiler l'application
build:
	go build -o $(BINARY_NAME) $(MAIN_PATH)

# Compiler pour différentes plateformes
build-all:
	GOOS=linux GOARCH=amd64 go build -o $(BINARY_NAME)-linux-amd64 $(MAIN_PATH)
	GOOS=windows GOARCH=amd64 go build -o $(BINARY_NAME)-windows-amd64.exe $(MAIN_PATH)
	GOOS=darwin GOARCH=amd64 go build -o $(BINARY_NAME)-darwin-amd64 $(MAIN_PATH)

# Lancer les tests
test:
	go test -v ./...

# Lancer les tests avec couverture
test-coverage:
	go test -v -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html

# Générer la documentation Swagger
swagger:
	go run github.com/swaggo/swag/cmd/swag@latest init
	@echo "Documentation Swagger générée dans le dossier docs/"
	@echo "Démarrez l'API et visitez http://localhost:8080/swagger/index.html"

# Nettoyer les fichiers générés
clean:
	go clean
	rm -f $(BINARY_NAME)
	rm -f $(BINARY_NAME)-*
	rm -f coverage.out coverage.html
	rm -f *.db
	rm -rf docs/

# Formater le code
fmt:
	go fmt ./...

# Vérifier le code
vet:
	go vet ./...

# Linter (nécessite golangci-lint)
lint:
	golangci-lint run

# Lancer l'application en mode production
prod:
	GIN_MODE=release go run $(MAIN_PATH)

# Aide
help:
	@echo "Commandes disponibles:"
	@echo "  install      - Installer les dépendances"
	@echo "  dev          - Lancer en mode développement"
	@echo "  run          - Lancer l'application"
	@echo "  build        - Compiler l'application"
	@echo "  build-all    - Compiler pour toutes les plateformes"
	@echo "  test         - Lancer les tests"
	@echo "  test-coverage- Tests avec couverture"
	@echo "  swagger      - Générer la documentation Swagger"
	@echo "  clean        - Nettoyer les fichiers générés"
	@echo "  fmt          - Formater le code"
	@echo "  vet          - Vérifier le code"
	@echo "  lint         - Linter le code"
	@echo "  prod         - Lancer en mode production"
	@echo "  help         - Afficher cette aide" 
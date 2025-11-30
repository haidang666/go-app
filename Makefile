APP_NAME = github.com/haidang666/go-app
CMD_PATH = ./cmd

install:
	go mod tidy
	go mod download

run:
	go run $(CMD_PATH)/server/main.go

# build:
# 	go build -o bin/$(APP_NAME) $(CMD_PATH)

format: 
	go fmt ./...

lint: 
	go vet ./...


# ent-create:
# 	go run -mod=mod entgo.io/ent/cmd/ent new ${name}

# ent-gen:
# 	go generate ./ent

# wire-gen:
# 	wire ./internal/app

# install-tools:
# 	go install github.com/swaggo/swag/cmd/swag@latest
# 	go install entgo.io/ent/cmd/ent
# 	go install github.com/google/wire/cmd/wire@latest
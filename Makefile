TARGET_DIR:=$(abspath $(lastword ./))
APP_NAME:=mypass

build:
	@sh -c "go mod tidy"
	@sh -c "go build -o $(TARGET_DIR)/$(APP_NAME) -tags=jsoniter -ldflags "-s -w"  $(TARGET_DIR)/main.go"

install:
	@sh -c "go mod tidy"
	@sh -c "go build -o $(TARGET_DIR)/$(APP_NAME) $(TARGET_DIR)/main.go"
	@sh -c "cp $(TARGET_DIR)/$(APP_NAME) ~/bin/$(APP_NAME)"

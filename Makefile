CONFIG_FILE = config.json
TARGET_DIR = /opt/coraserver
all:
	cp $(CONFIG_FILE) $(TARGET_DIR)
	go build -o $(TARGET_DIR)

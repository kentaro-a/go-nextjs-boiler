#!/bin/sh
SCRIPT_DIR=$(cd $(dirname $0); pwd)
migrate -path $SCRIPT_DIR/sqls -database "mysql://app:12345678abc@tcp(db:3306)/app" drop -f
migrate -path $SCRIPT_DIR/sqls -database "mysql://app:12345678abc@tcp(db:3306)/app" up

migrate -path $SCRIPT_DIR/sqls -database "mysql://root:12345678abc@tcp(db:3306)/test_app" drop -f
migrate -path $SCRIPT_DIR/sqls -database "mysql://root:12345678abc@tcp(db:3306)/test_app" up

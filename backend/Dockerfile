FROM golang:1.17.1
WORKDIR /backend

RUN go get -u github.com/cosmtrek/air
RUN go get -u github.com/stretchr/testify
RUN go install -tags 'mysql' github.com/golang-migrate/migrate/v4/cmd/migrate@latest

# RUN migrate create -dir migrations -ext sql -seq create_users
# RUN migrate -path migrations/sqls -database "mysql://app:12345678abc@tcp(db:3306)/app" up

# Start server
# $ cd /backend/src 
# $ air

# # seeds info
# # login salt
# salt: mysalt
# crypt: SHA256

# # user1
# user1@test.com
# 12345678abc
# # user2
# user2@test.com
# 12345678abcd

# Test
# go test -count=1 -v ./handler/ -run TestSignIn

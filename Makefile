run:
	go run main.go

schema:
	@read -p "Enter Schema Name: " name; \
	   go run entgo.io/ent/cmd/ent new $$name

generate:
	go generate ./ent

linux:
	env GOOS=linux GOARCH=amd64 go build -o ./bin/linux/biz

mac:
	env GOOS=mac GOARCH=amd64 go build -o ./bin/mac/biz

windows:
	env GOOS=windows GOARCH=amd64 go build -o ./bin/windows/biz
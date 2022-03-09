go env -w GO111MODULE=auto 
go build -o NewMoney.exe main.go
go build -ldflags="-H windowsgui" -o NewMoney_NoTerminal.exe main.go
go env -w GO111MODULE=auto 
go build -ldflags="-H windowsgui" -o NewMoney_NoTerminal.exe main.go
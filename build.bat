go generate -v
go build -o llama.cpp-ui-cpu.exe -ldflags "-H windowsgui" -v
go build -o llama.cpp-ui-cpu-cmd.exe -v
pause
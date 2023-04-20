go generate -v
go build -o llama.cpp-ui-cuda.exe -tags="cuda" -ldflags "-H windowsgui" -v
go build -o llama.cpp-ui-cuda-cmd.exe -tags="cuda" -v
pause
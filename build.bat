go generate -v
go build -o llama.cpp-ui-cpu.exe -ldflags "-H windowsgui" -v
go build -o llama.cpp-ui-cpu-cmd.exe -v
"C:\Program Files\7-Zip\7z.exe" a llama.cpp-ui-cpu.zip llama.cpp-ui-cpu.exe llama.cpp-ui-cpu-cmd.exe
pause
go generate -v
go build -o llama.cpp-ui-cuda.exe -tags="cuda" -ldflags "-H windowsgui" -v
go build -o llama.cpp-ui-cuda-cmd.exe -tags="cuda" -v
"C:\Program Files\7-Zip\7z.exe" a llama.cpp-ui-cuda.zip llama.cpp-ui-cuda.exe llama.cpp-ui-cuda-cmd.exe
pause


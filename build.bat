SET CGO_ENABLED=1
SET GOOS=linux
SET GOARCH=amd64
set CC=x86_64-linux-musl-gcc
set CXX=x86_64-linux-musl-g++
go build -o hidethread -trimpath -ldflags "-w -s -linkmode \"external\" -extldflags \"-static -O3\""  --tags "fts5"  


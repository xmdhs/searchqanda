SET CGO_ENABLED=0
SET GOOS=js
SET GOARCH=wasm
go build -o s.wasm  -trimpath -ldflags "-w -s"
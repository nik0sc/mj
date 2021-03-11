build.handcheck.wasm.gc:
	GOOS=js GOARCH=wasm go build -ldflags="-s -w" -o assets/handcheck.wasm cmd/handcheck_wasm/handcheck.go
	ls -lh assets/handcheck.wasm

build.handcheck.wasm.tinygo:
	tinygo build -o=assets_tinygo/handcheck.wasm -target=wasm -no-debug cmd/handcheck_wasm_tinygo/handcheck.go
	ls -lh assets_tinygo/handcheck.wasm

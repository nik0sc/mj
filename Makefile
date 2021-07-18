COMMITHASH := $(shell git rev-parse HEAD)
TINYGOROOT := $(shell readlink -f `which tinygo` | sed "s/\/bin\/tinygo//")

build.handcheck.wasm.gc:
	GOOS=js GOARCH=wasm go build \
		-ldflags="-s -w" -o assets/handcheck.wasm \
		cmd/handcheck_wasm/handcheck.go
	ls -lh assets/handcheck.wasm

build.handcheck.wasm.tinygo:
	tinygo build \
		-ldflags="-X main.commithash=$(COMMITHASH)" \
		-o=assets_tinygo/handcheck.wasm -target=wasm -no-debug \
		cmd/handcheck_wasm_tinygo/handcheck.go
# Try the container path first
	cp /tinygo/targets/wasm_exec.js assets_tinygo/ || cp $(TINYGOROOT)/targets/wasm_exec.js assets_tinygo/ 
	ls -lh assets_tinygo

clean:
	rm assets/handcheck.wasm
	rm assets_tinygo/handcheck.wasm
	rm assets_tinygo/wasm_exec.js
This example based these descriptions:
- https://tderflinger.com/en/how-to-integrate-go-library-js-webpage-webassembly
- https://eli.thegreenplace.net/2021/a-comprehensive-guide-to-go-generate

```bash
cp "$(go env GOROOT)/misc/wasm/wasm_exec.js" .
```

To generate the wasm file, run the following command:
```bash
go generate
```

Tart th server:
```bash
go run server.go
```




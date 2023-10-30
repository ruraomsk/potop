//go:build wasm

package main

import "embed"

//go:embed resources config
var resources embed.FS

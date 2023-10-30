//go:build !wasm

package main

import "embed"

//go:embed resources media config
var resources embed.FS

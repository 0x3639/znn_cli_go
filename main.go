package main

import (
	"github.com/0x3639/znn_cli_go/cmd"
	_ "github.com/0x3639/znn_cli_go/cmd/wallet" // Import for init() registration
)

func main() {
	cmd.Execute()
}

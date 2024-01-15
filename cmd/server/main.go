package main

import (
	"bbu/internal/blockchain"
	"bbu/internal/command_line"
	"log"
)

func main() {
	bc, err := blockchain.NewBlockChain()
	if err != nil {
		log.Fatal(err)
	}
	defer bc.Close()

	cli := command_line.NewCommandLine(bc)
	cli.Run()
}

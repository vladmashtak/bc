package command_line

import (
	"bbu/internal/blockchain"
	"flag"
	"fmt"
	"log"
	"os"
	"runtime"
)

type CommandLine struct {
	blockChain *blockchain.BlockChain
}

func NewCommandLine(bc *blockchain.BlockChain) *CommandLine {
	return &CommandLine{
		blockChain: bc,
	}
}

func (cli *CommandLine) PrintUsage() {
	fmt.Println("Usage: ")
	fmt.Println(" add -block BLOCK_DATA - add block to the chain")
	fmt.Println(" print - prints the blocks in the chain")
}

func (cli *CommandLine) AddBlock(data string) {
	cli.blockChain.AddBlock(data)
	fmt.Println("added block !")
}

func (cli *CommandLine) PrintChain() {
	iter := cli.blockChain.Iterator()

	for {
		block := iter.Next()
		if block == nil {
			break
		}

		fmt.Println(block)

		if len(block.PrevHash()) == 0 {
			break
		}
	}
}

func (cli *CommandLine) Run() {
	cli.ValidateArgs()

	addBlockCmd := flag.NewFlagSet("add", flag.ExitOnError)
	printBlockCmd := flag.NewFlagSet("print", flag.ExitOnError)

	addBlockData := addBlockCmd.String("block", "", "Block data")

	switch os.Args[1] {
	case "add":
		err := addBlockCmd.Parse(os.Args[2:])
		if err != nil {
			log.Println(err)
		}
	case "print":
		err := printBlockCmd.Parse(os.Args[2:])
		if err != nil {
			log.Println(err)
		}
	default:
		cli.PrintUsage()
		runtime.Goexit()
	}

	if addBlockCmd.Parsed() {
		if *addBlockData == "" {
			addBlockCmd.Usage()
			runtime.Goexit()
		}

		cli.AddBlock(*addBlockData)
	}

	if printBlockCmd.Parsed() {
		cli.PrintChain()
	}
}

func (cli *CommandLine) ValidateArgs() {
	if len(os.Args) < 2 {
		cli.PrintUsage()
		runtime.Goexit()
	}
}

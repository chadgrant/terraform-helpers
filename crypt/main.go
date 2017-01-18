package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/chadgrant/terraform-helpers/crypt/encryption"
)

var key = ""

func main() {
	var operation string
	var key = os.Getenv("TERRAFORM_DECRYPT")

	flags := flag.NewFlagSet("crypt", flag.ExitOnError)
	flags.Usage = printUsage
	flags.StringVar(&operation, "operation", "encrypt", "encrypt or decrypt")
	flags.StringVar(&key, "key", os.Getenv("TERRAFORM_DECRYPT"), "encryption key")

	if err := flags.Parse(os.Args[1:]); err != nil {
		flags.Usage()
		os.Exit(1)
		return
	}

	if len(key) <= 0 {
		fmt.Println("decryption key required")
		return
	}

	path, _ := os.Getwd()
	path += string(os.PathSeparator)
	if len(flags.Args()) > 0 {
		path = flags.Args()[0]
	}

	fmt.Printf("Searching %s for files to %s\n", path, operation)

	if operation == "encrypt" {
		if err := encryption.EncryptFiles([]byte(key), path); err != nil {
			fmt.Println("Error searching files")
			os.Exit(1)
			return
		}
	}

	if operation == "decrypt" {
		if err := encryption.DecryptFiles([]byte(key), path); err != nil {
			fmt.Println("Error searching files")
			os.Exit(1)
			return
		}
	}
}

const helpText = `Usage: crypt [options] [directory]
  crypt searches recursively for tfvar files stored under [directory]
  to decrypt files stored in source

Options:
  -operation        encrypt or decrypt (encrypt is default)
  -key              the encryption/decryption key

Output :
	files will be created, decrypted or encrypted depending on operation mode
`

func printUsage() {
	fmt.Fprintf(os.Stderr, helpText)
}

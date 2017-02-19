// Copyright (c) 2017 Dongsu Park <dpark@posteo.net>
//
// Permission to use, copy, modify, and distribute this software for any
// purpose with or without fee is hereby granted, provided that the above
// copyright notice and this permission notice appear in all copies.
//
// THE SOFTWARE IS PROVIDED "AS IS" AND THE AUTHOR DISCLAIMS ALL WARRANTIES
// WITH REGARD TO THIS SOFTWARE INCLUDING ALL IMPLIED WARRANTIES OF
// MERCHANTABILITY AND FITNESS. IN NO EVENT SHALL THE AUTHOR BE LIABLE FOR
// ANY SPECIAL, DIRECT, INDIRECT, OR CONSEQUENTIAL DAMAGES OR ANY DAMAGES
// WHATSOEVER RESULTING FROM LOSS OF USE, DATA OR PROFITS, WHETHER IN AN
// ACTION OF CONTRACT, NEGLIGENCE OR OTHER TORTIOUS ACTION, ARISING OUT OF
// OR IN CONNECTION WITH THE USE OR PERFORMANCE OF THIS SOFTWARE.

package cli

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/dongsupark/go-subtitle/parser"
)

var (
	readCmd = &cobra.Command{
		Use:   "read [-f|--file] FILENAME",
		Short: "A command line client for reading a subtitle file",
		Long: `A CLI for reading a subtitle file

To get help about a resource or command, please run "go-subtitle help"`,
		Run: runReadCmd,
	}

	readFileName string
)

func init() {
	goSubtitleCmd.AddCommand(readCmd)

	readCmd.Flags().StringVarP(&readFileName, "file", "f", "", "Subtitle file")
}

func runReadCmd(cmd *cobra.Command, args []string) {
	if len(readFileName) == 0 {
		cmd.Help()
		return
	}

	if _, err := os.Stat(readFileName); os.IsNotExist(err) {
		fmt.Printf("file not found: %s\n", readFileName)
		return
	}

	outSt, err := parser.ReadSubFromFile(readFileName)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Printf("Parsed text:\n%s\n", outSt.ToText())

	return
}

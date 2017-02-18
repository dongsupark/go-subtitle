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
	"github.com/dongsupark/go-subtitle/subtitle"
)

var (
	writeCmd = &cobra.Command{
		Use:   "write [-f|--file] FILENAME --rawtext DATA",
		Short: "A command line client for writing a subtitle file",
		Long: `A CLI for writing a subtitle file

To get help about a resource or command, please run "go-subtitle help"`,
		Run: runWriteCmd,
	}

	writeFileName string
)

func init() {
	goSubtitleCmd.AddCommand(writeCmd)

	writeCmd.Flags().StringVarP(&writeFileName, "file", "f", "", "Subtitle file")
}

func runWriteCmd(cmd *cobra.Command, args []string) {
	if len(writeFileName) == 0 {
		cmd.Help()
		return
	}

	if _, err := os.Stat(writeFileName); os.IsNotExist(err) {
		fmt.Printf("file not found: %s\n", writeFileName)
		return
	}

	subtitleFormat := parser.GetParserFormat(outputFileName)
	if subtitleFormat == "" {
		fmt.Println("unable to get subtitle format")
		return
	}

	writer := parser.GetParserWriter(subtitleFormat)
	if writer == nil {
		fmt.Println("unable to get parser writer")
		return
	}

	// TODO: make it work with data provided by user
	err := writer(writeFileName, subtitle.Subtitle{})
	//     err := writer(writeFileName, inSub)
	if err != nil {
		fmt.Printf("parse error reading %s: %v\n", readFileName, err)
		return
	}

	fmt.Printf("Wrote text to subtitle file %s\n", writeFileName)
	return
}

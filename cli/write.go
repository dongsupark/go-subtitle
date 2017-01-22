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

	_ "github.com/dongsupark/go-subtitle/parser"
)

var (
	writeCmd = &cobra.Command{
		Use:   "write --file FILENAME --rawtext DATA",
		Short: "A command line client for writing a subtitle file",
		Long: `A CLI for writing a subtitle file

To get help about a resource or command, please run "go-subtitle help"`,
		Run: runWriteCmd,
	}

	writeFileName string
)

func init() {
	goSubtitleCmd.AddCommand(writeCmd)

	writeCmd.Flags().StringVar(&writeFileName, "file", "", "Subtitle file")
	writeCmd.Flags().StringVar(&writeFileName, "f", "", "Shorthand for --file")
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

	subtitleFormat, err := cmd.PersistentFlags().GetString("format")
	if err != nil {
		return
	}

	if subtitleFormat == "subrip" || subtitleFormat == "srt" {
		// TODO: make it work with data provided by user
		//         var subripParser parser.SubripFormat
		//         err = subripParser.Write(writeFileName, data)
		if err != nil {
			fmt.Printf("parse error writing %s: %v\n", writeFileName, err)
			return
		}
	}

	fmt.Printf("Wrote text to subtitle file %s\n", writeFileName)
	return
}

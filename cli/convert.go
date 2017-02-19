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
	"path"

	"github.com/spf13/cobra"

	"github.com/dongsupark/go-subtitle/parser"
	"github.com/dongsupark/go-subtitle/subtitle"
)

var (
	convertCmd = &cobra.Command{
		Use:   "convert [-f|--inputfile] FILENAME [-o|--outputfile] FILENAME",
		Short: "A command line client for converting a subtitle file",
		Long: `A CLI for converting a subtitle file

To get help about a resource or command, please run "go-subtitle help"`,
		Run: runConvertCmd,
	}

	inputFileName  string
	outputFileName string
)

func init() {
	goSubtitleCmd.AddCommand(convertCmd)

	convertCmd.Flags().StringVarP(&inputFileName, "inputfile", "f", "", "Input subtitle file")
	convertCmd.Flags().StringVarP(&outputFileName, "outputfile", "o", "", "Input subtitle file")
}

func doReadFile(inputFileName string) (*subtitle.Subtitle, error) {
	inputFormat := parser.GetParserFormat(path.Base(inputFileName))
	if inputFormat == "" {
		return nil, fmt.Errorf("unable to get subtitle format\n")
	}

	reader := parser.GetParserReader(inputFormat)
	if reader == nil {
		return nil, fmt.Errorf("unable to get parser reader\n")
	}

	outSt, err := reader(inputFileName)
	if err != nil {
		return nil, fmt.Errorf("parse error converting %s: %v\n", inputFileName, err)
	}

	parsedText := ""
	for _, v := range outSt.Subtitles {
		parsedText += v.Text
	}

	fmt.Printf("Parsed text:\n%s\n", parsedText)
	return &outSt, nil
}

func doWriteFile(outputFileName string, outSt *subtitle.Subtitle) error {
	outputFormat := parser.GetParserFormat(path.Base(outputFileName))
	if outputFormat == "" {
		return fmt.Errorf("unable to get subtitle format\n")
	}

	writer := parser.GetParserWriter(outputFormat)
	if writer == nil {
		return fmt.Errorf("unable to get parser writer\n")
	}

	err := writer(outputFileName, *outSt)
	if err != nil {
		return fmt.Errorf("parse error reading %s: %v\n", outputFileName, err)
	}

	fmt.Printf("Wrote text to subtitle file %s\n", outputFileName)
	return nil
}

func runConvertCmd(cmd *cobra.Command, args []string) {
	if len(inputFileName) == 0 {
		cmd.Help()
		return
	}

	if len(outputFileName) == 0 {
		cmd.Help()
		return
	}

	if _, err := os.Stat(inputFileName); os.IsNotExist(err) {
		fmt.Printf("file not found: %s\n", inputFileName)
		return
	}

	outSt, err := doReadFile(inputFileName)
	if err != nil {
		fmt.Printf("Failed to read from %s: %v", inputFileName, err)
		return
	}

	if err := doWriteFile(outputFileName, outSt); err != nil {
		fmt.Printf("Failed to write to %s: %v", outputFileName, err)
		return
	}

	return
}

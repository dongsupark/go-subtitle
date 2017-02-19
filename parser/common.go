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

package parser

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/dongsupark/go-subtitle/subtitle"
)

var parserFormats = map[string]string{
	"subrip": "subrip",
	"srt":    "subrip",
	"sami":   "sami",
	"smi":    "sami",
}

var parserMap = map[string]interface{}{
	"subrip": SubripFormat{TypeName: "subrip"},
	"srt":    SubripFormat{TypeName: "subrip"},
	"sami":   SamiFormat{TypeName: "sami"},
	"smi":    SamiFormat{TypeName: "sami"},
}

type ReadFunc func(string) (subtitle.Subtitle, error)
type WriteFunc func(string, subtitle.Subtitle) error

func GetParserReader(formatName string) ReadFunc {
	switch formatName {
	case "subrip":
		fallthrough
	case "srt":
		sp := parserMap[formatName].(SubripFormat)
		return sp.Read
	case "sami":
		fallthrough
	case "smi":
		sp := parserMap[formatName].(SamiFormat)
		return sp.Read
	}

	return nil
}

func GetParserWriter(formatName string) WriteFunc {
	switch formatName {
	case "subrip":
		fallthrough
	case "srt":
		sp := parserMap[formatName].(SubripFormat)
		return sp.Write
	case "sami":
		fallthrough
	case "smi":
		sp := parserMap[formatName].(SamiFormat)
		return sp.Write
	}

	return nil
}

func GetParserFormat(filename string) string {
	for k, v := range parserFormats {
		if strings.ToLower(filepath.Ext(filename))[1:] == k {
			return v
		}
	}
	return ""
}

func ReadSubFromFile(readFileName string) (*subtitle.Subtitle, error) {
	subtitleFormat := GetParserFormat(readFileName)
	if subtitleFormat == "" {
		return nil, fmt.Errorf("unable to get subtitle format")
	}

	reader := GetParserReader(subtitleFormat)
	if reader == nil {
		return nil, fmt.Errorf("unable to get parser reader")
	}

	outSt, err := reader(readFileName)
	if err != nil {
		return nil, fmt.Errorf("parse error reading %s: %v\n", readFileName, err)
	}

	return &outSt, nil
}

func WriteSubToFile(writeFileName string, inSt subtitle.Subtitle) error {
	subtitleFormat := GetParserFormat(writeFileName)
	if subtitleFormat == "" {
		return fmt.Errorf("unable to get subtitle format")
	}

	writer := GetParserWriter(subtitleFormat)
	if writer == nil {
		return fmt.Errorf("unable to get parser writer")
	}

	err := writer(writeFileName, inSt)
	if err != nil {
		return fmt.Errorf("parse error writing %s: %v\n", writeFileName, err)
	}

	return nil
}

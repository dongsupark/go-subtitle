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
	"github.com/dongsupark/go-subtitle/subtitle"
)

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

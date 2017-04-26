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
	"strings"
	"testing"
)

const (
	inSamiData = `
<SYNC Start=20011><P Class=ENCC><i>this is a text</i>
<SYNC Start=22211><P Class=ENCC>&nbsp;
`
)

func TestSamiReadWrite(t *testing.T) {
	samiFormat := SamiFormat{TypeName: "sami"}

	inData := strings.TrimSpace(inSamiData)

	subtitle, err := samiFormat.Read(inData)
	if err != nil {
		t.Fatalf("Unable to read sami data\n")
	}

	outData, err := samiFormat.Write(subtitle)
	if err != nil {
		t.Fatalf("Unable to write sami data\n")
	}

	if inData != outData {
		t.Fatalf("does not match with original sami data\n")
	}
}

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
	"bufio"
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/dongsupark/go-subtitle/pkg"
	"github.com/dongsupark/go-subtitle/subtitle"
)

var (
	reSrtNum  = regexp.MustCompile("^\\d+$")
	reSrtTime = regexp.MustCompile("^(\\d+):(\\d+):(\\d+),(\\d+)\\s-->\\s(\\d+):(\\d+):(\\d+),(\\d+)")
)

type SubripFormat struct {
	TypeName string
}

func (sr *SubripFormat) Read(inputData string) (subtitle.Subtitle, error) {
	scanner := bufio.NewScanner(strings.NewReader(inputData))
	var st subtitle.Subtitle
	for scanner.Scan() {
		timeLine := scanner.Text()
		substrs := reSrtTime.FindStringSubmatch(timeLine)
		if substrs == nil || substrs[0] == "" {
			continue
		}

		var count = 0
		textLine := ""
		for scanner.Scan() {
			tmpLine := scanner.Text()
			if tmpLine == "" {
				break
			}
			if count > 0 {
				textLine += "\n"
			}
			textLine += tmpLine
			count++
		}

		var se subtitle.SubtitleEntry

		se.StartValue = pkg.ComposeTimeDuration(
			pkg.StringToInt(substrs[1]), pkg.StringToInt(substrs[2]),
			pkg.StringToInt(substrs[3]), pkg.StringToInt(substrs[4]))
		se.EndValue = pkg.ComposeTimeDuration(
			pkg.StringToInt(substrs[5]), pkg.StringToInt(substrs[6]),
			pkg.StringToInt(substrs[7]), pkg.StringToInt(substrs[8]))

		se.Text = textLine

		st.Subtitles = append(st.Subtitles, se)
	}

	if err := scanner.Err(); err != nil {
		return subtitle.Subtitle{}, err
	}

	return st, nil
}

func (sr *SubripFormat) Write(insub subtitle.Subtitle) (string, error) {
	dataStr := ""
	count := 1
	for _, v := range insub.Subtitles {
		dataStr += fmt.Sprintf("%d\n%s --> %s\n%s\n\n",
			count, timeToSubrip(v.StartValue), timeToSubrip(v.EndValue), v.Text)
		count++
	}

	return dataStr, nil
}

func timeToSubrip(inTime time.Duration) string {
	return fmt.Sprintf("%02d:%02d:%02d,%03d",
		int(inTime.Hours()), int(inTime.Minutes()), int(inTime.Seconds()), int(inTime.Nanoseconds()/1000%1000))
}

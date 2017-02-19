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
	"bytes"
	"fmt"
	"io"
	"os"
	"regexp"
	"strings"
	"time"

	"github.com/golang/glog"
	"golang.org/x/net/html"

	"github.com/dongsupark/go-subtitle/pkg"
	"github.com/dongsupark/go-subtitle/subtitle"
)

const (
	SamiStateInit      = 0
	SamiStateSyncStart = 1
	SamiStateText      = 2
	SamiStateSyncEnd   = 3
	SamiStateForceQuit = 99
)

type SamiStateType int

type SamiFormat struct {
	TypeName string
}

func (sr *SamiFormat) Read(fileName string) (subtitle.Subtitle, error) {
	glog.Infoln("reading sami file")

	fh, err := os.Open(fileName)
	if err != nil {
		return subtitle.Subtitle{}, err
	}
	defer fh.Close()

	var st subtitle.Subtitle
	se := new(subtitle.SubtitleEntry)
	ss := SamiStateType(SamiStateInit)

	doc, err := html.Parse(fh)
	if err != nil {
		glog.Infof("cannot parse data from %s\n", fileName)
		return subtitle.Subtitle{}, err
	}

	renl := regexp.MustCompile("\\n")

	var parseNode func(*html.Node, SamiStateType)
	parseNode = func(n *html.Node, samiState SamiStateType) {
		if n.Type == html.ElementNode && n.Data == "sync" {
			for _, a := range n.Attr {
				if a.Key == "start" {
					if samiState == SamiStateSyncStart || samiState == SamiStateText {
						se.EndValue = pkg.ComposeTimeDuration(0, 0, 0, pkg.StringToInt(a.Val))
						samiState = SamiStateSyncEnd

						st.Subtitles = append(st.Subtitles, *se)
						se = new(subtitle.SubtitleEntry)
					} else {
						se.StartValue = pkg.ComposeTimeDuration(0, 0, 0, pkg.StringToInt(a.Val))
						samiState = SamiStateSyncStart
					}
					break
				}
			}
		}
		if n.Type == html.TextNode {
			n.Data = stripComments(n.Data)

			inText := strings.TrimSpace(renl.ReplaceAllString(n.Data, " "))
			if len(inText) == 0 {
				if samiState == SamiStateSyncEnd {
					samiState = SamiStateInit
				}
			} else {
				se.Text = n.Data
				samiState = SamiStateText
			}
		}
		if n.Type == html.CommentNode {
			return
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			parseNode(c, samiState)
		}
	}
	parseNode(doc, ss)

	return st, nil
}

func (sr *SamiFormat) Write(fileName string, insub subtitle.Subtitle) error {
	glog.Infoln("writing sami file")

	doc := &html.Node{
		Type: html.DocumentNode,
	}
	for _, v := range insub.Subtitles {
		htmlText := strings.Replace(v.Text, "\n", "<br>", -1)

		sStartNode := &html.Node{
			Type: html.ElementNode,
			Data: fmt.Sprintf("<SYNC Start=%s>", timeToSami(v.StartValue)),
		}
		sPNode := &html.Node{
			Type: html.ElementNode,
			Data: fmt.Sprintf("<P Class=ENCC>\n"),
		}
		sPNode.AppendChild(&html.Node{
			Type: html.TextNode,
			Data: fmt.Sprintf("%s\n", htmlText),
		})
		sStartNode.AppendChild(sPNode)
		doc.AppendChild(sStartNode)

		sEndNode := &html.Node{
			Type: html.ElementNode,
			Data: fmt.Sprintf("<SYNC Start=%s>", timeToSami(v.EndValue)),
		}
		sPNode = &html.Node{
			Type: html.ElementNode,
			Data: fmt.Sprintf("<P Class=ENCC>\n"),
		}
		sPNode.AppendChild(&html.Node{
			Type: html.TextNode,
			Data: fmt.Sprintf("&nbsp;\n"),
		})
		sEndNode.AppendChild(sPNode)
		doc.AppendChild(sEndNode)
	}

	fh, err := os.Create(fileName)
	if err != nil {
		return err
	}
	defer fh.Close()

	if html.Render(fh, doc); err != nil {
		return err
	}

	return nil
}

// strip comments in every text node
func stripComments(inStr string) string {
	z := html.NewTokenizer(bytes.NewBufferString(inStr))
	for {
		tt := z.Next()
		if tt == html.ErrorToken {
			if err := z.Err(); err != nil && err != io.EOF {
				return inStr
			}
			break
		}
		if tt == html.CommentToken {
			return strings.Replace(inStr, string(z.Raw()), "", -1)
		}
	}

	return inStr
}

func timeToSami(inTime time.Duration) string {
	totalSec := inTime.Seconds()
	return fmt.Sprintf("%04d%03d", int(totalSec), int(inTime.Nanoseconds()/1000/1000%1000))
}

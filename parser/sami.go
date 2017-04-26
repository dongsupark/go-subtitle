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
	"bytes"
	"errors"
	"fmt"
	"io"
	"regexp"
	"strings"
	"time"

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

var textElemTags = []string{
	"i",
	"b",
}

func matchTextElemTag(input string) bool {
	for _, e := range textElemTags {
		if input == e {
			return true
		}
	}
	return false
}

var legitElemTags = []string{
	"&nbsp",
	"i>",
	"b>",
}

func hasLegitElemTag(input string) bool {
	for _, e := range legitElemTags {
		if strings.Contains(input, e) {
			return true
		}
	}
	return false
}

type SamiStateType int

type SamiFormat struct {
	TypeName string
}

func (sr *SamiFormat) Read(inputData string) (subtitle.Subtitle, error) {
	var st subtitle.Subtitle
	se := new(subtitle.SubtitleEntry)
	samiState := SamiStateType(SamiStateInit)

	inputData = strings.TrimSpace(inputData)

	renl := regexp.MustCompile("\\n")

	z := html.NewTokenizer(strings.NewReader(inputData))
	prevStartValue := time.Duration(0)
	for {
		tok := z.Next()
		switch tok {
		case html.ErrorToken:
			if z.Err() == io.EOF {
				break
			}
			return subtitle.Subtitle{}, fmt.Errorf("got error token")
		case html.StartTagToken:
			tn, hasAttr := z.TagName()
			tnStr := string(tn)

			if hasAttr && strings.ToLower(tnStr) == "sync" {
				key, value, _ := z.TagAttr()
				if strings.ToLower(string(key)) == "start" {
					if samiState == SamiStateSyncEnd {
						se.EndValue = prevStartValue
						samiState = SamiStateInit

						st.Subtitles = append(st.Subtitles, *se)
						se = new(subtitle.SubtitleEntry)
					}

					if samiState == SamiStateInit {
						se.StartValue = pkg.ComposeTimeDuration(0, 0, 0, pkg.StringToInt(string(value)))
						prevStartValue = se.StartValue
						samiState = SamiStateSyncStart
					} else if samiState == SamiStateSyncStart || samiState == SamiStateText {
						se.EndValue = pkg.ComposeTimeDuration(0, 0, 0, pkg.StringToInt(string(value)))
						samiState = SamiStateInit

						st.Subtitles = append(st.Subtitles, *se)
						se = new(subtitle.SubtitleEntry)
					}
					break
				}
			}

			// consider this node as a text node with an in-text tag
			if matchTextElemTag(tnStr) {
				se.Text += fmt.Sprintf("<%s>", tnStr)
			}
		case html.TextToken:
			toSyncEnd := false
			parsed := ""

			if strings.Contains(string(z.Raw()), "&nbsp") {
				parsed = string(z.Raw())
				toSyncEnd = true
			} else {
				parsed = string(z.Text())
			}

			if samiState == SamiStateSyncStart || samiState == SamiStateInit {
				textStr := stripComments(parsed)

				inText := strings.TrimSpace(renl.ReplaceAllString(textStr, " "))
				if len(inText) > 0 {
					se.Text += parsed

					if toSyncEnd {
						samiState = SamiStateSyncEnd
					} else {
						samiState = SamiStateText
					}
				}
			}
		case html.EndTagToken:
			tn, _ := z.TagName()
			tnStr := string(tn)

			if matchTextElemTag(tnStr) {
				se.Text += fmt.Sprintf("</%s>", tnStr)
			}
		case html.SelfClosingTagToken, html.CommentToken, html.DoctypeToken:
			// do nothing
		}

		if z.Err() == io.EOF {
			break
		}
	}

	return st, nil
}

func (sr *SamiFormat) Write(insub subtitle.Subtitle) (string, error) {
	doc := &html.Node{
		Type: html.DocumentNode,
	}
	for _, v := range insub.Subtitles {
		htmlText := strings.TrimSpace(html.UnescapeString(v.Text))

		sStartNode := &html.Node{
			Type: html.ElementNode,
			Data: fmt.Sprintf("SYNC Start=%s", timeToSami(v.StartValue)),
		}
		sPNode := &html.Node{
			Type: html.ElementNode,
			Data: "P Class=ENCC",
		}
		sPNode.AppendChild(&html.Node{
			Type: html.TextNode,
			Data: fmt.Sprintf("%s\n", htmlText),
		})
		sStartNode.AppendChild(sPNode)
		doc.AppendChild(sStartNode)

		sEndNode := &html.Node{
			Type: html.ElementNode,
			Data: fmt.Sprintf("SYNC Start=%s", timeToSami(v.EndValue)),
		}
		sPNode = &html.Node{
			Type: html.ElementNode,
			Data: "P Class=ENCC",
		}
		sPNode.AppendChild(&html.Node{
			Type: html.TextNode,
			Data: "&nbsp;\n",
		})
		sEndNode.AppendChild(sPNode)
		doc.AppendChild(sEndNode)
	}

	b := new(bytes.Buffer)
	if err := samiRender(b, doc); err != nil {
		return "", err
	}

	return strings.TrimSpace(b.String()), nil
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
	totalMsec := (int(totalSec) * 1000) + int(inTime.Nanoseconds()/1000/1000%1000)
	return fmt.Sprintf("%d", totalMsec)
}

type writer interface {
	io.Writer
	io.ByteWriter
	WriteString(string) (int, error)
}

func samiRender(w io.Writer, n *html.Node) error {
	if x, ok := w.(writer); ok {
		return doSamiRender(x, n)
	}
	buf := bufio.NewWriter(w)
	if err := doSamiRender(buf, n); err != nil {
		return err
	}
	return buf.Flush()
}

func doSamiRender(w writer, n *html.Node) error {
	// Render non-element nodes; these are the easy cases.
	switch n.Type {
	case html.ErrorNode:
		return errors.New("html: cannot render an ErrorNode node")
	case html.TextNode:
		// NOTE: we need to prevent several strings from being escaped,
		// for example, &nbsp", a special termination tag in sami format.
		if hasLegitElemTag(n.Data) {
			_, err := w.WriteString(n.Data)
			return err
		} else {
			return escape(w, n.Data)
		}
	case html.DocumentNode:
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			if err := doSamiRender(w, c); err != nil {
				return err
			}
		}
		return nil
	case html.ElementNode:
		// No-op.
	case html.CommentNode:
		if _, err := w.WriteString("<!--"); err != nil {
			return err
		}
		if _, err := w.WriteString(n.Data); err != nil {
			return err
		}
		if _, err := w.WriteString("-->"); err != nil {
			return err
		}
		return nil
	case html.DoctypeNode:
		if _, err := w.WriteString("<!DOCTYPE "); err != nil {
			return err
		}
		if _, err := w.WriteString(n.Data); err != nil {
			return err
		}
		if n.Attr != nil {
			var p, s string
			for _, a := range n.Attr {
				switch a.Key {
				case "public":
					p = a.Val
				case "system":
					s = a.Val
				}
			}
			if p != "" {
				if _, err := w.WriteString(" PUBLIC "); err != nil {
					return err
				}
				if err := writeQuoted(w, p); err != nil {
					return err
				}
				if s != "" {
					if err := w.WriteByte(' '); err != nil {
						return err
					}
					if err := writeQuoted(w, s); err != nil {
						return err
					}
				}
			} else if s != "" {
				if _, err := w.WriteString(" SYSTEM "); err != nil {
					return err
				}
				if err := writeQuoted(w, s); err != nil {
					return err
				}
			}
		}
		return w.WriteByte('>')
	default:
		return errors.New("html: unknown node type")
	}

	// Render the <xxx> opening tag.
	if err := w.WriteByte('<'); err != nil {
		return err
	}
	if _, err := w.WriteString(n.Data); err != nil {
		return err
	}
	for _, a := range n.Attr {
		if err := w.WriteByte(' '); err != nil {
			return err
		}
		if a.Namespace != "" {
			if _, err := w.WriteString(a.Namespace); err != nil {
				return err
			}
			if err := w.WriteByte(':'); err != nil {
				return err
			}
		}
		if _, err := w.WriteString(a.Key); err != nil {
			return err
		}
		if _, err := w.WriteString(`="`); err != nil {
			return err
		}
		if err := escape(w, a.Val); err != nil {
			return err
		}
		if err := w.WriteByte('"'); err != nil {
			return err
		}
	}
	if err := w.WriteByte('>'); err != nil {
		return err
	}

	// Add initial newline where there is danger of a newline beging ignored.
	if c := n.FirstChild; c != nil && c.Type == html.TextNode && strings.HasPrefix(c.Data, "\n") {
		switch n.Data {
		case "pre", "listing", "textarea":
			if err := w.WriteByte('\n'); err != nil {
				return err
			}
		}
	}

	// Render any child nodes.
	switch n.Data {
	case "iframe", "noembed", "noframes", "noscript", "plaintext", "script", "style", "xmp":
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			if c.Type == html.TextNode {
				if _, err := w.WriteString(c.Data); err != nil {
					return err
				}
			} else {
				if err := doSamiRender(w, c); err != nil {
					return err
				}
			}
		}
	default:
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			if err := doSamiRender(w, c); err != nil {
				return err
			}
		}
	}

	// NOTE: don't render a closing tag at all, for sami files.
	return nil
}

const escapedChars = "&'<>\"\r"

func escape(w writer, s string) error {
	i := strings.IndexAny(s, escapedChars)
	for i != -1 {
		if _, err := w.WriteString(s[:i]); err != nil {
			return err
		}
		var esc string
		switch s[i] {
		case '&':
			esc = "&amp;"
		case '\'':
			// "&#39;" is shorter than "&apos;" and apos was not in HTML until HTML5.
			esc = "&#39;"
		case '<':
			esc = "&lt;"
		case '>':
			esc = "&gt;"
		case '"':
			// "&#34;" is shorter than "&quot;".
			esc = "&#34;"
		case '\r':
			esc = "&#13;"
		default:
			panic("unrecognized escape character")
		}
		s = s[i+1:]
		if _, err := w.WriteString(esc); err != nil {
			return err
		}
		i = strings.IndexAny(s, escapedChars)
	}
	_, err := w.WriteString(s)
	return err
}

// writeQuoted writes s to w surrounded by quotes. Normally it will use double
// quotes, but if s contains a double quote, it will use single quotes.
// It is used for writing the identifiers in a doctype declaration.
// In valid HTML, they can't contain both types of quotes.
func writeQuoted(w writer, s string) error {
	var q byte = '"'
	if strings.Contains(s, `"`) {
		q = '\''
	}
	if err := w.WriteByte(q); err != nil {
		return err
	}
	if _, err := w.WriteString(s); err != nil {
		return err
	}
	if err := w.WriteByte(q); err != nil {
		return err
	}
	return nil
}

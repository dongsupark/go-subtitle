package main

import (
	"fmt"
	"log"
	"os"
	"strings"

	"golang.org/x/net/html"
)

func main() {
	s := "<BODY>\n<SYNC Start=0><P Class=KRCC>&nbsp;\n<SYNC Start=20011><P Class=KRCC>\n<i>this is italic text</i>\n<SYNC Start=22211><P Class=KRCC>&nbsp;\n</BODY>"

	doc, err := html.Parse(strings.NewReader(s))
	if err != nil {
		log.Fatal(err)
	}

	var parseNode func(*html.Node)
	parseNode = func(n *html.Node) {
		if n.Type == html.ElementNode && n.Data == "sync" {
			for _, a := range n.Attr {
				if a.Key == "start" {
					fmt.Println(a.Val)
					break
				}
			}
		}
		if n.Type == html.TextNode {
			fmt.Println(n.Data)
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			parseNode(c)
		}
	}
	parseNode(doc)

	f, err := os.Create("./res.html")
	if err != nil {
		log.Fatal(err)
	}
	if err := html.Render(f, doc); err != nil {
		log.Fatal(err)
	}
}

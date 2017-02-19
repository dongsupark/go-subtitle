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

package backend

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"path"
	"strings"

	"github.com/spf13/cobra"

	"github.com/dongsupark/go-subtitle/parser"
	"github.com/dongsupark/go-subtitle/subtitle"
)

type Flags struct {
	WebPort string
}

var (
	globalFlags = Flags{}

	webCommand = &cobra.Command{
		Use:   "go-subtitle-web",
		Short: "A simple web app for editing a subtitle file",
		Long: `A simple web app for editing a subtitle file

To get help about a resource or command, please run "go-subtitle-web help"`,
		Run: CommandFunc,
	}
)

const fixtureDir = "./tests/fixtures"

func init() {
	webCommand.PersistentFlags().StringVarP(&globalFlags.WebPort, "port", "p", "8000", "port to serve the web interface")
}

func Execute() {
	if err := webCommand.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
}

func CommandFunc(cmd *cobra.Command, args []string) {
	http.HandleFunc("/readsub/", handleReadSub)

	lAddr := fmt.Sprintf("localhost:%s", globalFlags.WebPort)
	log.Printf("started serving on %q", lAddr)
	if err := http.ListenAndServe(lAddr, nil); err != nil {
		log.Fatalf("unable to listen and service (%v)", err)
		os.Exit(0)
	}
}

func handleReadSub(w http.ResponseWriter, r *http.Request) {
	fname := strings.SplitN(r.URL.Path, "/", 3)[2]

	st, err := readSubFromFile(path.Join(fixtureDir, fname))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	if err := json.NewEncoder(w).Encode(st); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func readSubFromFile(readFileName string) (*subtitle.Subtitle, error) {
	if len(readFileName) == 0 {
		return nil, fmt.Errorf("empty file name")
	}

	if _, err := os.Stat(readFileName); os.IsNotExist(err) {
		return nil, fmt.Errorf("file not found: %s\n", readFileName)
	}

	outSt, err := parser.ReadSubFromFile(readFileName)
	if err != nil {
		return nil, fmt.Errorf("parse error reading %s: %v\n", readFileName, err)
	}

	return outSt, nil
}

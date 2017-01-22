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
)

var (
	goSubtitleCmd = &cobra.Command{
		Use:   "go-subtitle",
		Short: "A command line client for parsing a subtitle file",
		Long: `A CLI for parsing a subtitle file

To get help about a resource or command, please run "go-subtitle help"`,
	}

	globalFlags = struct {
		format string
	}{}
)

func init() {
	goSubtitleCmd.PersistentFlags().StringVar(&globalFlags.format, "format", "subrip", "Subtitle format")
}

func Execute() {
	if err := goSubtitleCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
}

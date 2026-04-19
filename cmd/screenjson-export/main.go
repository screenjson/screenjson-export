// screenjson-export is the free, open-source reference CLI that converts
// Final Draft (.fdx), Fountain (.fountain, .spmd), and FadeIn (.fadein)
// screenplays into ScreenJSON documents.
//
// It does one thing: reads a supported writer format and emits a valid
// ScreenJSON document on stdout or to a file. For the full toolchain —
// PDF import, exports, validation, encryption, REST API, storage
// drivers — see https://screenjson.com/tools/screenjson-cli/
package main

import (
	"context"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	fadeinbridge "screenjson/export/internal/formats/fadein/bridge"
	fadeincodec "screenjson/export/internal/formats/fadein/codec"
	fdxbridge "screenjson/export/internal/formats/fdx/bridge"
	fdxcodec "screenjson/export/internal/formats/fdx/codec"
	fountainbridge "screenjson/export/internal/formats/fountain/bridge"
	fountaincodec "screenjson/export/internal/formats/fountain/codec"
	"screenjson/export/internal/model"
)

// version is the CLI version. Injected at build time via
//   go build -ldflags="-X main.version=v1.2.3"
// Defaults to "dev" for local unreleased builds.
var version = "dev"

func main() {
	fs := flag.NewFlagSet("screenjson-export", flag.ExitOnError)

	var (
		input       string
		output      string
		format      string
		lang        string
		pretty      bool
		showVersion bool
		showHelp    bool
	)

	// Short and long forms share a single variable each.
	fs.StringVar(&input, "i", "", "Input file path (.fdx, .fountain, .spmd, or .fadein). Required.")
	fs.StringVar(&input, "input", "", "Alias for -i.")
	fs.StringVar(&output, "o", "", "Output file path. Defaults to stdout.")
	fs.StringVar(&output, "output", "", "Alias for -o.")
	fs.StringVar(&format, "f", "", "Force input format: fdx | fountain | fadein. Auto-detected if omitted.")
	fs.StringVar(&format, "format", "", "Alias for -f.")
	fs.StringVar(&lang, "lang", "en", "BCP 47 primary language tag for the output.")
	fs.BoolVar(&pretty, "pretty", true, "Pretty-print the JSON output.")
	fs.BoolVar(&showVersion, "v", false, "Print version and exit.")
	fs.BoolVar(&showVersion, "version", false, "Alias for -v.")
	fs.BoolVar(&showHelp, "h", false, "Print help and exit.")
	fs.BoolVar(&showHelp, "help", false, "Alias for -h.")

	fs.Usage = printUsage

	args := rewriteLongFlags(os.Args[1:])
	if err := fs.Parse(args); err != nil {
		os.Exit(2)
	}

	if showVersion {
		fmt.Println("screenjson-export", version)
		return
	}
	if showHelp || input == "" {
		printUsage()
		if input == "" {
			os.Exit(2)
		}
		return
	}

	data, err := os.ReadFile(input)
	if err != nil {
		fatal("read input: %v", err)
	}

	detected := format
	if detected == "" {
		detected = detectFormat(input)
		if detected == "" {
			fatal("could not detect input format from extension. Use -f <fdx|fountain|fadein>")
		}
	}

	doc, err := convert(context.Background(), data, detected, lang)
	if err != nil {
		fatal("convert: %v", err)
	}

	var out []byte
	if pretty {
		out, err = json.MarshalIndent(doc, "", "  ")
	} else {
		out, err = json.Marshal(doc)
	}
	if err != nil {
		fatal("marshal: %v", err)
	}

	if output == "" {
		if _, err := os.Stdout.Write(out); err != nil {
			fatal("write stdout: %v", err)
		}
		if _, err := os.Stdout.Write([]byte("\n")); err != nil {
			fatal("write stdout: %v", err)
		}
		return
	}
	if err := os.WriteFile(output, out, 0o644); err != nil {
		fatal("write output: %v", err)
	}
}

func convert(ctx context.Context, data []byte, format, lang string) (*model.Document, error) {
	switch format {
	case "fdx":
		fdx, err := fdxcodec.NewDecoder().Decode(ctx, data)
		if err != nil {
			return nil, err
		}
		return fdxbridge.ToScreenJSON(fdx, lang), nil
	case "fountain":
		ftn, err := fountaincodec.NewDecoder().Decode(ctx, data)
		if err != nil {
			return nil, err
		}
		return fountainbridge.ToScreenJSON(ftn, lang), nil
	case "fadein":
		osf, err := fadeincodec.NewDecoder().Decode(ctx, data)
		if err != nil {
			return nil, err
		}
		return fadeinbridge.ToScreenJSON(osf, lang), nil
	default:
		return nil, errors.New("unsupported format: " + format + " (use fdx, fountain, or fadein)")
	}
}

func detectFormat(path string) string {
	ext := strings.ToLower(filepath.Ext(path))
	switch ext {
	case ".fdx":
		return "fdx"
	case ".fountain", ".spmd":
		return "fountain"
	case ".fadein":
		return "fadein"
	}
	return ""
}

// rewriteLongFlags maps the GNU-style --flag to Go's -flag.
func rewriteLongFlags(in []string) []string {
	out := make([]string, 0, len(in))
	for _, a := range in {
		if strings.HasPrefix(a, "--") {
			out = append(out, a[1:])
		} else {
			out = append(out, a)
		}
	}
	return out
}

func fatal(format string, args ...any) {
	fmt.Fprintf(os.Stderr, "screenjson-export: "+format+"\n", args...)
	os.Exit(1)
}

// printUsage writes the help text to stderr.
func printUsage() {
	fmt.Fprintln(os.Stderr, `screenjson-export — convert screenplay formats to ScreenJSON.

USAGE
    screenjson-export -i <file> [options]

OPTIONS
    -i, --input <path>      Input file (.fdx, .fountain, .spmd, .fadein).   (required)
    -o, --output <path>     Output file. Defaults to stdout.
    -f, --format <name>     Force input format: fdx | fountain | fadein.
        --lang <tag>        BCP 47 primary language tag. Default: en.
        --pretty            Pretty-print JSON. Default: true.
    -v, --version           Print version.
    -h, --help              Show this help.

EXAMPLES
    screenjson-export -i script.fdx -o script.json
    screenjson-export -i script.fountain | jq '.title.en'
    screenjson-export -i screenplay.fadein --lang fr -o screenplay.json

This is the free, open-source reference implementation. It handles
writer-format import only. For full ScreenJSON tooling — PDF import,
exports, validation, encryption, REST API, and storage drivers — see
https://screenjson.com/tools/screenjson-cli/
`)
	_ = io.Discard
}

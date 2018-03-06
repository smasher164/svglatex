// svglatex takes a given latex document and produces an SVG image.
// Usage:
//
//     svglatex [-inline] < foo.tex
//
// It is a wrapper over `latex` and `dvisvgm`and
// makes the following assumptions and promises:
// * latex is already installed and has an
//   "-output-directory" flag (used to store temporary files).
// * dvisvgm is already installed.
// * The latex document is input through Stdin.
// * When successful, it will output the SVG markup to
//   Stdout (everything else is diverted to Stderr).
//
// The -inline flag causes svglatex to assume the input
// is an inline equation, and wraps it in the following:
//
// \documentclass{standalone}
// \begin{document}
//
// %THE INPUT GOES HERE%
//
// \end{document}
//
package main

import (
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"
)

func usage() {
	fmt.Fprintf(os.Stderr, "usage: svglatex [-inline]\n")
	os.Exit(2)
}

func writeTex(w io.Writer, r io.Reader, inline bool) error {
	if inline {
		if _, err := w.Write([]byte("\\documentclass{standalone}\n\\begin{document}\n")); err != nil {
			return err
		}
	}
	if _, err := io.Copy(w, r); err != nil {
		return err
	}
	if inline {
		if _, err := w.Write([]byte("\n\\end{document}")); err != nil {
			return err
		}
	}
	return nil
}

func svglatex(inline bool) error {
	dirName, err := ioutil.TempDir("", "svglatex")
	if err != nil {
		return err
	}
	defer os.RemoveAll(dirName)

	tex := filepath.Join(dirName, "in.tex")
	f, err := os.Create(tex)
	if err != nil {
		return err
	}
	defer f.Close()
	if err := writeTex(f, os.Stdin, inline); err != nil {
		return err
	}

	lacmd := exec.Command("latex", "-output-directory", dirName, tex)
	lacmd.Stdout = os.Stderr
	lacmd.Stderr = os.Stderr
	if err := lacmd.Run(); err != nil {
		return err
	}

	dvi := filepath.Join(dirName, "in.dvi")
	dvcmd := exec.Command("dvisvgm", dvi, "--no-fonts", "--stdout", "--verbosity=1")
	dvcmd.Stdout = os.Stdout
	dvcmd.Stderr = os.Stderr
	if err := dvcmd.Run(); err != nil {
		return err
	}
	return nil
}

func main() {
	log.SetPrefix("svglatex: ")
	log.SetFlags(0)
	flag.Usage = usage
	inline := flag.Bool("inline", false, "May be preferable when passing in inline equations.")
	flag.Parse()
	if err := svglatex(*inline); err != nil {
		log.Fatalln(err)
	}
}

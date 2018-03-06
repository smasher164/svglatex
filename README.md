# svglatex

svglatex takes a given latex document and produces an SVG image.
Usage:
```
go get -u github.com/smasher164/svglatex

svglatex [-inline] < foo.tex
```
It is a wrapper over the `latex` and `dvisvgm` commands and makes the following assumptions and promises:
* `latex` is already installed and has an "-output-directory" flag (used to store temporary files).
* `dvisvgm` is already installed.
* The latex document is input through Stdin.
* When successful, it will output the SVG markup to Stdout (everything else is diverted to Stderr).

The `-inline` flag causes svglatex to assume the input is an inline equation, and wraps it in the following
```latex
\documentclass{standalone}
\begin{document}

%THE INPUT GOES HERE%

\end{document}
```
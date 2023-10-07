package compile

import (
	"bufio"
	_ "embed"
	"fmt"
	ms2 "gwcoffey/otis/ms"
	"gwcoffey/otis/ms/compile/html"
	"gwcoffey/otis/ms/compile/rtf"
	"gwcoffey/otis/ms/compile/tex"
	"gwcoffey/otis/msfs"
	"gwcoffey/otis/text"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

type Format int

type Args struct {
	ProjectPath *string `arg:"positional"`
	Format      string  `arg:"-f" help:"the compiled output format (PDF, RTF, HTML, or TEX)" default:"PDF"`
	Tag         *string `arg:"-t" help:"tag to append to the filename, [default: <current date>]"`
}

func generateTex(fileName string, manuscript ms2.Manuscript) (err error) {
	outDir, err := msfs.DistDir(manuscript.Path())
	if err != nil {
		return
	}

	path := filepath.Join(outDir, fileName+".tex")
	err = os.MkdirAll(filepath.Dir(path), os.ModePerm)
	if err != nil {
		return
	}

	return writeTex(path, manuscript)
}

func generatePdf(fileName string, manuscript ms2.Manuscript) (err error) {
	tmpDir, err := msfs.TmpDir(manuscript.Path())
	if err != nil {
		return
	}

	distDir, err := msfs.DistDir(manuscript.Path())
	if err != nil {
		return
	}

	texPath := filepath.Join(tmpDir, "compile", "tmp-for-pdf.tex")
	err = os.MkdirAll(filepath.Dir(texPath), os.ModePerm)
	if err != nil {
		return
	}

	pdfPath := filepath.Join(tmpDir, "compile", "tmp-for-pdf.pdf")
	err = os.MkdirAll(filepath.Dir(pdfPath), os.ModePerm)
	if err != nil {
		return
	}

	outPath := filepath.Join(distDir, fileName+".pdf")
	err = os.MkdirAll(filepath.Dir(pdfPath), os.ModePerm)
	if err != nil {
		return
	}

	err = writeTex(texPath, manuscript)
	if err != nil {
		return
	}

	err = execPdfLatex(texPath, manuscript)
	if err != nil {
		return
	}

	err = os.Rename(pdfPath, outPath)
	return
}

func generateHtml(fileName string, manuscript ms2.Manuscript) (err error) {
	outDir, err := msfs.DistDir(manuscript.Path())
	if err != nil {
		return
	}

	path := filepath.Join(outDir, fileName+".html")
	err = os.MkdirAll(filepath.Dir(path), os.ModePerm)
	if err != nil {
		return
	}

	file, err := os.Create(path)
	if err != nil {
		return
	}

	w := bufio.NewWriter(file)
	defer func() {
		if err = file.Close(); err != nil {
			return
		}
	}()

	htmlContent, err := html.ManuscriptToHtml(manuscript)
	if err != nil {
		return
	}
	_, err = w.WriteString(htmlContent)
	if err != nil {
		return
	}
	err = w.Flush()

	return
}

func generateRtf(fileName string, manuscript ms2.Manuscript) (err error) {
	outDir, err := msfs.DistDir(manuscript.Path())
	if err != nil {
		return
	}

	path := filepath.Join(outDir, fileName+".rtf")
	err = os.MkdirAll(filepath.Dir(path), os.ModePerm)
	if err != nil {
		return
	}

	file, err := os.Create(path)
	if err != nil {
		return
	}

	w := bufio.NewWriter(file)
	defer func() {
		if err = file.Close(); err != nil {
			return
		}
	}()

	rtfContent, err := rtf.ManuscriptToHtml(manuscript)
	if err != nil {
		return
	}
	_, err = w.WriteString(rtfContent)
	if err != nil {
		return
	}
	err = w.Flush()

	return
}

func execPdfLatex(texPath string, manuscript ms2.Manuscript) (err error) {
	tmpDir, err := msfs.TmpDir(manuscript.Path())
	if err != nil {
		return
	}

	cmd := exec.Command("pdflatex", "-output-directory", filepath.Join(tmpDir, "compile"), texPath)
	cmd.Stdin = strings.NewReader("some input")
	var out strings.Builder
	cmd.Stdout = &out

	err = cmd.Run()
	if err != nil {
		return
	}

	return
}

func writeTex(path string, manuscript ms2.Manuscript) (err error) {
	file, err := os.Create(path)
	if err != nil {
		return
	}

	w := bufio.NewWriter(file)
	defer func() {
		if err = file.Close(); err != nil {
			return
		}
	}()

	latex, err := tex.ManuscriptToTex(manuscript)
	if err != nil {
		return
	}
	_, err = w.WriteString(latex)
	if err != nil {
		return
	}
	err = w.Flush()

	return
}

func Compile(args *Args) (err error) {
	var manuscript ms2.Manuscript

	if &args.ProjectPath == nil {
		manuscript, err = ms2.LoadHere()
	} else {
		manuscript, err = ms2.Load(*args.ProjectPath)
	}
	if err != nil {
		return
	}

	var fileName string
	if args.Tag != nil {
		fileName = fmt.Sprintf("%s-%s", text.ToKebab(manuscript.Title()), *args.Tag)
	} else {
		fileName = fmt.Sprintf("%s-%s", text.ToKebab(manuscript.Title()), time.Now().Format("2006-01-02"))
	}

	switch strings.ToUpper(args.Format) {
	case "PDF":
		err = generatePdf(fileName, manuscript)
	case "RTF":
		err = generateRtf(fileName, manuscript)
	case "HTML":
		err = generateHtml(fileName, manuscript)
	case "TEX":
		err = generateTex(fileName, manuscript)
	}
	if err != nil {
		return
	}

	return nil
}

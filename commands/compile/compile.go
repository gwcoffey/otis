package compile

import (
	"bufio"
	_ "embed"
	"fmt"
	"gwcoffey/otis/shared/ms"
	"gwcoffey/otis/shared/ms/compile/html"
	"gwcoffey/otis/shared/ms/compile/rtf"
	"gwcoffey/otis/shared/ms/compile/tex"
	"gwcoffey/otis/shared/o"
	"gwcoffey/otis/shared/text"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

type Format int

type Args struct {
	WorkName *string `arg:"positional" help:"the work to compile, required if the manuscript has multiple works"`
	Format   string  `arg:"-f" help:"the compiled output format (PDF, RTF, HTML, or TEX)" default:"PDF"`
	Tag      *string `arg:"-t" help:"tag to append to the filename, [default: <current date>]"`
}

func selectWork(args *Args, manuscript ms.Manuscript) ms.Work {
	var work ms.Work
	if args.WorkName != nil {
		for _, w := range manuscript.Works() {
			if filepath.Base(w.Path()) == *args.WorkName {
				work = w
				break
			}
		}
	} else if len(manuscript.Works()) == 1 {
		work = manuscript.Works()[0]
	} else {
		panic("specify a work")
	}
	return work
}

func generateTex(fileName string, otis o.Otis, work ms.Work) (err error) {
	outDir, err := otis.DistDir()
	if err != nil {
		return
	}

	path := filepath.Join(outDir, fileName+".tex")
	err = os.MkdirAll(filepath.Dir(path), os.ModePerm)
	if err != nil {
		return
	}

	return writeTex(path, otis, work)
}

func generatePdf(fileName string, otis o.Otis, work ms.Work) (err error) {
	tmpDir, err := otis.TmpDir()
	if err != nil {
		return
	}

	distDir, err := otis.DistDir()
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

	err = writeTex(texPath, otis, work)
	if err != nil {
		return
	}

	err = execPdfLatex(texPath, otis)
	if err != nil {
		return
	}

	err = os.Rename(pdfPath, outPath)
	return
}

func generateHtml(fileName string, otis o.Otis, work ms.Work) (err error) {
	outDir, err := otis.DistDir()
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

	htmlContent, err := html.WorkToHtml(otis, work)
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

func generateRtf(fileName string, otis o.Otis, work ms.Work) (err error) {
	outDir, err := otis.DistDir()
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

	rtfContent, err := rtf.WorkToRtf(work, otis)
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

func execPdfLatex(texPath string, otis o.Otis) (err error) {
	tmpDir, err := otis.TmpDir()
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

func writeTex(path string, otis o.Otis, work ms.Work) (err error) {
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

	latex, err := tex.WorkToTex(work, otis)
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

func Compile(otis o.Otis, args *Args) {
	manuscript, err := otis.Manuscript()
	if err != nil {
		panic(err)
	}

	work := selectWork(args, manuscript)

	var fileName string
	if args.Tag != nil {
		fileName = fmt.Sprintf("%s-%s", text.ToKebab(work.Title()), *args.Tag)
	} else {
		fileName = fmt.Sprintf("%s-%s", text.ToKebab(work.Title()), time.Now().Format("2006-01-02"))
	}

	switch strings.ToUpper(args.Format) {
	case "PDF":
		err = generatePdf(fileName, otis, work)
	case "RTF":
		err = generateRtf(fileName, otis, work)
	case "HTML":
		err = generateHtml(fileName, otis, work)
	case "TEX":
		err = generateTex(fileName, otis, work)
	}
	if err != nil {
		panic(err)
	}
}

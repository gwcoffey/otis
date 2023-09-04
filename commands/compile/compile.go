package compile

import (
	"bufio"
	_ "embed"
	"fmt"
	"gwcoffey/otis/shared/cfg"
	"gwcoffey/otis/shared/ms"
	"gwcoffey/otis/shared/ms/compile/html"
	"gwcoffey/otis/shared/ms/compile/rtf"
	"gwcoffey/otis/shared/ms/compile/tex"
	"gwcoffey/otis/shared/text"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

type Format int

type Args struct {
	WorkName   *string `arg:"positional" help:"the work to compile, required if the manuscript has multiple works"`
	Submission bool    `help:"compile for submission"`
	Format     string  `help:"the compiled output format (PDF, RTF, HTML, or TEX)" default:"PDF"`
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

func generateTex(fileName string, config cfg.Config, work ms.Work) (err error) {
	outDir, err := config.DistDir()
	if err != nil {
		return
	}

	path := filepath.Join(outDir, fileName+".tex")
	err = os.MkdirAll(filepath.Dir(path), os.ModePerm)
	if err != nil {
		return
	}

	return writeTex(path, config, work)
}

func generatePdf(fileName string, config cfg.Config, work ms.Work) (err error) {
	tmpDir, err := config.TmpDir()
	if err != nil {
		return
	}

	distDir, err := config.DistDir()
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

	err = writeTex(texPath, config, work)
	if err != nil {
		return
	}

	err = execPdfLatex(texPath, config)
	if err != nil {
		return
	}

	err = os.Rename(pdfPath, outPath)
	return
}

func generateHtml(fileName string, config cfg.Config, work ms.Work) (err error) {
	outDir, err := config.DistDir()
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

	htmlContent, err := html.WorkToHtml(config, work)
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

func generateRtf(fileName string, config cfg.Config, work ms.Work) (err error) {
	outDir, err := config.DistDir()
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

	rtfContent, err := rtf.WorkToRtf(work, config)
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

func execPdfLatex(texPath string, config cfg.Config) (err error) {
	tmpDir, err := config.TmpDir()
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

func writeTex(path string, config cfg.Config, work ms.Work) (err error) {
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

	latex, err := tex.WorkToTex(work, config)
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

func Compile(config cfg.Config, args *Args) {
	manuscript, err := ms.Load(filepath.Join(config.ProjectRoot, "manuscript"))
	if err != nil {
		panic(err)
	}

	work := selectWork(args, manuscript)

	fileName := fmt.Sprintf("%s-%s", text.ToKebab(work.Title()), time.Now().Format("2006-01-02"))

	switch strings.ToUpper(args.Format) {
	case "PDF":
		err = generatePdf(fileName, config, work)
	case "RTF":
		err = generateRtf(fileName, config, work)
	case "HTML":
		err = generateHtml(fileName, config, work)
	case "TEX":
		err = generateTex(fileName, config, work)
	}
	if err != nil {
		panic(err)
	}
}

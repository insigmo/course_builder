package docx

import (
	"archive/zip"
	"encoding/xml"
	"fmt"
	"html"
	"io"
	"strings"
)

// Quiz represents a self-check question extracted from a .docx file.
type Quiz struct {
	Question string   `json:"question"`
	Options  []string `json:"options"`
	Answer   int      `json:"answer"`
}

// Result holds the parsed output of a .docx file.
type Result struct {
	HTML    string
	Quizzes []Quiz
}

// ── Minimal OOXML structs ─────────────────────────────────────────

type wDoc struct {
	Body struct {
		Elems []wElem `xml:",any"`
	} `xml:"body"`
}

type wElem struct {
	XMLName xml.Name
	PPr     *wPPr  `xml:"pPr"`
	Runs    []wRun `xml:"r"`
	Rows    []wRow `xml:"tr"`
}

type wPPr struct {
	Style *struct {
		Val string `xml:"val,attr"`
	} `xml:"pStyle"`
	NumPr *struct {
		Ilvl *struct {
			Val string `xml:"val,attr"`
		} `xml:"ilvl"`
		NumID *struct {
			Val string `xml:"val,attr"`
		} `xml:"numId"`
	} `xml:"numPr"`
}

type wRun struct {
	RPr *struct {
		Bold      *struct{} `xml:"b"`
		Italic    *struct{} `xml:"i"`
		Strike    *struct{} `xml:"strike"`
		Underline *struct {
			Val string `xml:"val,attr"`
		} `xml:"u"`
	} `xml:"rPr"`
	Text string    `xml:"t"`
	Tab  *struct{} `xml:"tab"`
}

type wRow struct {
	Cells []struct {
		Paras []wElem `xml:"p"`
	} `xml:"tc"`
}

// Parse reads a .docx file and returns HTML + quizzes.
func Parse(path string) (Result, error) {
	zr, err := zip.OpenReader(path)
	if err != nil {
		return Result{}, fmt.Errorf("open docx zip: %w", err)
	}
	defer zr.Close()

	var xmlData []byte
	for _, f := range zr.File {
		if f.Name == "word/document.xml" {
			rc, err := f.Open()
			if err != nil {
				return Result{}, fmt.Errorf("open document.xml: %w", err)
			}
			xmlData, err = io.ReadAll(rc)
			rc.Close()
			if err != nil {
				return Result{}, fmt.Errorf("read document.xml: %w", err)
			}
			break
		}
	}
	if xmlData == nil {
		return Result{}, fmt.Errorf("word/document.xml not found in %s", path)
	}

	var doc wDoc
	if err := xml.Unmarshal(xmlData, &doc); err != nil {
		return Result{}, fmt.Errorf("unmarshal xml: %w", err)
	}

	return convert(doc.Body.Elems), nil
}

func convert(elems []wElem) Result {
	var parts []string
	var quizzes []Quiz
	var curQ *Quiz
	openList := ""

	closeList := func() {
		if openList != "" {
			parts = append(parts, "</"+openList+">")
			openList = ""
		}
	}
	finalizeQ := func() {
		if curQ != nil && curQ.Question != "" && len(curQ.Options) > 0 {
			if curQ.Answer >= len(curQ.Options) {
				curQ.Answer = 0
			}
			quizzes = append(quizzes, *curQ)
		}
		curQ = nil
	}

	for _, el := range elems {
		switch el.XMLName.Local {
		case "p":
			raw := paragraphText(el)
			if raw == "" {
				continue
			}

			if strings.HasPrefix(raw, "??") {
				closeList()
				finalizeQ()
				q := &Quiz{Question: strings.TrimSpace(raw[2:])}
				curQ = q
				continue
			}

			if curQ != nil {
				if strings.HasPrefix(raw, "+") {
					curQ.Answer = len(curQ.Options)
					curQ.Options = append(curQ.Options, strings.TrimSpace(raw[1:]))
					continue
				}
				if strings.HasPrefix(raw, "-") {
					curQ.Options = append(curQ.Options, strings.TrimSpace(raw[1:]))
					continue
				}
				finalizeQ()
			}

			if lvl := headingLevel(el); lvl > 0 {
				closeList()
				parts = append(parts, fmt.Sprintf("<h%d>%s</h%d>", lvl, html.EscapeString(raw), lvl))
				continue
			}

			if kind := listKind(el); kind != "" {
				inner := innerHTML(el)
				if openList != kind {
					closeList()
					parts = append(parts, "<"+kind+">")
					openList = kind
				}
				if strings.TrimSpace(inner) != "" {
					parts = append(parts, "<li>"+inner+"</li>")
				}
				continue
			}

			closeList()
			inner := innerHTML(el)
			if strings.TrimSpace(inner) != "" {
				parts = append(parts, "<p>"+inner+"</p>")
			}

		case "tbl":
			finalizeQ()
			closeList()
			if t := tableHTML(el); t != "" {
				parts = append(parts, t)
			}
		}
	}

	finalizeQ()
	closeList()

	return Result{
		HTML:    strings.Join(parts, "\n"),
		Quizzes: quizzes,
	}
}

func paragraphText(p wElem) string {
	var sb strings.Builder
	for _, r := range p.Runs {
		sb.WriteString(r.Text)
		if r.Tab != nil {
			sb.WriteString(" ")
		}
	}
	return strings.TrimSpace(sb.String())
}

func innerHTML(p wElem) string {
	var sb strings.Builder
	for _, r := range p.Runs {
		t := html.EscapeString(r.Text)
		if r.Tab != nil {
			t += " "
		}
		if t == "" {
			continue
		}
		if r.RPr != nil {
			if r.RPr.Bold != nil {
				t = "<strong>" + t + "</strong>"
			}
			if r.RPr.Italic != nil {
				t = "<em>" + t + "</em>"
			}
			if r.RPr.Strike != nil {
				t = "<s>" + t + "</s>"
			}
			if r.RPr.Underline != nil && r.RPr.Underline.Val != "none" {
				t = "<u>" + t + "</u>"
			}
		}
		sb.WriteString(t)
	}
	return sb.String()
}

func headingLevel(p wElem) int {
	if p.PPr == nil || p.PPr.Style == nil {
		return 0
	}
	switch strings.ToLower(p.PPr.Style.Val) {
	case "heading1", "1":
		return 2
	case "heading2", "2":
		return 3
	case "heading3", "3":
		return 4
	case "heading4", "4":
		return 5
	case "heading5", "heading6", "5", "6":
		return 6
	}
	return 0
}

func listKind(p wElem) string {
	if p.PPr == nil || p.PPr.NumPr == nil {
		return ""
	}
	return "ul"
}

func tableHTML(tbl wElem) string {
	if len(tbl.Rows) == 0 {
		return ""
	}
	var sb strings.Builder
	sb.WriteString("<table>\n")
	for i, row := range tbl.Rows {
		tag := "td"
		if i == 0 {
			tag = "th"
		}
		sb.WriteString("<tr>\n")
		for _, cell := range row.Cells {
			var lines []string
			for _, p := range cell.Paras {
				if t := paragraphText(p); t != "" {
					lines = append(lines, html.EscapeString(t))
				}
			}
			fmt.Fprintf(&sb, "<%s>%s</%s>\n", tag, strings.Join(lines, "<br>"), tag)
		}
		sb.WriteString("</tr>\n")
	}
	sb.WriteString("</table>")
	return sb.String()
}

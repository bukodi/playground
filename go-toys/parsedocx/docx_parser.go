package parsedocx

import (
	"fmt"
	"github.com/fumiama/go-docx"
	"io"
	"strconv"
	"strings"
)

type Section struct {
	Heading  string
	Level    int
	Index    int
	Parent   *Section
	Children []*Section
	Items    []any
}

func (s *Section) FindChild(heading string) *Section {
	if s == nil {
		return nil
	}
	for _, child := range s.Children {
		if child.Heading == heading {
			return child
		}
	}
	return nil
}

func (s *Section) DumpTOC(w io.Writer, tab string) {
	if s == nil {
		return
	}
	for idx, child := range s.Children {
		fmt.Fprintf(w, "%s %d. %s\n", tab, idx+1, child.Heading)
		child.DumpTOC(w, tab+"  ")
	}
}

func parseToSections(parent *Section, startIdx int, items []any) (nextIdx int) {
	for idx := startIdx; idx < len(items); {
		it := items[idx]
		if p, ok := it.(*docx.Paragraph); ok {
			hl := extractHeadingLevel(p)
			text := p.String()
			if strings.HasPrefix(text, "Appendix") {
				fmt.Printf(" %s\n", text)
			}
			if hl == 0 {
				// Skip non-heading paragraph
				parent.Items = append(parent.Items, p)
				idx++
				continue
			}
			if hl > parent.Level {
				newChild := &Section{
					Heading: text,
					Level:   hl,
					Parent:  parent,
				}
				parent.Children = append(parent.Children, newChild)
				idx = parseToSections(newChild, idx+1, items)
				continue
			} else if hl <= parent.Level {
				return idx
			}
		} else if t, ok := it.(*docx.Table); ok {
			parent.Items = append(parent.Items, t)
			idx++
			continue
		} else {
			idx++
			continue
		}
	}
	return len(items)
}

func extractHeadingLevel(p *docx.Paragraph) int {
	if p.Properties != nil {
		if p.Properties.Style != nil {
			name := p.Properties.Style.Val
			if strings.HasPrefix(name, "Heading") {
				idxtxt := name[len("Heading"):]
				idx, err := strconv.Atoi(idxtxt)
				if err != nil {
					panic(fmt.Errorf("cannot parse heading level from %s", name))
				}
				return idx
			}
		}
	}
	return 0
}

package main

import (
	"github.com/digitorus/pdf"
	"testing"
)

func TestParseDAPPdf(t *testing.T) {
	r, err := pdf.Open("/home/lbukodi/Downloads/Lakcimkartya_202502121714.pdf")
	if err != nil {
		t.Fatalf("%+v", err)
	}

	p1 := r.Page(1)
	c := p1.Content()

	s := ""
	for _, text := range c.Text {
		s = s + text.S
	}
	t.Logf("Text: %s", s)

	//t.Skip("Skip this test")
	//t.Parallel()
}

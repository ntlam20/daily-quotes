package utils

import (
	_ "embed"
	"encoding/base64"
	"fmt"
	"os"
	"strings"
	"time"
)

//go:embed sniglet.ttf
var snigletFont []byte

type Quote struct {
	Q string `json:"q"`
	A string `json:"a"`
}

func wrapText(text string, maxChars int) []string {
	words := strings.Fields(text)
	var lines []string
	var line strings.Builder
	for _, w := range words {
		if line.Len() == 0 {
			line.WriteString(w)
		} else if line.Len()+1+len(w) <= maxChars {
			line.WriteByte(' ')
			line.WriteString(w)
		} else {
			lines = append(lines, line.String())
			line.Reset()
			line.WriteString(w)
		}
	}
	if line.Len() > 0 {
		lines = append(lines, line.String())
	}
	return lines
}

func xmlEscape(s string) string {
	s = strings.ReplaceAll(s, "&", "&amp;")
	s = strings.ReplaceAll(s, "<", "&lt;")
	s = strings.ReplaceAll(s, ">", "&gt;")
	s = strings.ReplaceAll(s, `"`, "&quot;")
	return s
}

func generateSVG(q Quote, date string) (string, error) {
	fontB64 := base64.StdEncoding.EncodeToString(snigletFont)

	pngBytes, err := os.ReadFile("assets/quote-brush-mask.png")
	if err != nil {
		return "", fmt.Errorf("failed to read brush mask: %w", err)
	}
	pngB64 := base64.StdEncoding.EncodeToString(pngBytes)

	const (
		width    = 900
		cx       = width / 2
		headerH  = 72  // height of header section
		quoteH   = 260 // height of quote section
		lineH    = 32
		maxChars = 52
		padTop   = 36
	)
	totalH := headerH + quoteH

	lines := wrapText(q.Q, maxChars)

	// Quote section y-positions (offset by headerH)
	qyMark      := headerH + padTop + 26
	qyFirstLine := qyMark + 6 + lineH
	qyLastLine  := qyFirstLine + (len(lines)-1)*lineH
	qyAuthor    := qyLastLine + 16 + 24
	qyDate      := qyAuthor + 12 + 18

	var sb strings.Builder
	sb.WriteString(fmt.Sprintf(`<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 %d %d" width="%d" height="%d">`, width, totalH, width, totalH))
	sb.WriteString(fmt.Sprintf(`
  <defs>
    <style>
      @font-face {
        font-family: 'Sniglet';
        src: url('data:font/truetype;base64,%s') format('truetype');
      }
      .bg    { fill: #ffffff; }
      .htxt  { fill: #1a1a1a; font-family: 'Sniglet', sans-serif; font-size: 26px; text-anchor: middle; dominant-baseline: middle; }
      .hline { stroke: #e0e0e0; stroke-width: 1; }
      .qt    { fill: #1a1a1a; font-family: 'Sniglet', sans-serif; font-size: 20px; font-style: italic; text-anchor: middle; }
      .mark  { fill: #c8c0b8; font-family: 'Sniglet', sans-serif; font-size: 38px; text-anchor: middle; }
      .au    { fill: #555555; font-family: 'Sniglet', sans-serif; font-size: 16px; text-anchor: middle; }
      .dt    { fill: #999999; font-family: 'Sniglet', sans-serif; font-size: 13px; text-anchor: middle; }
      .brush { opacity: 0.35; }
    </style>
    <filter id="colorize">
      <feFlood flood-color="#c8b89a" flood-opacity="1" result="color"/>
      <feComposite in="color" in2="SourceGraphic" operator="in"/>
    </filter>
  </defs>`, fontB64))

	// Full background
	sb.WriteString(fmt.Sprintf("\n  <rect class=\"bg\" width=\"%d\" height=\"%d\"/>", width, totalH))

	// ── Header section ──
	sb.WriteString(fmt.Sprintf("\n  <line class=\"hline\" x1=\"40\" y1=\"14\" x2=\"860\" y2=\"14\"/>"))
	sb.WriteString(fmt.Sprintf("\n  <text class=\"htxt\" x=\"%d\" y=\"40\">&#10077; Today&#8217;s Quote</text>", cx))
	sb.WriteString(fmt.Sprintf("\n  <line class=\"hline\" x1=\"40\" y1=\"62\" x2=\"860\" y2=\"62\"/>"))

	// ── Quote section ──
	sb.WriteString(fmt.Sprintf("\n  <image class=\"brush\" filter=\"url(#colorize)\" href=\"data:image/png;base64,%s\" x=\"0\" y=\"%d\" width=\"%d\" height=\"%d\"/>", pngB64, headerH, width, quoteH))
	sb.WriteString(fmt.Sprintf("\n  <text class=\"mark\" style=\"font-size:48px\" x=\"%d\" y=\"%d\">&#10077;</text>", cx, qyMark))

	for i, line := range lines {
		y := qyFirstLine + i*lineH
		sb.WriteString(fmt.Sprintf("\n  <text class=\"qt\" x=\"%d\" y=\"%d\">%s</text>", cx, y, xmlEscape(line)))
	}

	sb.WriteString(fmt.Sprintf("\n  <text class=\"au\" x=\"%d\" y=\"%d\">&#8212; %s</text>", cx, qyAuthor, xmlEscape(q.A)))
	sb.WriteString(fmt.Sprintf("\n  <text class=\"dt\" x=\"%d\" y=\"%d\">%s</text>", cx, qyDate, xmlEscape(date)))

	sb.WriteString("\n</svg>")
	return sb.String(), nil
}

func generateHiThereSVG() string {
	fontB64 := base64.StdEncoding.EncodeToString(snigletFont)
	const (
		width  = 900
		height = 72
		cx     = width / 2
	)
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf(`<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 %d %d" width="%d" height="%d">`, width, height, width, height))
	sb.WriteString(fmt.Sprintf(`
  <defs>
    <style>
      @font-face {
        font-family: 'Sniglet';
        src: url('data:font/truetype;base64,%s') format('truetype');
      }
      .bg   { fill: #ffffff; }
      .text { fill: #1a1a1a; font-family: 'Sniglet', sans-serif; font-size: 26px; text-anchor: middle; dominant-baseline: middle; }
      .line { stroke: #e0e0e0; stroke-width: 1; }
    </style>
  </defs>`, fontB64))
	sb.WriteString(fmt.Sprintf("\n  <rect class=\"bg\" width=\"%d\" height=\"%d\"/>", width, height))
	sb.WriteString("\n  <line class=\"line\" x1=\"40\" y1=\"14\" x2=\"860\" y2=\"14\"/>")
	sb.WriteString(fmt.Sprintf("\n  <text class=\"text\" x=\"%d\" y=\"%d\">Hi there &#128075;</text>", cx, height/2))
	sb.WriteString("\n  <line class=\"line\" x1=\"40\" y1=\"62\" x2=\"860\" y2=\"62\"/>")
	sb.WriteString("\n</svg>")
	return sb.String()
}

func UpdateREADME(q Quote) error {
	templateContent, err := os.ReadFile("README.template.md")
	if err != nil {
		return fmt.Errorf("failed to read template: %w", err)
	}

	currentTime := time.Now().UTC().Format("January 2, 2006")

	svg, err := generateSVG(q, currentTime)
	if err != nil {
		return fmt.Errorf("failed to generate SVG: %w", err)
	}
	if err := os.WriteFile("assets/daily-quote.svg", []byte(svg), 0644); err != nil {
		return fmt.Errorf("failed to write assets/daily-quote.svg: %w", err)
	}

	hiThere := generateHiThereSVG()
	if err := os.WriteFile("assets/hi-there.svg", []byte(hiThere), 0644); err != nil {
		return fmt.Errorf("failed to write assets/hi-there.svg: %w", err)
	}

	if err := os.WriteFile("README.md", templateContent, 0644); err != nil {
		return fmt.Errorf("failed to write README.md: %w", err)
	}

	return nil
}

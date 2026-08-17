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

func generateQuoteSVG(q Quote, date string) (string, error) {
	fontB64 := base64.StdEncoding.EncodeToString(snigletFont)

	pngBytes, err := os.ReadFile("assets/quote-brush-mask.png")
	if err != nil {
		return "", fmt.Errorf("failed to read brush mask: %w", err)
	}
	pngB64 := base64.StdEncoding.EncodeToString(pngBytes)

	const (
		width    = 900
		height   = 220
		cx       = width / 2
		lineH    = 32
		maxChars = 52
		padTop   = 36
		padBot   = 20
	)

	lines := wrapText(q.Q, maxChars)

	yMark := padTop + 26
	yFirstLine := yMark + 6 + lineH
	yLastLine := yFirstLine + (len(lines)-1)*lineH
	yAuthor := yLastLine + 16 + 24
	yDate := yAuthor + 12 + 18

	var sb strings.Builder
	sb.WriteString(fmt.Sprintf(`<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 %d %d" width="%d" height="%d">`, width, height, width, height))
	sb.WriteString(fmt.Sprintf(`
  <defs><style>
    @font-face {
      font-family: 'Sniglet';
      src: url('data:font/truetype;base64,%s') format('truetype');
    }
    .bg   { fill: #ffffff; }
    .qt   { fill: #1a1a1a; font-family: 'Sniglet', sans-serif; font-size: 20px; font-style: italic; text-anchor: middle; }
    .mark { fill: #c8c0b8; font-family: 'Sniglet', sans-serif; font-size: 38px; text-anchor: middle; }
    .au   { fill: #555555; font-family: 'Sniglet', sans-serif; font-size: 16px; text-anchor: middle; }
    .dt   { fill: #999999; font-family: 'Sniglet', sans-serif; font-size: 13px; text-anchor: middle; }
    .brush { opacity: 0.35; }
  </style>
  <filter id="colorize">
    <feFlood flood-color="#c8b89a" flood-opacity="1" result="color"/>
    <feComposite in="color" in2="SourceGraphic" operator="in"/>
  </filter>
  </defs>`, fontB64))

	// Background
	sb.WriteString(fmt.Sprintf("\n  <rect class=\"bg\" width=\"%d\" height=\"%d\"/>", width, height))

	// Brush mask PNG overlay
	sb.WriteString(fmt.Sprintf("\n  <image class=\"brush\" filter=\"url(#colorize)\" href=\"data:image/png;base64,%s\" x=\"0\" y=\"0\" width=\"%d\" height=\"%d\"/>", pngB64, width, height))

	// Opening quote mark
	sb.WriteString(fmt.Sprintf("\n  <text class=\"mark\" x=\"%d\" y=\"%d\">&#10077;</text>", cx, yMark))

	// Quote lines
	for i, line := range lines {
		y := yFirstLine + i*lineH
		sb.WriteString(fmt.Sprintf("\n  <text class=\"qt\" x=\"%d\" y=\"%d\">%s</text>", cx, y, xmlEscape(line)))
	}

	// Author
	sb.WriteString(fmt.Sprintf("\n  <text class=\"au\" x=\"%d\" y=\"%d\">&#8212; %s</text>", cx, yAuthor, xmlEscape(q.A)))

	// Date
	sb.WriteString(fmt.Sprintf("\n  <text class=\"dt\" x=\"%d\" y=\"%d\">%s</text>", cx, yDate, xmlEscape(date)))

	sb.WriteString("\n</svg>")
	return sb.String(), nil
}

func UpdateREADME(q Quote) error {
	templateContent, err := os.ReadFile("README.template.md")
	if err != nil {
		return fmt.Errorf("failed to read template: %w", err)
	}

	currentTime := time.Now().UTC().Format("January 2, 2006")

	quoteSVG, err := generateQuoteSVG(q, currentTime)
	if err != nil {
		return fmt.Errorf("failed to generate quote SVG: %w", err)
	}
	if err := os.WriteFile("assets/quote.svg", []byte(quoteSVG), 0644); err != nil {
		return fmt.Errorf("failed to write assets/quote.svg: %w", err)
	}

	if err := os.WriteFile("README.md", templateContent, 0644); err != nil {
		return fmt.Errorf("failed to write README.md: %w", err)
	}

	return nil
}

package handlers

import (
	"ascii-art-web/internal/renderer"
	"bytes"
	"errors"
	"html/template"
	"net/http"
	"strings"
	"sync"
)

var (
	tmpl     *template.Template
	tmplErr  error
	tmplOnce = new(sync.Once)
)

type PageData struct {
	Art       string
	Lines     []renderer.RenderedLine
	Color     string
	Alignment string
	Input     string
}

func Home(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}
	if r.Method != http.MethodGet {
		w.Header().Set("Allow", http.MethodGet)
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	data := PageData{
		Art: "",
	}

	err := renderPage(w, data)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}

func Generate(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.Header().Set("Allow", http.MethodPost)
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}
	text := r.FormValue("text")
	text = strings.ReplaceAll(text, "\r\n", "\n")
	if strings.TrimSpace(text) == "" {
		http.Error(w, "Text cannot be empty", http.StatusBadRequest)
		return
	}
	bannerName := r.FormValue("banner")
	switch bannerName {
	case "standard", "shadow", "thinkertoy":
		// Supported banner.
	default:
		http.Error(w, "Invalid banner", http.StatusBadRequest)
		return
	}
	color := r.FormValue("color")
	color, err := validateColor(color)
	if err != nil {
		http.Error(w, "Invalid color", http.StatusBadRequest)
		return
	}
	alignment := validateAlignment(r.FormValue("alignment"))
	target := r.FormValue("target")
	target = strings.ReplaceAll(target, "\r\n", "\n")
	if len(target) > len(text) {
		http.Error(w, "Target cannot be longer than text", http.StatusBadRequest)
		return
	}
	shapes, err := renderer.GetShapes(bannerName)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	err = validateTarget(target, shapes)
	if err != nil {
		if errors.Is(err, renderer.ErrUnsupportedCharacter) {
			http.Error(w, "Unsupported character", http.StatusBadRequest)
			return
		}
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	data := PageData{Color: color, Alignment: alignment, Input: text}
	if alignment == "justify" || strings.TrimSpace(target) != "" {
		artLines, err := renderer.GenerateBlocks(text, bannerName, target)
		if err != nil {
			if errors.Is(err, renderer.ErrUnsupportedCharacter) {
				http.Error(w, "Unsupported character", http.StatusBadRequest)
				return
			}
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
		data.Lines = artLines

	} else {
		art, err := renderer.GenerateArt(text, bannerName)
		if err != nil {
			if errors.Is(err, renderer.ErrUnsupportedCharacter) {
				http.Error(w, "Unsupported character", http.StatusBadRequest)
				return
			}

			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
		data.Art = art
	}

	err = renderPage(w, data)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

}

func validateColor(color string) (string, error) {
	switch color {
	case "black", "red", "green", "yellow", "blue", "magenta", "cyan", "white":
		return color, nil
	case "":
		return "black", nil
	default:
		return "", errors.New("invalid color")
	}
}
func validateTarget(target string, shapes map[rune][]string) error {
	for _, char := range target {
		if char == '\n' {
			continue
		}
		_, err := renderer.RenderCharacter(shapes, char)
		if err != nil {
			return err
		}
	}
	return nil
}

func validateAlignment(alignment string) string {
	switch alignment {
	case "left", "center", "right", "justify":
		return alignment
	default:
		return "left"
	}
}

func renderPage(w http.ResponseWriter, data PageData) error {
	tmplOnce.Do(func() {
		tmpl, tmplErr = template.ParseFiles("templates/index.html")
	})
	if tmplErr != nil {
		return tmplErr
	}
	if tmpl == nil {
		return errors.New("template not loaded")
	}

	var buf bytes.Buffer
	err := tmpl.Execute(&buf, data)
	if err != nil {
		return err
	}

	_, err = buf.WriteTo(w)
	if err != nil {
		return err
	}
	return nil
}

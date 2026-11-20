package handlers

import (
	"ascii-art-web/internal/banner"
	"ascii-art-web/internal/renderer"
	"errors"
	"html"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"strings"
	"sync"
	"testing"
)

func TestMain(m *testing.M) {
	// If running from internal/handlers directory, change working directory
	// to project root so that templates/ and banners/ can be found.
	if _, err := os.Stat("templates/index.html"); os.IsNotExist(err) {
		if err := os.Chdir("../.."); err != nil {
			panic(err)
		}
	}
	os.Exit(m.Run())
}

func TestHome(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	recorder := httptest.NewRecorder()

	Home(recorder, req)

	if recorder.Code != http.StatusOK {
		t.Errorf("expected status %d, got %d", http.StatusOK, recorder.Code)
	}

	body := recorder.Body.String()
	if !strings.Contains(body, "ASCII Art Generator") {
		t.Errorf("expected response body to contain 'ASCII Art Generator', got %q", body)
	}
	if !strings.Contains(body, "></textarea>") {
		t.Errorf("expected empty textarea on home page, got %q", body)
	}
}

func TestGenerate_MethodNotAllowed(t *testing.T) {
	methods := []string{
		http.MethodGet,
		http.MethodPut,
		http.MethodDelete,
		http.MethodPatch,
	}

	for _, method := range methods {
		t.Run(method, func(t *testing.T) {
			req := httptest.NewRequest(method, "/ascii-art", nil)
			recorder := httptest.NewRecorder()

			Generate(recorder, req)

			if recorder.Code != http.StatusMethodNotAllowed {
				t.Errorf("expected status %d, got %d", http.StatusMethodNotAllowed, recorder.Code)
			}
			if allow := recorder.Header().Get("Allow"); allow != http.MethodPost {
				t.Errorf("expected Allow header %q, got %q", http.MethodPost, allow)
			}
		})
	}
}

func TestGenerate_InvalidBanner(t *testing.T) {
	banners := []string{
		"invalid",
		"standard.txt",
		"",
		"random",
	}

	for _, banner := range banners {
		t.Run("banner_"+banner, func(t *testing.T) {
			formData := url.Values{}
			formData.Set("text", "hello")
			formData.Set("banner", banner)

			req := httptest.NewRequest(http.MethodPost, "/ascii-art", strings.NewReader(formData.Encode()))
			req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
			recorder := httptest.NewRecorder()

			Generate(recorder, req)

			if recorder.Code != http.StatusBadRequest {
				t.Errorf("expected status %d, got %d", http.StatusBadRequest, recorder.Code)
			}

			if !strings.Contains(recorder.Body.String(), "Invalid banner") {
				t.Errorf("expected body to contain 'Invalid banner', got %q", recorder.Body.String())
			}
		})
	}
}

func TestGenerate_EmptyText(t *testing.T) {
	cases := []struct {
		name string
		text string
	}{
		{name: "empty string", text: ""},
		{name: "spaces only", text: "   "},
		{name: "tabs and newlines only", text: "\t\n  "},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			formData := url.Values{}
			formData.Set("text", tc.text)
			formData.Set("banner", "standard")

			req := httptest.NewRequest(http.MethodPost, "/ascii-art", strings.NewReader(formData.Encode()))
			req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
			recorder := httptest.NewRecorder()

			Generate(recorder, req)

			if recorder.Code != http.StatusBadRequest {
				t.Errorf("expected status %d, got %d", http.StatusBadRequest, recorder.Code)
			}

			if !strings.Contains(recorder.Body.String(), "Text cannot be empty") {
				t.Errorf("expected body to contain 'Text cannot be empty', got %q", recorder.Body.String())
			}
		})
	}
}

func TestGenerate_UnsupportedCharacter(t *testing.T) {
	cases := []struct {
		name   string
		text   string
		target string
	}{
		{name: "emoji in text", text: "hello 😀", target: ""},
		{name: "non-ASCII characters in text", text: "hello 世界", target: ""},
		{name: "accented characters in text", text: "café", target: ""},
		{name: "emoji in target", text: "hello", target: "😀"},
		{name: "non-ASCII characters in target", text: "hello world", target: "世界"},
		{name: "accented characters in target", text: "hello", target: "café"},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			formData := url.Values{}
			formData.Set("text", tc.text)
			if tc.target != "" {
				formData.Set("target", tc.target)
			}
			formData.Set("banner", "standard")

			req := httptest.NewRequest(http.MethodPost, "/ascii-art", strings.NewReader(formData.Encode()))
			req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
			recorder := httptest.NewRecorder()

			Generate(recorder, req)

			if recorder.Code != http.StatusBadRequest {
				t.Errorf("expected status %d, got %d", http.StatusBadRequest, recorder.Code)
			}

			if !strings.Contains(recorder.Body.String(), "Unsupported character") {
				t.Errorf("expected body to contain 'Unsupported character', got %q", recorder.Body.String())
			}
		})
	}
}

func TestGenerate_Success(t *testing.T) {
	banners := []string{
		"standard",
		"shadow",
		"thinkertoy",
	}

	for _, banner := range banners {
		t.Run(banner+"_full_art", func(t *testing.T) {
			formData := url.Values{}
			formData.Set("text", "Hello 123!")
			formData.Set("banner", banner)

			req := httptest.NewRequest(http.MethodPost, "/ascii-art", strings.NewReader(formData.Encode()))
			req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
			recorder := httptest.NewRecorder()

			Generate(recorder, req)

			if recorder.Code != http.StatusOK {
				t.Fatalf("expected status %d, got %d, body: %q", http.StatusOK, recorder.Code, recorder.Body.String())
			}

			body := recorder.Body.String()
			if !strings.Contains(body, `<pre class="color-black">`) {
				t.Error("expected generated art with default black color")
			}
			if strings.Contains(body, "<pre></pre>") {
				t.Error("expected generated ASCII art, got an empty result")
			}
			if !strings.Contains(body, `<figure class="art-container align-left">`) {
				t.Error("expected container with default align-left class")
			}
			if !strings.Contains(body, ">Hello 123!</textarea>") {
				t.Error("expected submitted text to be preserved in textarea")
			}
		})

		t.Run(banner+"_with_target", func(t *testing.T) {
			formData := url.Values{}
			formData.Set("text", "Hello 123!")
			formData.Set("target", "123")
			formData.Set("color", "red")
			formData.Set("banner", banner)

			req := httptest.NewRequest(http.MethodPost, "/ascii-art", strings.NewReader(formData.Encode()))
			req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
			recorder := httptest.NewRecorder()

			Generate(recorder, req)

			if recorder.Code != http.StatusOK {
				t.Fatalf("expected status %d, got %d, body: %q", http.StatusOK, recorder.Code, recorder.Body.String())
			}

			body := recorder.Body.String()
			if !strings.Contains(body, `<span class="color-red">`) {
				t.Error("expected target to be wrapped with color-red span")
			}
			if !strings.Contains(body, `<figure class="art-container align-left">`) {
				t.Error("expected container with default align-left class")
			}
			if !strings.Contains(body, ">Hello 123!</textarea>") {
				t.Error("expected submitted text to be preserved in textarea")
			}
		})

		t.Run(banner+"_with_alignment", func(t *testing.T) {
			formData := url.Values{}
			formData.Set("text", "Hello 123!")
			formData.Set("banner", banner)
			formData.Set("alignment", "center")

			req := httptest.NewRequest(http.MethodPost, "/ascii-art", strings.NewReader(formData.Encode()))
			req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
			recorder := httptest.NewRecorder()

			Generate(recorder, req)

			if recorder.Code != http.StatusOK {
				t.Fatalf("expected status %d, got %d, body: %q", http.StatusOK, recorder.Code, recorder.Body.String())
			}

			body := recorder.Body.String()
			if !strings.Contains(body, `<figure class="art-container align-center">`) {
				t.Error("expected container with align-center class")
			}
			if !strings.Contains(body, ">Hello 123!</textarea>") {
				t.Error("expected submitted text to be preserved in textarea")
			}
		})
	}
}

func TestGenerate_PreserveInput(t *testing.T) {
	cases := []struct {
		name string
		text string
	}{
		{name: "standard text", text: "Hello 123!"},
		{name: "multiline text", text: "Hello\nWorld"},
		{name: "special characters", text: "Text with & < > \" ' symbols"},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			formData := url.Values{}
			formData.Set("text", tc.text)
			formData.Set("banner", "standard")

			req := httptest.NewRequest(http.MethodPost, "/ascii-art", strings.NewReader(formData.Encode()))
			req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
			recorder := httptest.NewRecorder()

			Generate(recorder, req)

			if recorder.Code != http.StatusOK {
				t.Fatalf("expected status 200, got %d", recorder.Code)
			}

			body := recorder.Body.String()
			expected := ">" + html.EscapeString(tc.text) + "</textarea>"
			if !strings.Contains(body, expected) {
				t.Errorf("expected textarea to contain %q, not found in body", expected)
			}
		})
	}
}

func TestGenerate_Target(t *testing.T) {
	tests := []struct {
		name           string
		text           string
		target         string
		color          string
		expectedStatus int
		expectedBody   string
		unexpectedBody string
	}{
		{
			name:           "valid target matching substring",
			text:           "Hello World",
			target:         "World",
			color:          "blue",
			expectedStatus: http.StatusOK,
			expectedBody:   `<span class="color-blue">`,
		},
		{
			name:           "valid target with default color black",
			text:           "Hello World",
			target:         "Hello",
			color:          "",
			expectedStatus: http.StatusOK,
			expectedBody:   `<span class="color-black">`,
		},
		{
			name:           "target not matching text",
			text:           "Hello World",
			target:         "xyz",
			color:          "green",
			expectedStatus: http.StatusOK,
			unexpectedBody: `<span class="color-green">`,
		},
		{
			name:           "target longer than text",
			text:           "Short",
			target:         "Much Longer Target",
			color:          "red",
			expectedStatus: http.StatusBadRequest,
			expectedBody:   "Target cannot be longer than text",
		},
		{
			name:           "target longer than multiline text with newlines",
			text:           "A\nB",
			target:         "ABCD",
			color:          "red",
			expectedStatus: http.StatusBadRequest,
			expectedBody:   "Target cannot be longer than text",
		},
		{
			name:           "target with CRLF matching text",
			text:           "Hello\r\nWorld",
			target:         "World",
			color:          "cyan",
			expectedStatus: http.StatusOK,
			expectedBody:   `<span class="color-cyan">`,
		},
		{
			name:           "whitespace only target falls back to Art",
			text:           "Hello World",
			target:         "   ",
			color:          "magenta",
			expectedStatus: http.StatusOK,
			expectedBody:   `<pre class="color-magenta">`,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			formData := url.Values{}
			formData.Set("text", tc.text)
			formData.Set("target", tc.target)
			formData.Set("banner", "standard")
			if tc.color != "" {
				formData.Set("color", tc.color)
			}

			req := httptest.NewRequest(
				http.MethodPost,
				"/ascii-art",
				strings.NewReader(formData.Encode()),
			)
			req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

			recorder := httptest.NewRecorder()

			Generate(recorder, req)

			if recorder.Code != tc.expectedStatus {
				t.Errorf(
					"expected status %d, got %d, body: %q",
					tc.expectedStatus,
					recorder.Code,
					recorder.Body.String(),
				)
			}

			body := recorder.Body.String()
			if tc.expectedBody != "" && !strings.Contains(body, tc.expectedBody) {
				t.Errorf("expected response body to contain %q, got %q", tc.expectedBody, body)
			}
			if tc.unexpectedBody != "" && strings.Contains(body, tc.unexpectedBody) {
				t.Errorf("expected response body to NOT contain %q, got %q", tc.unexpectedBody, body)
			}
		})
	}
}

func TestGenerate_SelectedColoringScenarios(t *testing.T) {
	// Scenario 1: Exact colored span count - World produces 8 colored spans
	t.Run("Exact colored span count", func(t *testing.T) {
		formData := url.Values{}
		formData.Set("text", "Hello World")
		formData.Set("target", "World")
		formData.Set("banner", "standard")
		formData.Set("color", "red")

		req := httptest.NewRequest(http.MethodPost, "/ascii-art", strings.NewReader(formData.Encode()))
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		recorder := httptest.NewRecorder()

		Generate(recorder, req)

		if recorder.Code != http.StatusOK {
			t.Fatalf("expected status 200, got %d, body: %q", recorder.Code, recorder.Body.String())
		}

		body := recorder.Body.String()
		spanOpenCount := strings.Count(body, `<span class="color-red">`)
		spanCloseCount := strings.Count(body, `</span>`)

		if spanOpenCount != 8 {
			t.Fatalf("expected exactly 8 colored spans for 'World', got %d", spanOpenCount)
		}
		if spanCloseCount != 8 {
			t.Fatalf("expected exactly 8 closing spans, got %d", spanCloseCount)
		}
	})

	// Scenario 2: Exact segment content - Only World is colored; Hello remains uncolored
	t.Run("Exact segment content", func(t *testing.T) {
		formData := url.Values{}
		formData.Set("text", "Hello World")
		formData.Set("target", "World")
		formData.Set("banner", "standard")
		formData.Set("color", "red")

		req := httptest.NewRequest(http.MethodPost, "/ascii-art", strings.NewReader(formData.Encode()))
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		recorder := httptest.NewRecorder()

		Generate(recorder, req)

		if recorder.Code != http.StatusOK {
			t.Fatalf("expected status 200, got %d", recorder.Code)
		}

		shapes, err := banner.Load("banners/standard.txt")
		if err != nil {
			t.Fatalf("failed to load banner: %v", err)
		}
		helloRows, err := renderer.RenderLine(shapes, "Hello")
		if err != nil {
			t.Fatalf("failed to render 'Hello': %v", err)
		}
		worldRows, err := renderer.RenderLine(shapes, "World")
		if err != nil {
			t.Fatalf("failed to render 'World': %v", err)
		}

		body := recorder.Body.String()
		// Verify every row has uncolored "Hello" in the word block and colored "World" in the target block
		for i := 0; i < 8; i++ {
			expectedHello := html.EscapeString(helloRows[i])
			if !strings.Contains(body, expectedHello) {
				t.Fatalf("row %d: expected uncolored 'Hello' %q not found in body", i, expectedHello)
			}
			expectedWorld := `<span class="color-red">` + html.EscapeString(worldRows[i]) + `</span>`
			if !strings.Contains(body, expectedWorld) {
				t.Fatalf("row %d: expected colored 'World' %q not found in body", i, expectedWorld)
			}
		}

		// Also verify non-empty rows of "Hello" are not wrapped inside a colored span
		for i := 0; i < 8; i++ {
			trimmed := strings.TrimSpace(helloRows[i])
			if trimmed != "" {
				unwantedPattern := `<span class="color-red">` + html.EscapeString(trimmed)
				if strings.Contains(body, unwantedPattern) {
					t.Fatalf("row %d: 'Hello' was unexpectedly colored", i)
				}
			}
		}
	})

	// Scenario 3: Multiple occurrences - Every occurrence is colored correctly
	t.Run("Multiple occurrences", func(t *testing.T) {
		formData := url.Values{}
		formData.Set("text", "World Hello World")
		formData.Set("target", "World")
		formData.Set("banner", "standard")
		formData.Set("color", "blue")

		req := httptest.NewRequest(http.MethodPost, "/ascii-art", strings.NewReader(formData.Encode()))
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		recorder := httptest.NewRecorder()

		Generate(recorder, req)

		if recorder.Code != http.StatusOK {
			t.Fatalf("expected status 200, got %d", recorder.Code)
		}

		body := recorder.Body.String()
		spanCount := strings.Count(body, `<span class="color-blue">`)
		if spanCount != 16 {
			t.Fatalf("expected 16 colored spans (2 occurrences * 8 rows), got %d", spanCount)
		}

		shapes, err := banner.Load("banners/standard.txt")
		if err != nil {
			t.Fatalf("failed to load banner: %v", err)
		}
		worldRows, err := renderer.RenderLine(shapes, "World")
		if err != nil {
			t.Fatalf("failed to render 'World': %v", err)
		}
		helloRows, err := renderer.RenderLine(shapes, "Hello")
		if err != nil {
			t.Fatalf("failed to render 'Hello': %v", err)
		}

		for i := 0; i < 8; i++ {
			expectedWorld := `<span class="color-blue">` + html.EscapeString(worldRows[i]) + `</span>`
			if !strings.Contains(body, expectedWorld) {
				t.Fatalf("row %d: expected colored 'World' not found in body", i)
			}
			expectedHello := html.EscapeString(helloRows[i])
			if !strings.Contains(body, expectedHello) {
				t.Fatalf("row %d: expected uncolored 'Hello' not found in body", i)
			}
		}
	})

	// Scenario 4: Case-sensitive matching - World does not match world
	t.Run("Case-sensitive matching", func(t *testing.T) {
		formData := url.Values{}
		formData.Set("text", "Hello World")
		formData.Set("target", "world")
		formData.Set("banner", "standard")
		formData.Set("color", "green")

		req := httptest.NewRequest(http.MethodPost, "/ascii-art", strings.NewReader(formData.Encode()))
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		recorder := httptest.NewRecorder()

		Generate(recorder, req)

		if recorder.Code != http.StatusOK {
			t.Fatalf("expected status 200, got %d", recorder.Code)
		}

		body := recorder.Body.String()
		if strings.Contains(body, `<span class="color-`) {
			t.Fatalf("expected no colored spans when target casing does not match, found in body")
		}

		shapes, err := banner.Load("banners/standard.txt")
		if err != nil {
			t.Fatalf("failed to load banner: %v", err)
		}
		helloRows, err := renderer.RenderLine(shapes, "Hello")
		if err != nil {
			t.Fatalf("failed to render 'Hello': %v", err)
		}
		worldRows, err := renderer.RenderLine(shapes, "World")
		if err != nil {
			t.Fatalf("failed to render 'World': %v", err)
		}

		// Ensure all rows are completely uncolored
		for i := 0; i < 8; i++ {
			if !strings.Contains(body, html.EscapeString(helloRows[i])) {
				t.Fatalf("row %d: uncolored text 'Hello' %q not found in body", i, helloRows[i])
			}
			if !strings.Contains(body, html.EscapeString(worldRows[i])) {
				t.Fatalf("row %d: uncolored text 'World' %q not found in body", i, worldRows[i])
			}
		}
	})

	// Scenario 5: Multiline selected coloring - Matches are handled independently on each line
	t.Run("Multiline selected coloring", func(t *testing.T) {
		formData := url.Values{}
		formData.Set("text", "Hello World\nWorld Hello")
		formData.Set("target", "World")
		formData.Set("banner", "standard")
		formData.Set("color", "yellow")

		req := httptest.NewRequest(http.MethodPost, "/ascii-art", strings.NewReader(formData.Encode()))
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		recorder := httptest.NewRecorder()

		Generate(recorder, req)

		if recorder.Code != http.StatusOK {
			t.Fatalf("expected status 200, got %d", recorder.Code)
		}

		body := recorder.Body.String()
		spanCount := strings.Count(body, `<span class="color-yellow">`)
		if spanCount != 16 {
			t.Fatalf("expected 16 colored spans across both lines (1 per line * 8 rows * 2 lines), got %d", spanCount)
		}

		shapes, err := banner.Load("banners/standard.txt")
		if err != nil {
			t.Fatalf("failed to load banner: %v", err)
		}
		helloRows, err := renderer.RenderLine(shapes, "Hello")
		if err != nil {
			t.Fatalf("failed to render 'Hello': %v", err)
		}
		worldRows, err := renderer.RenderLine(shapes, "World")
		if err != nil {
			t.Fatalf("failed to render 'World': %v", err)
		}

		// Verify line 1 and line 2 both contain uncolored Hello and colored World
		for i := 0; i < 8; i++ {
			expectedWorld := `<span class="color-yellow">` + html.EscapeString(worldRows[i]) + `</span>`
			if !strings.Contains(body, expectedWorld) {
				t.Fatalf("row %d: expected colored 'World' not found in body", i)
			}
			expectedHello := html.EscapeString(helloRows[i])
			if !strings.Contains(body, expectedHello) {
				t.Fatalf("row %d: expected uncolored 'Hello' not found in body", i)
			}
		}
	})
}

func TestHome_NotFound(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/unknown", nil)
	recorder := httptest.NewRecorder()

	Home(recorder, req)

	if recorder.Code != http.StatusNotFound {
		t.Errorf("expected status %d, got %d", http.StatusNotFound, recorder.Code)
	}

}

func TestHome_TemplateFailure(t *testing.T) {
	oldTmpl, oldErr, oldOnce := tmpl, tmplErr, tmplOnce
	defer func() {
		tmpl, tmplErr, tmplOnce = oldTmpl, oldErr, oldOnce
	}()
	tmplOnce = new(sync.Once)
	tmpl = nil
	tmplErr = nil

	t.Chdir(t.TempDir())

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()

	Home(rec, req)

	if rec.Code != http.StatusInternalServerError {
		t.Errorf("expected status 500, got %d", rec.Code)
	}
}

func TestValidateColor(t *testing.T) {
	tests := []struct {
		name          string
		inputColor    string
		expectedColor string
		expectError   bool
	}{
		{name: "empty color defaults to black", inputColor: "", expectedColor: "black", expectError: false},
		{name: "black", inputColor: "black", expectedColor: "black", expectError: false},
		{name: "red", inputColor: "red", expectedColor: "red", expectError: false},
		{name: "green", inputColor: "green", expectedColor: "green", expectError: false},
		{name: "yellow", inputColor: "yellow", expectedColor: "yellow", expectError: false},
		{name: "blue", inputColor: "blue", expectedColor: "blue", expectError: false},
		{name: "magenta", inputColor: "magenta", expectedColor: "magenta", expectError: false},
		{name: "cyan", inputColor: "cyan", expectedColor: "cyan", expectError: false},
		{name: "white", inputColor: "white", expectedColor: "white", expectError: false},
		{name: "unsupported color purple", inputColor: "purple", expectedColor: "", expectError: true},
		{name: "unsupported color orange", inputColor: "orange", expectedColor: "", expectError: true},
		{name: "case sensitive uppercase RED", inputColor: "RED", expectedColor: "", expectError: true},
		{name: "numeric value", inputColor: "123", expectedColor: "", expectError: true},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got, err := validateColor(tc.inputColor)
			if tc.expectError {
				if err == nil {
					t.Fatalf("expected error for input %q, got nil", tc.inputColor)
				}
				if got != "" {
					t.Errorf("expected empty string on error, got %q", got)
				}
				return
			}

			if err != nil {
				t.Fatalf("unexpected error for input %q: %v", tc.inputColor, err)
			}
			if got != tc.expectedColor {
				t.Errorf("expected color %q, got %q", tc.expectedColor, got)
			}
		})
	}
}

func TestValidateTarget(t *testing.T) {
	shapes := map[rune][]string{
		'A': {"A0", "A1", "A2", "A3", "A4", "A5", "A6", "A7"},
		'B': {"B0", "B1", "B2", "B3", "B4", "B5", "B6", "B7"},
		' ': {"S0", "S1", "S2", "S3", "S4", "S5", "S6", "S7"},
	}

	tests := []struct {
		name        string
		target      string
		expectError bool
	}{
		{name: "empty target", target: "", expectError: false},
		{name: "single supported character", target: "A", expectError: false},
		{name: "multiple supported characters", target: "AB", expectError: false},
		{name: "supported characters with spaces", target: "A B", expectError: false},
		{name: "supported characters with newline", target: "A\nB", expectError: false},
		{name: "target consisting only of newlines", target: "\n\n", expectError: false},
		{name: "unsupported character", target: "Z", expectError: true},
		{name: "unsupported emoji", target: "😀", expectError: true},
		{name: "valid prefix with unsupported character", target: "AB😀", expectError: true},
		{name: "unsupported character in middle", target: "A Z B", expectError: true},
		{name: "unsupported accented character", target: "café", expectError: true},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			err := validateTarget(tc.target, shapes)
			if tc.expectError {
				if err == nil {
					t.Fatalf("expected error for target %q, got nil", tc.target)
				}
				if !errors.Is(err, renderer.ErrUnsupportedCharacter) {
					t.Errorf("expected ErrUnsupportedCharacter, got %v", err)
				}
				return
			}

			if err != nil {
				t.Fatalf("unexpected error for target %q: %v", tc.target, err)
			}
		})
	}
}

func TestGenerate_Color(t *testing.T) {
	tests := []struct {
		name           string
		color          string
		expectedStatus int
		expectedBody   string
	}{
		{
			name:           "valid red color",
			color:          "red",
			expectedStatus: http.StatusOK,
			expectedBody:   `class="color-red"`,
		},
		{
			name:           "empty color defaults to black",
			color:          "",
			expectedStatus: http.StatusOK,
			expectedBody:   `class="color-black"`,
		},
		{
			name:           "invalid purple color",
			color:          "purple",
			expectedStatus: http.StatusBadRequest,
			expectedBody:   "Invalid color",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			formData := url.Values{}
			formData.Set("text", "Hello")
			formData.Set("banner", "standard")
			formData.Set("color", tc.color)

			req := httptest.NewRequest(
				http.MethodPost,
				"/ascii-art",
				strings.NewReader(formData.Encode()),
			)
			req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

			recorder := httptest.NewRecorder()

			Generate(recorder, req)

			if recorder.Code != tc.expectedStatus {
				t.Errorf(
					"expected status %d, got %d",
					tc.expectedStatus,
					recorder.Code,
				)
			}

			body := recorder.Body.String()

			if !strings.Contains(body, tc.expectedBody) {
				t.Errorf(
					"expected response body to contain %q",
					tc.expectedBody,
				)
			}
		})
	}
}

func TestValidateAlignment(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{input: "left", expected: "left"},
		{input: "center", expected: "center"},
		{input: "right", expected: "right"},
		{input: "justify", expected: "justify"},
		{input: "", expected: "left"},
		{input: "unknown", expected: "left"},
	}

	for _, tc := range tests {
		t.Run(tc.input, func(t *testing.T) {
			got := validateAlignment(tc.input)
			if got != tc.expected {
				t.Errorf("validateAlignment(%q) = %q, expected %q", tc.input, got, tc.expected)
			}
		})
	}
}

func TestGenerate_Alignment(t *testing.T) {
	tests := []struct {
		name         string
		alignment    string
		expectedBody string
	}{
		{
			name:         "center alignment",
			alignment:    "center",
			expectedBody: `class="art-container align-center"`,
		},
		{
			name:         "right alignment",
			alignment:    "right",
			expectedBody: `class="art-container align-right"`,
		},
		{
			name:         "justify alignment",
			alignment:    "justify",
			expectedBody: `class="art-container align-justify"`,
		},
		{
			name:         "left alignment",
			alignment:    "left",
			expectedBody: `class="art-container align-left"`,
		},
		{
			name:         "empty alignment defaults to left",
			alignment:    "",
			expectedBody: `class="art-container align-left"`,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			formData := url.Values{}
			formData.Set("text", "Hello")
			formData.Set("banner", "standard")
			if tc.alignment != "" {
				formData.Set("alignment", tc.alignment)
			}

			req := httptest.NewRequest(
				http.MethodPost,
				"/ascii-art",
				strings.NewReader(formData.Encode()),
			)
			req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

			recorder := httptest.NewRecorder()

			Generate(recorder, req)

			if recorder.Code != http.StatusOK {
				t.Fatalf("expected status 200, got %d", recorder.Code)
			}

			body := recorder.Body.String()
			if !strings.Contains(body, tc.expectedBody) {
				t.Errorf("expected response body to contain %q, got %q", tc.expectedBody, body)
			}
		})
	}
}

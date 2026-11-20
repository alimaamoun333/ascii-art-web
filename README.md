# ASCII Art Web

A web-based graphical interface for generating ASCII art banners from text inputs with advanced styling, color customization, alignment options, and targeted substring highlighting, built entirely with the Go standard library.

---

## Table of Contents

- [Overview](#overview)
- [Features](#features)
- [Project Architecture](#project-architecture)
  - [System Flow](#system-flow)
  - [Package Breakdown](#package-breakdown)
- [Endpoints & Request Reference](#endpoints--request-reference)
  - [HTTP Routes](#http-routes)
  - [Form Parameters](#form-parameters)
  - [HTTP Status Codes & Error Handling](#http-status-codes--error-handling)
- [Usage Examples](#usage-examples)
  - [Web Interface](#web-interface)
  - [cURL / Command Line](#curl--command-line)
- [Prerequisites](#prerequisites)
- [Getting Started](#getting-started)
- [Running Tests](#running-tests)
- [Project Structure](#project-structure)

---

## Overview

**ASCII Art Web** is an HTTP server application written in Go that transforms standard text into decorative ASCII art using standard ASCII banner font templates. Users interact with the application via an HTML interface or directly through HTTP POST requests to customize banner styles, text alignment, foreground colors, and targeted substring coloring.

The project is built without external third-party frameworks, relying strictly on Go's standard library packages (`net/http`, `html/template`, `os`, `strings`, `sync`, `bytes`, etc.).

---

## Features

- **Interactive Web Interface**: Semantic HTML5 form for user input, font selection, alignment, color customization, and instant generation.
- **Multiple Banner Styles**:
  - `standard`
  - `shadow`
  - `thinkertoy`
- **Text Alignment**:
  - `left` (default)
  - `center`
  - `right`
  - `justify` (computes inter-word justification spaces and flexbox layout across lines)
- **Color Customization**:
  - Supports 8 standard terminal colors: `black`, `red`, `green`, `yellow`, `blue`, `magenta`, `cyan`, and `white`.
  - Contrast optimization: Light colors automatically apply a dark contrast background (`#1F2937`) for readability.
- **Targeted Substring Coloring**:
  - Color the entire banner or highlight specific target words, substrings, or individual characters.
  - Case-sensitive matching with support for repeated occurrences, partial-word coloring, and multiline targets.
- **Multiline & Line Break Support**:
  - Handles line breaks and normalizes CRLF (`\r\n`) and LF (`\n`).
  - Correctly renders consecutive blank lines.
- **Strict Input Validation**:
  - Rejects empty strings or whitespace-only inputs.
  - Restricts text and target to printable ASCII characters (decimal 32 to 126) and newline characters.
  - Ensures target substring does not exceed the text input length.
  - Validates banner font and color selections against allowed options.
- **Robust HTTP Architecture**:
  - Standard status codes (`200 OK`, `400 Bad Request`, `404 Not Found`, `405 Method Not Allowed`, `500 Internal Server Error`).
  - Includes appropriate `Allow` response headers on `405 Method Not Allowed`.
  - Thread-safe banner loading and template parsing via `sync.Once`.
  - Buffered template execution (`bytes.Buffer`) to prevent partial or corrupted responses on render errors.
- **Comprehensive Test Suite**:
  - Extensive unit tests covering routing, handlers, ASCII rendering, character boundary calculations, justification layouts, and banner file parsing.

---

## Project Architecture

### System Flow

```
                      +-------------------+
                      |    Web Browser    |
                      +---------+---------+
                                |
                        HTTP GET / POST
                                v
                      +-------------------+
                      |   HTTP Handlers   | (internal/handlers)
                      +---------+---------+
                                |
             +------------------+------------------+
             |                                     |
   (Standard Rendering)                 (Justify / Targeted)
             v                                     v
+------------------------+             +------------------------+
|    internal/renderer   |             |   internal/justify     |
|   - GenerateArt        |             |   - MarkColored        |
+-----------+------------+             |   - SplitTextBlocks    |
            |                          +-----------+------------+
            |                                      |
            |                          +-----------v------------+
            |                          |    internal/layout     |
            |                          |   - CalculateBoundaries|
            |                          |   - MergeRanges        |
            |                          +-----------+------------+
            |                                      |
            |                          +-----------v------------+
            |                          |   internal/renderer    |
            |                          |   - GenerateBlocks     |
            |                          |   - RenderSelected     |
            |                          +-----------+------------+
            |                                      |
            +------------------+-------------------+
                               v
                     +-------------------+
                     |   Banner Loader   | (internal/banner)
                     +---------+---------+
                               |
                         banners/*.txt
```

### Package Breakdown

- **`main.go`**: Configures the HTTP multiplexer (`http.ServeMux`), mounts static file serving on `/static/`, registers routes, and starts the server on port `:8080`.
- **`internal/handlers`**:
  - `Home`: Handles `GET /` to render the initial webpage.
  - `Generate`: Handles `POST /ascii-art`, validates form inputs (`text`, `banner`, `alignment`, `color`, `target`), invokes appropriate renderer pipelines, and executes the HTML template safely using an in-memory buffer.
- **`internal/renderer`**:
  - Loads and caches banner shapes once thread-safely using `sync.Once`.
  - Converts text lines into 8-row ASCII glyphs (`RenderCharacter`, `RenderLine`, `Render`, `GenerateArt`).
  - Provides segment-based rendering (`SegmentRow`, `RenderSelected`, `RenderBlocks`, `GenerateBlocks`) for targeted coloring and block-based justification.
- **`internal/justify`**:
  - Tokenizes lines into words and spacing gaps (`LineLayout`, `ParseLine`).
  - Marks character match masks (`MarkColored`).
  - Splits lines into `TextBlock` units, identifying internal justification spaces.
- **`internal/layout`**:
  - Calculates exact pixel/column boundary widths for character glyphs row by row (`CalculateCharacterBoundaries`, `CalculateMultilineBoundaries`).
  - Identifies substring match positions (`FindRanges`, `MergeRanges`).
  - Maps character indices to rendered row column spans (`CharacterRangeToRenderedRange`).
- **`internal/banner`**:
  - Parses banner font files (`standard.txt`, `shadow.txt`, `thinkertoy.txt`) into memory indexed by `rune` (ASCII 32–126).
- **`templates/`**:
  - `index.html`: Go HTML template with semantic layout, options form, and result display supporting both monolithic ASCII blocks and flexbox block lines.
- **`static/css/`**:
  - `style.css`: Presentation styles for font colors, dark contrast backgrounds, flexbox containers, and alignment rules (`align-left`, `align-center`, `align-right`, `align-justify`).

---

## Endpoints & Request Reference

### HTTP Routes

| Endpoint | Method | Description |
|----------|--------|-------------|
| `/` | `GET` | Serves the main generator form |
| `/ascii-art` | `POST` | Processes text and rendering options, returns the rendered HTML page |
| `/static/*` | `GET` | Serves static assets such as CSS stylesheets |

### Form Parameters

When submitting a `POST` request to `/ascii-art`, use `application/x-www-form-urlencoded`:

| Field | Type | Required | Allowed Values | Default | Description |
|-------|------|----------|----------------|---------|-------------|
| `text` | string | **Yes** | Printable ASCII (32–126) & newlines | None | The input text to convert into ASCII art. Must not be empty or whitespace-only. |
| `banner` | string | **Yes** | `standard`, `shadow`, `thinkertoy` | `standard` | The banner font style to use. |
| `alignment` | string | No | `left`, `center`, `right`, `justify` | `left` | Alignment of the rendered ASCII art. |
| `color` | string | No | `black`, `red`, `green`, `yellow`, `blue`, `magenta`, `cyan`, `white` | `black` | Foreground color of the text or targeted substring. |
| `target` | string | No | Printable ASCII (32–126) & newlines | `""` (empty) | Substring to color. If empty, the entire text is colored. Case-sensitive. |

### HTTP Status Codes & Error Handling

- **`200 OK`**: Request successfully processed and page rendered.
- **`400 Bad Request`**:
  - Text input is missing, empty, or whitespace-only (`"Text cannot be empty"`).
  - Selected banner font is unrecognized (`"Invalid banner"`).
  - Selected color is unrecognized (`"Invalid color"`).
  - Target string length exceeds the total text length (`"Target cannot be longer than text"`).
  - Input text or target contains characters outside printable ASCII (32–126) (`"Unsupported character"`).
- **`404 Not Found`**:
  - Request path is unmapped or non-root (e.g. `GET /unknown`, `GET /hello`).
- **`405 Method Not Allowed`**:
  - Method other than `GET` sent to `/` (e.g. `POST /`).
  - Method other than `POST` sent to `/ascii-art` (e.g. `GET /ascii-art`).
  - Includes `Allow` response header specifying permitted methods.
- **`500 Internal Server Error`**:
  - Failure loading or parsing banner files or executing HTML templates.

---

## Usage Examples

### Web Interface

1. Start the server (see [Getting Started](#getting-started)) and open `http://localhost:8080` in your web browser.
2. Enter your text into the textarea (e.g., `Hello World`).
3. Select your desired banner style (`standard`, `shadow`, or `thinkertoy`).
4. Select an alignment (`left`, `center`, `right`, or `justify`).
5. Choose a text color, and optionally specify a substring in **Text to color** to highlight specific words.
6. Click **Generate ASCII Art** to render your banner.

### cURL / Command Line

- **Load home page:**
  ```bash
  curl -i http://localhost:8080/
  ```

- **Generate basic ASCII art (`standard` font, left-aligned):**
  ```bash
  curl -i -X POST http://localhost:8080/ascii-art \
       -d "text=Hello+World&banner=standard"
  ```

- **Generate colored ASCII art with center alignment:**
  ```bash
  curl -i -X POST http://localhost:8080/ascii-art \
       -d "text=Antigravity&banner=shadow&color=cyan&alignment=center"
  ```

- **Targeted substring coloring with justify alignment:**
  ```bash
  curl -i -X POST http://localhost:8080/ascii-art \
       -d "text=Hello+Beautiful+World&banner=standard&color=red&target=Beautiful&alignment=justify"
  ```

- **Multiline text input:**
  ```bash
  curl -i -X POST http://localhost:8080/ascii-art \
       --data-urlencode "text=Hello"$'\n'"World" \
       -d "banner=standard"
  ```

- **Inspect error: Unsupported character (400 Bad Request):**
  ```bash
  curl -i -X POST http://localhost:8080/ascii-art \
       -d "text=Hello+🚀&banner=standard"
  ```

- **Inspect error: Invalid banner (400 Bad Request):**
  ```bash
  curl -i -X POST http://localhost:8080/ascii-art \
       -d "text=Hello&banner=comic-sans"
  ```

- **Inspect error: Method Not Allowed (405):**
  ```bash
  curl -i -X GET http://localhost:8080/ascii-art
  ```

- **Inspect error: Not Found (404):**
  ```bash
  curl -i http://localhost:8080/nonexistent
  ```

---

## Prerequisites

- [Go](https://go.dev/dl/) (version 1.20 or later).

---

## Getting Started

1. **Clone the repository:**
   ```bash
   git clone https://01.nexteraeducation.net/git/amaamoun/ascii-art-web.git
   cd ascii-art-web
   ```

2. **Run the server:**
   ```bash
   go run .
   ```

3. **Access the application:**
   Open [http://localhost:8080](http://localhost:8080) in your browser.

---

## Running Tests

Run all unit tests across the entire project:

```bash
go test -v ./...
```

Run tests with code coverage analysis:

```bash
go test -cover ./...
```

Run tests for specific packages:

- **Root router tests:**
  ```bash
  go test -v .
  ```

- **HTTP handlers tests:**
  ```bash
  go test -v ./internal/handlers
  ```

- **ASCII renderer tests:**
  ```bash
  go test -v ./internal/renderer
  ```

- **Text justification tests:**
  ```bash
  go test -v ./internal/justify
  ```

- **Glyph layout & boundary tests:**
  ```bash
  go test -v ./internal/layout
  ```

- **Banner parser tests:**
  ```bash
  go test -v ./internal/banner
  ```

---

## Project Structure

```text
ascii-art-web/
├── banners/
│   ├── shadow.txt          # Shadow banner font glyphs
│   ├── standard.txt        # Standard banner font glyphs
│   └── thinkertoy.txt      # Thinkertoy banner font glyphs
├── internal/
│   ├── banner/             # Banner file loading and ASCII rune parsing
│   │   ├── banner.go
│   │   └── banner_test.go
│   ├── handlers/           # HTTP handlers, validation, and template execution
│   │   ├── handlers.go
│   │   └── handlers_test.go
│   ├── justify/            # Line tokenization and justification logic
│   │   ├── justify.go
│   │   └── justify_test.go
│   ├── layout/             # Substring range matching and character width boundaries
│   │   ├── layout.go
│   │   └── layout_test.go
│   └── renderer/           # ASCII art rendering and segment generation
│       ├── renderer.go
│       └── renderer_test.go
├── static/
│   └── css/
│       └── style.css       # Visual styling, colors, and layout flexbox rules
├── templates/
│   └── index.html          # HTML5 form and result template
├── go.mod                  # Go module definition
├── main.go                 # Server entry point and router setup
├── main_test.go            # Top-level routing and method validation tests
└── README.md               # Project documentation
```

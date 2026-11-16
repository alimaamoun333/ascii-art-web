# ASCII Art Web

A web-based graphical interface for generating ASCII art banners from text inputs, built with the Go standard library.

---

## Table of Contents

- [Overview](#overview)
- [Features](#features)
- [Project Architecture](#project-architecture)
- [Endpoints & Status Codes](#endpoints--status-codes)
  - [HTTP Status Codes](#http-status-codes)
  - [Error Conditions](#error-conditions)
- [Usage Examples](#usage-examples)
  - [Web Interface](#web-interface)
  - [cURL / Command Line](#curl--command-line)
- [Prerequisites](#prerequisites)
- [Getting Started](#getting-started)
- [Running Tests](#running-tests)
- [Project Structure](#project-structure)

---

## Overview

**ASCII Art Web** is an HTTP server application that converts user-provided text into decorative ASCII art using standard ASCII banner font templates. Users interact with the service through an HTML interface to submit text and choose between available banner styles.

The application is built cleanly without external frameworks, relying solely on Go's standard library packages (`net/http`, `html/template`, `os`, `strings`, `sync`, etc.).

---

## Features

- **Interactive Web Interface**: Semantic HTML5 form to enter text and select banner fonts.
- **Multiple Banner Styles**:
  - `standard`
  - `shadow`
  - `thinkertoy`
- **Multiline Rendering**: Supports text containing line breaks (`\r\n` / `\n`) and consecutive empty lines.
- **Strict Input Validation**:
  - Rejects empty strings or whitespace-only inputs.
  - Validates characters against printable ASCII (decimal 32 to 126) and newlines.
  - Validates banner names against known styles.
- **Robust HTTP Error Handling**:
  - Correct HTTP status codes (`200 OK`, `400 Bad Request`, `404 Not Found`, `405 Method Not Allowed`, `500 Internal Server Error`).
  - Sets appropriate `Allow` response headers on `405 Method Not Allowed`.
  - Buffered template rendering prevents partial responses and duplicate HTTP status writes.
- **Comprehensive Unit Testing**: High test coverage across router, handlers, renderer, and banner parser.

---

## Project Architecture

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
                    +---------v---------+
                    |  ASCII Renderer   | (internal/renderer)
                    +---------+---------+
                              |
                    +---------v---------+
                    |   Banner Loader   | (internal/banner)
                    +---------+---------+
                              |
                        banners/*.txt
```

- **`main.go`**: Configures routes on `http.ServeMux` and starts the HTTP server on port `:8080`.
- **`internal/handlers`**: Manages HTTP request lifecycles, validates method & input, executes templates safely via buffer, and returns appropriate status codes.
- **`internal/renderer`**: Parses input lines and combines 8-line-high character glyphs from the loaded banner map.
- **`internal/banner`**: Reads and parses raw banner files into memory indexed by `rune` (ASCII 32–126).
- **`templates/`**: HTML template rendered by the server.
- **`banners/`**: Source banner font files (`standard.txt`, `shadow.txt`, `thinkertoy.txt`).

---

## Endpoints & Status Codes

### HTTP Status Codes

| Endpoint | Method | Description | Possible Status Codes |
|----------|--------|-------------|-----------------------|
| `/` | `GET` | Renders the home page with the input form | `200 OK` (success)<br>`404 Not Found` (non-root URL)<br>`405 Method Not Allowed` (non-GET)<br>`500 Internal Server Error` (template failure) |
| `/ascii-art` | `POST` | Processes input text and banner, renders result | `200 OK` (success)<br>`400 Bad Request` (invalid input or banner)<br>`405 Method Not Allowed` (non-POST)<br>`500 Internal Server Error` (server error) |

### Error Conditions

- **`400 Bad Request`**:
  - Banner name is missing or invalid (supported: `standard`, `shadow`, `thinkertoy`).
  - Text input is empty or consists solely of whitespace.
  - Text contains characters outside the printable ASCII range (`32`–`126`).
- **`404 Not Found`**:
  - Request sent to an unmapped path (e.g., `/unknown`, `/about`, `/foo`).
- **`405 Method Not Allowed`**:
  - Method other than `GET` sent to `/` (e.g., `POST /`).
  - Method other than `POST` sent to `/ascii-art` (e.g., `GET /ascii-art`, `PUT /ascii-art`).
  - The response includes the `Allow` header indicating supported methods.
- **`500 Internal Server Error`**:
  - Failure reading or parsing banner files or HTML templates.

---

## Usage Examples

### Web Interface

1. Start the server and navigate to `http://localhost:8080` in your browser.
2. Enter your text into the text area (e.g., `Hello World`).
3. Select a banner style (`standard`, `shadow`, or `thinkertoy`).
4. Click **Generate ASCII Art** to view the rendered graphic.

### cURL / Command Line

- **Load home page:**
  ```bash
  curl -i http://localhost:8080/
  ```

- **Generate ASCII art (`standard` banner):**
  ```bash
  curl -i -X POST http://localhost:8080/ascii-art \
       -d "text=Hello&banner=standard"
  ```

- **Generate ASCII art (`shadow` banner):**
  ```bash
  curl -i -X POST http://localhost:8080/ascii-art \
       -d "text=Go+Language&banner=shadow"
  ```

- **Trigger 400 Bad Request (unsupported character):**
  ```bash
  curl -i -X POST http://localhost:8080/ascii-art \
       -d "text=Hello 🚀&banner=standard"
  ```

- **Trigger 404 Not Found:**
  ```bash
  curl -i http://localhost:8080/unknown
  ```

- **Trigger 405 Method Not Allowed:**
  ```bash
  curl -i -X GET http://localhost:8080/ascii-art
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

3. **Open in browser:**
   ```
   http://localhost:8080
   ```

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

- **Banner parser tests:**
  ```bash
  go test -v ./internal/banner
  ```

---

## Project Structure

```text
ascii-art-web/
├── banners/
│   ├── shadow.txt
│   ├── standard.txt
│   └── thinkertoy.txt
├── internal/
│   ├── banner/
│   │   ├── banner.go
│   │   └── banner_test.go
│   ├── handlers/
│   │   ├── handlers.go
│   │   └── handlers_test.go
│   └── renderer/
│       ├── renderer.go
│       └── renderer_test.go
├── templates/
│   └── index.html
├── go.mod
├── main.go
├── main_test.go
└── README.md
```

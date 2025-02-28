package server

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	common "ocr-api/common/utils"
	"os/exec"
)

func (s *Server) RegisterRoutes() http.Handler {
	mux := http.NewServeMux()

	// Register routes
	mux.HandleFunc("/pdf-parser", s.PDFParser)

	// Wrap the mux with CORS middleware
	return s.corsMiddleware(mux)
}

func (s *Server) corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Set CORS headers
		w.Header().Set("Access-Control-Allow-Origin", "*") // Replace "*" with specific origins if needed
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS, PATCH")
		w.Header().Set("Access-Control-Allow-Headers", "Accept, Authorization, Content-Type, X-CSRF-Token")
		w.Header().Set("Access-Control-Allow-Credentials", "false") // Set to "true" if credentials are required

		// Handle preflight OPTIONS requests
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		// Proceed with the next handler
		next.ServeHTTP(w, r)
	})
}

func (s *Server) PDFParser(w http.ResponseWriter, r *http.Request) {
	// Parse the request body
	if err := r.ParseMultipartForm(10 << 20); err != nil {
		http.Error(w, "Failed to parse the request body", http.StatusBadRequest)
		return
	}
	files := r.MultipartForm.File["files[]"]
	var responseText []map[string]string
	for _, file := range files {
		file, err := file.Open()
		if err != nil {
			http.Error(w, "Failed to open the file", http.StatusInternalServerError)
			return
		}
		defer file.Close()

		args := []string{
			"-layout",  // Maintain (as best as possible) the original physical layout of the text.
			"-nopgbrk", // Don't insert page breaks (form feed characters) between pages.
			"-",        // Read from stdin
			"-",        // Send the output to stdout.
		}

		cmd := exec.CommandContext(context.Background(), "pdftotext", args...)

		// Set the input to the command as the uploaded file
		cmd.Stdin = file

		var buf bytes.Buffer
		cmd.Stdout = &buf
		cmd.Stderr = &buf // Capture any errors

		if err := cmd.Run(); err != nil {
			log.Printf("pdftotext error: %v, output: %s", err, buf.String())
			http.Error(w, fmt.Sprintf("Failed to execute pdftotext: %v", err), http.StatusInternalServerError)
			return
		}

		text := string(buf.String())
		OrderId := common.ExtractOrderID(text)
		ContentOrder := common.ExtractContentOrder(text)

		if OrderId == nil || ContentOrder == nil {
			continue
		}
		responseText = append(responseText, map[string]string{
			"OrderId":      *OrderId,
			"ContentOrder": *ContentOrder,
		})
		// response = append(response, text)
	}
	responseBody, err := json.Marshal(responseText)
	if err != nil {
		log.Printf("Failed to marshal response: %v", err)
		http.Error(w, "Failed to marshal response", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	if _, err := w.Write(responseBody); err != nil {
		log.Printf("Failed to write response: %v", err)
	}
}

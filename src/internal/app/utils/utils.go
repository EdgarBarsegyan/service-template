package utils

import (
	"fmt"
	"net/http"
	"strings"
)

func GetHeaders(header http.Header) string {
	rowHeaders := make([]string, 0, len(header))

	for key, values := range header {
		row := fmt.Sprintf("%s %s", key, values)
		rowHeaders = append(rowHeaders, row)
	}

	return strings.Join(rowHeaders, ";\n")
}

type CustomResponseWriter struct {
	http.ResponseWriter
	StatusCode int
}

func (lrw *CustomResponseWriter) WriteHeader(code int) {
	lrw.StatusCode = code
	lrw.ResponseWriter.WriteHeader(code)
}

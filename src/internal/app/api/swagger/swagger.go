package swagger

import (
	"embed"
	"io/fs"
	"net/http"
	"service-template/pkg/api"
)

//go:embed static/*
var swaggerUIFiles embed.FS

func serveOpenAPISpec(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", "application/json")
	d,_ := api.GetSwagger()
	spec, _ := d.MarshalJSON()
    w.Write(spec)
}

func getSwaggerUIFS() fs.FS {
    subFS, err := fs.Sub(swaggerUIFiles, "static")
    if err != nil {
        panic(err)
    }
    return subFS
}

func RegisterHandlers(mux *http.ServeMux) {
	mux.Handle(
		"GET /swagger/", 
		http.StripPrefix(
			"/swagger",
			http.FileServer(http.FS(getSwaggerUIFS())),
		),
	)

	mux.Handle(
		"/openapi.json", 
		http.HandlerFunc(serveOpenAPISpec),
	)
}
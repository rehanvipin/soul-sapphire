package urlshort

import (
	"net/http"

	"gopkg.in/yaml.v2"
)

type mapper struct {
	switches map[string]string
	fallback http.Handler
}

func (m *mapper) redirect(w http.ResponseWriter, req *http.Request) {
	var og = req.URL.Path
	if _, ok := m.switches[og]; !ok {
		m.fallback.ServeHTTP(w, req)
	} else {
		target := m.switches[og]
		http.Redirect(w, req, target,
			http.StatusTemporaryRedirect)
	}
}

// MapHandler will return an http.HandlerFunc (which also
// implements http.Handler) that will attempt to map any
// paths (keys in the map) to their corresponding URL (values
// that each key in the map points to, in string format).
// If the path is not provided in the map, then the fallback
// http.Handler will be called instead.
func MapHandler(pathsToUrls map[string]string, fallback http.Handler) http.HandlerFunc {

	var movement = mapper{
		switches: pathsToUrls,
		fallback: fallback,
	}

	return http.HandlerFunc(movement.redirect)
}

type wicker struct {
	Path, URL string
}

// YAMLHandler will parse the provided YAML and then return
// an http.HandlerFunc (which also implements http.Handler)
// that will attempt to map any paths to their corresponding
// URL. If the path is not provided in the YAML, then the
// fallback http.Handler will be called instead.
//
// YAML is expected to be in the format:
//
//     - path: /some-path
//       url: https://www.some-url.com/demo
//
// The only errors that can be returned all related to having
// invalid YAML data.
//
// See MapHandler to create a similar http.HandlerFunc via
// a mapping of paths to urls.
func YAMLHandler(yml []byte, fallback http.Handler) (http.HandlerFunc, error) {
	var f = []wicker{}
	err := yaml.Unmarshal(yml, &f)
	if err != nil {
		return nil, err
	}
	var pathToUrls = make(map[string]string)
	for _, mapped := range f {
		pathToUrls[mapped.Path] = mapped.URL
	}

	return MapHandler(pathToUrls, fallback), nil
}

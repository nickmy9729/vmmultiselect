package main

import (
	"fmt"
	"html/template"
	"net/http"
	"os"
	"time"

	"github.com/yourusername/vmmultiselect/config"
	"github.com/yourusername/vmmultiselect/internal/api"
	"github.com/yourusername/vmmultiselect/internal/vmmultiselect"
)

// Define a simple template for rendering the health status.
const healthTemplate = `
<!DOCTYPE html>
<html>
<head>
	<title>VictoriaMetrics Health Status</title>
</head>
<body>
	<h1>VictoriaMetrics Health Status</h1>
	<p>Group: {{ .GroupName }}</p>
	{{ range $endpoint, $status := .HealthMap }}
		<p>{{ $endpoint }}: {{ $status }}</p>
	{{ end }}
</body>
</html>
`

func main() {
	cfg, err := config.LoadConfig("config/config.json")
	if err != nil {
		fmt.Printf("Error loading config: %s\n", err)
		os.Exit(1)
	}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		groupName := r.URL.Query().Get("group")
		if groupName == "" {
			http.Error(w, "Missing 'group' query parameter", http.StatusBadRequest)
			return
		}

		headers := http.Header{
			"X-GROUP": []string{groupName},
		}

		// Specify the timeout for health check requests
		timeout := 2 * time.Second

		// Perform health checks with timeout
		statusMap := make(map[string]string)
		for _, endpoint := range cfg.Groups[groupName] {
			_, err := api.MakeAPIRequestWithTimeout(endpoint+"/health", headers, timeout)
			if err != nil {
				statusMap[endpoint] = "Unhealthy"
			} else {
				statusMap[endpoint] = "Healthy"
			}
		}

		// Render the template with the health status data
		tmpl := template.Must(template.New("health").Parse(healthTemplate))
		tmplData := struct {
			GroupName string
			HealthMap map[string]string
		}{
			GroupName: groupName,
			HealthMap: statusMap,
		}
		err := tmpl.Execute(w, tmplData)
		if err != nil {
			http.Error(w, fmt.Sprintf("Error rendering template: %s", err), http.StatusInternalServerError)
		}
	})

	fmt.Println("Listening on :8080...")
	err = http.ListenAndServe(":8080", nil)
	if err != nil {
		fmt.Printf("Error starting server: %s\n", err)
	}
}


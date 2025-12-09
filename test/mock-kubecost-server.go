package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Mock Kubecost Server\nAvailable endpoints:\n  - /model/savings\n")
	})

	http.HandleFunc("/model/savings", func(w http.ResponseWriter, r *http.Request) {
		log.Printf("Request: %s %s from %s", r.Method, r.URL.Path, r.RemoteAddr)

		response := map[string]interface{}{
			"recommendedRightSizing": []map[string]interface{}{
				{
					"clusterId":      "prod-aks",
					"namespace":      "production",
					"controllerKind": "deployment",
					"controllerName": "payment-api",
					"container":      "payment-api",
					"recommendations": map[string]interface{}{
						"cpu": map[string]interface{}{
							"current":        1000.0,
							"recommended":    500.0,
							"monthlySavings": 200.0,
						},
						"memory": map[string]interface{}{
							"current":        2048.0,
							"recommended":    1024.0,
							"monthlySavings": 300.0,
						},
					},
					"totalMonthlySavings": 500.0,
				},
				{
					"clusterId":      "prod-aks",
					"namespace":      "monitoring",
					"controllerKind": "deployment",
					"controllerName": "logging-service",
					"container":      "fluentd",
					"recommendations": map[string]interface{}{
						"cpu": map[string]interface{}{
							"current":        500.0,
							"recommended":    200.0,
							"monthlySavings": 150.0,
						},
						"memory": map[string]interface{}{
							"current":        1024.0,
							"recommended":    512.0,
							"monthlySavings": 250.0,
						},
					},
					"totalMonthlySavings": 400.0,
				},
			},
			"underutilizedNodes": []map[string]interface{}{
				{
					"node":                    "aks-nodepool1-12345678-0",
					"monthlySavings":          1200.0,
					"providerId":              "azure:///node/id",
					"recommendedInstanceType": "Standard_D4s_v3",
				},
			},
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
		log.Printf("Responded with %d recommendations",
			len(response["recommendedRightSizing"].([]map[string]interface{})))
	})

	port := ":8080"
	fmt.Println("===================================================================")
	fmt.Println("  Mock Kubecost Server Running")
	fmt.Println("===================================================================")
	fmt.Println()
	fmt.Printf("Listening on: http://localhost%s\n", port)
	fmt.Println()
	fmt.Println("Endpoints:")
	fmt.Printf("  - http://localhost%s/\n", port)
	fmt.Printf("  - http://localhost%s/model/savings\n", port)
	fmt.Println()
	fmt.Println("Test with:")
	fmt.Printf("  curl http://localhost%s/model/savings\n", port)
	fmt.Printf("  ./kosva check --kubecost-url http://localhost%s\n", port)
	fmt.Println()
	fmt.Println("Press Ctrl+C to stop")
	fmt.Println("-------------------------------------------------------------------")

	log.Fatal(http.ListenAndServe(port, nil))
}

package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type Node struct {
	Type     string `json:"type"`
	Label    string `json:"label"`
	Children []Node `json:"children,omitempty"`
}

func main() {
	http.HandleFunc("/generatexdomea", Generatexdomea)
	fmt.Println("Server running on :8080")
	http.ListenAndServe(":8080", nil)
}

// Converts a Node tree to XML string
func nodeToXML(n Node) string {
	xml := fmt.Sprintf("<%s>%s", n.Type, n.Label)
	for _, child := range n.Children {
		xml += nodeToXML(child)
	}
	xml += fmt.Sprintf("</%s>", n.Type)
	return xml
}

func Generatexdomea(w http.ResponseWriter, r *http.Request) {

    w.Header().Set("Access-Control-Allow-Origin", "*")
    w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
    w.Header().Set("Access-Control-Allow-Headers", "Content-Type")


    if r.Method == http.MethodOptions {
        w.WriteHeader(http.StatusOK)
        return
    }

	if r.Method != http.MethodPost {
		http.Error(w, "Only POST allowed", http.StatusMethodNotAllowed)
		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Failed to read body", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	var nodes []Node
	if err := json.Unmarshal(body, &nodes); err != nil {
		http.Error(w, "Failed to parse JSON: "+err.Error(), http.StatusBadRequest)
		return
	}

	xml := ""
	for _, node := range nodes {
		xml += nodeToXML(node)
	}

	w.Header().Set("Content-Type", "application/xml")
	w.Write([]byte(xml))
}

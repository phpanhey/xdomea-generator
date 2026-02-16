package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"strings"
)

type Node struct {
	Type     string `json:"type"`
	Label    string `json:"label"`
	Children []Node `json:"children,omitempty"`
}

func main() {
	http.HandleFunc("/generatexdomea", Generatexdomea)
	fmt.Println("Server running on :8080")
	fmt.Println("/generatexdomea endpoint is exposed")
	http.ListenAndServe(":8080", nil)

	/*
	   		filesMap := readFiles()
	   		jsondata := `[
	   	   {
	   	        "type": "akte",
	   	        "label": "Akte 1",
	   	        "children": [
	   	            {
	   	                "type": "vorgang",
	   	                "label": "Vorgang 1",
	   	                "children": [
	   	                    {
	   	                        "type": "untervorgang",
	   	                        "label": "Untervorgang 1",
	   	                        "children": [
	   	                            {
	   	                                "type": "dokument",
	   	                                "label": "Dokument 1"
	   	                            }
	   	                        ]
	   	                    }
	   	                ]
	   	            }
	   	        ]
	   	    }

	   ]`

	   	var nodes []Node
	   	if err := json.Unmarshal([]byte(jsondata), &nodes); err != nil {
	   		panic("json could not be parsed")
	   	}

	   	xml := filesMap["base_prefix"]

	   	for _, node := range nodes {
	   		xml += nodeToXML(node, filesMap)
	   	}

	   	xml += filesMap["base_postfix"]

	   	exportFile := "output.xml"
	   	stringToFile(exportFile, xml)
	   	format(exportFile)
	*/
}

func nodeToXML(n Node, filesMap map[string]string) string {

	nodeType := n.Type

	prefixContent := filesMap[nodeType+"_prefix"]
	postfixContent := filesMap[nodeType+"_postfix"]

	xml := prefixContent

	for _, child := range n.Children {
		xml += nodeToXML(child, filesMap)
	}

	xml += postfixContent

	return xml
}

func Generatexdomea(w http.ResponseWriter, r *http.Request) {

	filesMap := readFiles()

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
	if err := json.Unmarshal([]byte(body), &nodes); err != nil {
		panic("json could not be parsed")
	}

	xml := filesMap["base_prefix"]

	for _, node := range nodes {
		xml += nodeToXML(node, filesMap)
	}

	xml += filesMap["base_postfix"]

	exportFile := "output.xml"
	stringToFile(exportFile, xml)
	format(exportFile)

	w.Header().Set("Content-Type", "application/xml")
	w.Write([]byte(readFile("output.xml")))
}

func readFiles() map[string]string {
	filesMap := make(map[string]string)
	entries, err := os.ReadDir("./template")
	if err != nil {
		panic(err)
	}
	for _, e := range entries {
		path := "template/"
		fileName := e.Name()

		fileContent := readFile(path + fileName)

		filesMap[getKeyFromFileName(fileName)] = fileContent
	}
	return filesMap
}

func getKeyFromFileName(fileName string) string {
	return strings.Split(fileName, ".")[0]
}

func stringToFile(path string, content string) error {
	return os.WriteFile(path, []byte(content), 0644)
}

func format(fileName string) {
	cmd := exec.Command("xmllint", "--format", fileName, "-o", fileName)
	err := cmd.Run()
	if err != nil {
		panic(err)
	}
}

func readFile(path string) string {
	fileContent, error := os.ReadFile(path)
	if error != nil {
		panic("file not found")
	}
	return string(fileContent)
}

package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
)

const openAIEndpoint = "https://api.openai.com/v1/chat/completions"
const openAIModel = "gpt-4"

func main() {
	var filePath string
	var apiKey string
	var outputPath string
	var explainMode bool

	var rootCmd = &cobra.Command{
		Use:   "kubesage",
		Short: "DevOps CLI assistant powered by LLM",
	}

	var auditCmd = &cobra.Command{
		Use:   "audit",
		Short: "Audit a Kubernetes YAML file using a LLM",
		Run: func(cmd *cobra.Command, args []string) {
			if filePath == "" || apiKey == "" {
				log.Fatal("--file and --api-key are required")
			}
			runAudit(filePath, apiKey, outputPath, explainMode)
		},
	}

	auditCmd.Flags().StringVarP(&filePath, "file", "f", "", "Path to the YAML file")
	auditCmd.Flags().StringVarP(&apiKey, "api-key", "k", "", "OpenAI API key")
	auditCmd.Flags().StringVarP(&outputPath, "output", "o", "", "Optional output file to write improved YAML")
	auditCmd.Flags().BoolVarP(&explainMode, "explain", "e", false, "Explain YAML contents instead of auditing")

	var generateCmd = &cobra.Command{
		Use:   "generate",
		Short: "Generate a new Kubernetes YAML from a natural language description",
		Run: func(cmd *cobra.Command, args []string) {
			if apiKey == "" || len(args) == 0 {
				log.Fatal("--api-key and prompt are required")
			}
			runGenerate(apiKey, strings.Join(args, " "))
		},
	}
	generateCmd.Flags().StringVarP(&apiKey, "api-key", "k", "", "OpenAI API key")

	rootCmd.AddCommand(auditCmd)
	rootCmd.AddCommand(generateCmd)
	rootCmd.Execute()
}

func runAudit(filePath, apiKey, outputPath string, explainMode bool) {
	content, err := ioutil.ReadFile(filePath)
	if err != nil {
		log.Fatalf("Error reading file: %v", err)
	}

	var prompt string
	if explainMode {
		prompt = fmt.Sprintf("You are a Kubernetes expert. Please explain the following YAML file line by line:\n\n---\n%s", content)
	} else {
		prompt = fmt.Sprintf("You are a Kubernetes expert. Here is a YAML file:\n\n---\n%s\n\nAnalyze it, suggest improvements, and provide explanations.", content)
	}

	reqBody := map[string]interface{}{
		"model": openAIModel,
		"messages": []map[string]string{
			{"role": "system", "content": "You are a senior DevOps engineer."},
			{"role": "user", "content": prompt},
		},
	}

	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		log.Fatalf("Error encoding JSON: %v", err)
	}

	req, err := http.NewRequest("POST", openAIEndpoint, bytes.NewBuffer(jsonData))
	if err != nil {
		log.Fatalf("Error creating HTTP request: %v", err)
	}
	req.Header.Set("Authorization", "Bearer "+apiKey)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatalf("Error calling API: %v", err)
	}
	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)
	var res map[string]interface{}
	if err := json.Unmarshal(body, &res); err != nil {
		log.Fatalf("Error decoding response: %v", err)
	}

	fmt.Println("\n💡 LLM Response:")
	if choices, ok := res["choices"].([]interface{}); ok && len(choices) > 0 {
		msg := choices[0].(map[string]interface{})["message"].(map[string]interface{})["content"].(string)
		fmt.Println(msg)
		if outputPath != "" && !explainMode {
			outputFile := outputPath
			if filepath.Ext(outputFile) == "" {
				outputFile += ".yaml"
			}
			ioutil.WriteFile(outputFile, []byte(msg), 0644)
			fmt.Printf("\n✅ Improved YAML written to %s\n", outputFile)
		}
	} else {
		fmt.Println("No valid response from API.")
	}
}

func runGenerate(apiKey, prompt string) {
	fullPrompt := fmt.Sprintf("You are a Kubernetes expert. Generate a complete and valid Kubernetes YAML manifest based on this description:\n\n%s", prompt)

	reqBody := map[string]interface{}{
		"model": openAIModel,
		"messages": []map[string]string{
			{"role": "system", "content": "You are a senior DevOps engineer."},
			{"role": "user", "content": fullPrompt},
		},
	}

	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		log.Fatalf("Error encoding JSON: %v", err)
	}

	req, err := http.NewRequest("POST", openAIEndpoint, bytes.NewBuffer(jsonData))
	if err != nil {
		log.Fatalf("Error creating HTTP request: %v", err)
	}
	req.Header.Set("Authorization", "Bearer "+apiKey)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatalf("Error calling API: %v", err)
	}
	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)
	var res map[string]interface{}
	if err := json.Unmarshal(body, &res); err != nil {
		log.Fatalf("Error decoding response: %v", err)
	}

	fmt.Println("\n📦 Generated YAML:")
	if choices, ok := res["choices"].([]interface{}); ok && len(choices) > 0 {
		msg := choices[0].(map[string]interface{})["message"].(map[string]interface{})["content"].(string)
		fmt.Println(msg)
	} else {
		fmt.Println("No valid response from API.")
	}
}

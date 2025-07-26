package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
)

// sample YAML içeriği
const sampleYAML = `tests:
  - name: Get Users
    method: GET
    url: https://jsonplaceholder.typicode.com/users
    headers:
      Content-Type: application/json
    expect:
      status: 200

  - name: Get Single User
    method: GET
    url: https://jsonplaceholder.typicode.com/users/1
    headers:
      Content-Type: application/json
    expect:
      status: 200
      body: | 
        {
          "id": 1,
          "name": "Leanne Graham",
          "username": "Bret",
          "email": "Sincere@april.biz",
          "address": {
            "street": "Kulas Light",
            "suite": "Apt. 556",
            "city": "Gwenborough",
            "zipcode": "92998-3874",
            "geo": {
              "lat": "-37.3159",
              "lng": "81.1496"
            }
          },
          "phone": "1-770-736-8031 x56442",
          "website": "hildegard.org",
          "company": {
            "name": "Romaguera-Crona",
            "catchPhrase": "Multi-layered client-server neural-net",
            "bs": "harness real-time e-markets"
          }
        }
`

var outputPath string

var generateCmd = &cobra.Command{
	Use:   "generate",
	Short: "Generate a sample YAML test suite file",
	Long:  "Creates an example YAML file that defines HTTP test cases for Tapir.",
	Run: func(cmd *cobra.Command, args []string) {
		path := outputPath
		if path == "" {
			path = "sample.yaml"
		}
		abs, _ := filepath.Abs(path)

		err := os.WriteFile(abs, []byte(sampleYAML), 0644)
		if err != nil {
			fmt.Printf("❌ Failed to create sample file: %v\n", err)
			return
		}
		fmt.Printf("✅ Sample test suite written to %s\n", abs)
	},
}

func init() {
	generateCmd.Flags().StringVarP(&outputPath, "output", "o", "", "Output file path (default: sample.yaml)")
	rootCmd.AddCommand(generateCmd)
}

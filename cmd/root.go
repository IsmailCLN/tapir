package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "tapir",
	Short: "YAML tabanlı HTTP API test runner",
	Long: `Tapir, HTTP tabanlı REST API'leri YAML dosyaları aracılığıyla test etmenizi sağlayan bir CLI aracıdır.

YAML formatında tanımladığınız test senaryolarını sırayla çalıştırır, 
her isteğin yanıt süresini, boyutunu ve beklenen ile alınan HTTP durum kodlarını karşılaştırır.
İsteğe bağlı olarak response body içeriğini de doğrulayabilir.

Test sonuçlarını terminalde renkli ve tablolu olarak sunar, dilerseniz dosyaya veya panoya da aktarabilirsiniz.

Kullanım örneği:
  tapir run test-data/sample.yaml
`,
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

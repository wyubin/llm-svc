// vcfgo commanline entry
package main

import (
	"fmt"
	"os"

	"github.com/wyubin/llm-svc/cmd/rag/cmd"
)

func main() {
	err := cmd.Execute()
	if err != nil && err.Error() != "" {
		fmt.Println(err)
		os.Exit(99)
	}
}

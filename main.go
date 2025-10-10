package main

import (
	"fmt"
	"log"

	"github.com/ntlam20/daily-quotes/pkg/utils"
)

func main() {
	fmt.Println("🚀 Fetching daily quote...")

	q, err := utils.FetchQuote()
	if err != nil {
		log.Fatalf("failed to fetch quote: %v", err)
	}

	if err := utils.UpdateREADME(q); err != nil {
		log.Fatalf("failed to update README: %v", err)
	}

	fmt.Println("✅ README updated successfully!")
}

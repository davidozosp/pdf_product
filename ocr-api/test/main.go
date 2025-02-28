package main

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"os/exec"
	"regexp"
	"strings"
)

func main() {

	// count file in folder input
	files, err := os.ReadDir("input")
	if err != nil {
		panic(err)
	}

	for _, file := range files {
		fmt.Println(file.Name())

		args := []string{
			"-layout",              // Maintain (as best as possible) the original physical layout of the text.
			"-nopgbrk",             // Don't insert page breaks (form feed characters) between pages.
			"input/" + file.Name(), // The input file.
			"-",                    // Send the output to stdout.
		}
		cmd := exec.CommandContext(context.Background(), "pdftotext", args...)
		var buf bytes.Buffer
		cmd.Stdout = &buf

		if err := cmd.Run(); err != nil {
			fmt.Println(err)
			return
		}
		text := string(buf.String())
		// Extract Order ID
		orderIDRegex := regexp.MustCompile(`Mã đơn hàng: (\w+)`)
		orderIDMatches := orderIDRegex.FindStringSubmatch(text)

		var orderID string
		if len(orderIDMatches) > 1 {
			orderID = orderIDMatches[1]
		}

		// Extract Product Details
		productRegex := regexp.MustCompile(`(\d+\.\s[^\n]+)`)
		productMatches := productRegex.FindAllString(text, -1)

		fmt.Println("Order ID:", orderID)
		fmt.Println("Product Details:")
		for _, product := range productMatches {
			fmt.Println(strings.TrimSpace(product))
		}
	}

}

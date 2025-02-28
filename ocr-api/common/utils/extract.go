package common

import "regexp"

func ExtractOrderID(text string) *string {
	orderIDRegex := regexp.MustCompile(`Mã đơn hàng: (\w+)`)
	orderIDMatches := orderIDRegex.FindStringSubmatch(text)
	var orderID string
	if len(orderIDMatches) > 1 {
		orderID = orderIDMatches[1]
	}
	return &orderID
}

func ExtractContentOrder(text string) *string {
	productRegex := regexp.MustCompile(`(\d+\.\s[^\n]+)`)
	productMatches := productRegex.FindAllString(text, -1)
	var content string
	for _, product := range productMatches {
		content += product + "\n"
	}
	return &content
}

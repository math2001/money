package main

import (
	"bufio"
	"encoding/base64"
	"fmt"
	"log"
	"os"
	"strings"
)

func ParseSession() {
	reader := bufio.NewReader(os.Stdin)
	line, err := reader.ReadString('\n')
	if err != nil {
		log.Fatal(err)
	}
	parts := strings.Split(line, ".")
	if len(parts) != 3 {
		log.Fatalf("Should have 3 parts, %d in %q", len(parts), line)
	}
	content, err := base64.StdEncoding.DecodeString(parts[1])
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(string(content))
}

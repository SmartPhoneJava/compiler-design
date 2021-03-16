package repository

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

// GetInt Получить число из консоли
func GetInt(v *int, label string) {
	var stop = false

	fmt.Println(label)
	for !stop {
		_, err := fmt.Scanf("%d\n", v)
		if v != nil {
			stop = err == nil
		}
	}
}

// Получить текст из консоли
func GetString(str *string, label string) {
	var (
		reader = bufio.NewReader(os.Stdin)
		stop   bool
	)
	fmt.Println(label)
	for !stop {
		v, err := reader.ReadString('\n')
		v = strings.Trim(v, "\n")
		if err == nil && v != "" {
			*str = v
			break
		}
	}
}

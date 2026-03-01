package application

import (
	"fmt"
	"strings"
)

func NoInteractiveGame(args []string) {
	quest := []rune(strings.ToLower(args[1]))
	ans := strings.ToLower(args[2])
	result := make([]rune, len(quest))

	ansSet := map[rune]struct{}{}
	for _, c := range ans {
		ansSet[c] = struct{}{}
	}

	for i, c := range quest {
		if _, ok := ansSet[c]; ok {
			result[i] = c
		}
	}

	res := ";POS"
	for _, r := range result {
		if r == 0 {
			fmt.Print("*")
			res = ";NEG"
		} else {
			fmt.Print(string(r))
		}
	}
	fmt.Println(res)
}

package input

import (
	"bufio"
	"fmt"
	"os"
)

/*
Scan ...
*/
func Scan(msg string) string {
	scanner := bufio.NewScanner(os.Stdin)
	fmt.Print(msg)
	scanner.Scan()
	return scanner.Text()
}

package main
import (
	"fmt"
)

func main() {
	fmt.Print("Write expression (e.g. 10 + 5): ")
	var firstNum float64
	var operator string
	var secondNum float64
	fmt.Scanln(&firstNum, &operator, &secondNum)
	var result float64
	switch operator {
	case "+":
		result = firstNum + secondNum
	case "-":
		result = firstNum - secondNum
	case "*":
		result = firstNum * secondNum
	case "/":
		if secondNum == 0 {
			fmt.Println("Error: Division by zero is not allowed.")
			return
		}
		result = firstNum / secondNum
	default:
		fmt.Println("Error: Invalid operator.")
		return
	}
	fmt.Printf("%.2f %s %.2f = %.2f\n", firstNum, operator, secondNum, result)
}

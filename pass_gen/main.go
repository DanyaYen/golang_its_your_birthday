package main 

import (
	"fmt"
	"flag"
	"crypto/rand"
	"math/big"
)

func main() {
	const lowercase = "abcdefghijklmnopqrstuvwxyz"
	const uppercase = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	const digits = "0123456789"
	const symbols = "!@#$%^&*(){}[],./;:?_-+=|`~"
	length := flag.Int("len", 12, "Length of password")

	includeDigits := flag.Bool("digits", true, "Include digits")
	includeLowercase := flag.Bool("lower", true, "Include lowercase letters")
	includeSymbols := flag.Bool("symbols", true, "Include symbols")
	includeUppercase := flag.Bool("upper", true, "Include uppercase letters")
	flag.Parse()
	
	var characterPool []rune
	
	if *includeLowercase {
		characterPool = append(characterPool, []rune(lowercase)...)
	}
	if *includeUppercase {
		characterPool = append(characterPool, []rune(uppercase)...)
	}
	if *includeDigits {
		characterPool = append(characterPool, []rune(digits)...)
	}
	if *includeSymbols {
		characterPool = append(characterPool, []rune(symbols)...)
	}
	if len(characterPool) == 0 {
	    fmt.Println("Error not my type for password")
	    return 
	}

	
	poolLength := len(characterPool) 
	var passwordRunes []rune
	
	for i := 0; i < *length; i++ {
		randomIndex, err := rand.Int(rand.Reader, big.NewInt(int64(poolLength)))
		if err != nil {
			fmt.Println("Error in gen process:", err)
			return 
		}
		index := int(randomIndex.Int64())
    		randomRune := characterPool[index]
    		passwordRunes = append(passwordRunes, randomRune)		    
	}

	generatedPassword := string(passwordRunes)
	fmt.Println(generatedPassword)
}
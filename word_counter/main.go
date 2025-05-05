package main 
import ( 
	"fmt"
	"bufio"
	"io"
	"os"
	"strings"
)

func main() {
	var input io.Reader
	var sourceName string = "stdin"
	
	if len(os.Args) > 1 {
	
		filename := os.Args[1]
		file , err := os.Open(filename)
		
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error while opening file %s: %v\n", filename, err)
			os.Exit(1)
		}
		defer file.Close()
		input = file
		sourceName = filename
	} else {
		input = os.Stdin
	}
	
	scanner := bufio.NewScanner(input)
	var lineCount, wordCount, byteCount int 
	
	for scanner.Scan() {
		line := scanner.Text()
		lineCount++
		words := strings.Fields(line)
		wordCount += len(words)
		byteCount += len(line) + 1
	}
	if err := scanner.Err(); err != nil {
		fmt.Fprintf(os.Stderr, "Error while opening file %s: %v\n", sourceName, err)
		os.Exit(1)
	}
	fmt.Printf("%d %d %d %s\n", lineCount, wordCount, byteCount, sourceName)
}
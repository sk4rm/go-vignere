package main

import (
	"errors"
	"fmt"
	"os"
	"strings"
)

const YEL = "\033[0;33m"
const WHT = "\033[0;37m"
const DESCRIPTION = `Go-vigenere is a encryption and decryption tool based on the Vigenere cipher written in Go.

Usage:
	go <command> [arguments]

The commands are:

	encrypt    encipher text from a file with a specified key
	decrypt    decipher text from a file with a specified key`

func check(err error) {
	if err != nil {
		// panic(err)
		fmt.Println(err)
		os.Exit(0)
	}
}

type Vigenere struct {
	table []string
}

// Generates a tabula recta for vigenere substitution.
//   - `valid_chars` is the set of runes that makes up the plaintext and key/
func (v *Vigenere) generate(valid_chars string) error {
	v.table = make([]string, len(valid_chars))

	for i := range valid_chars {
		v.table[i] = str_rotate_left(valid_chars, i)
	}

	return nil
}

func (v Vigenere) available() bool {
	return len(v.table) != 0
}

func str_rotate_left(str string, n int) string {
	return str[n:] + str[:n]
}

func (v Vigenere) substitute(char rune, keychar rune) (rune, error) {
	if !v.available() {
		return -1, errors.New("no vigenere table generated")
	}

	row := strings.IndexRune(v.table[0], char)
	if row < 0 {
		return -1, fmt.Errorf("substitute: character '%c' not found in table rows", char)
	}

	col := strings.IndexRune(v.table[0], keychar)
	if col < 0 {
		return -1, fmt.Errorf("substitute: character '%c' not found in table columns", keychar)
	}

	// Convert string to rune.
	substituted := []rune(v.table[row])[col]

	return substituted, nil
}

func (v Vigenere) reverse_substitute(char rune, keychar rune) (rune, error) {
	if !v.available() {
		return -1, errors.New("no vigenere table generated")
	}

	row := strings.IndexRune(v.table[0], keychar)
	if row < 0 {
		return -1, fmt.Errorf("substitute: character '%c' not found in table rows", keychar)
	}

	col := strings.IndexRune(v.table[row], char)
	if col < 0 {
		return -1, fmt.Errorf("substitute: character '%c' not found in table columns", char)
	}

	// Convert string to rune.
	substituted := []rune(v.table[0])[col]

	return substituted, nil
}

func (v Vigenere) encrypt(plaintext string, keystring string) (string, error) {
	ciphertext := ""

	if len(plaintext) > len(keystring) {
		original := keystring
		for len(keystring) < len(plaintext) {
			keystring += original
		}
	}

	for i := range plaintext {
		tmp, err := v.substitute(rune(plaintext[i]), rune(keystring[i]))
		check(err)

		ciphertext += string(tmp)

	}
	return ciphertext, nil
}

func (v Vigenere) decrypt(ciphertext string, keystring string) (string, error) {
	plaintext := ""

	if len(ciphertext) > len(keystring) {
		original := keystring
		for len(keystring) < len(ciphertext) {
			keystring += original
		}
	}

	for i := range ciphertext {
		tmp, err := v.reverse_substitute(rune(ciphertext[i]), rune(keystring[i]))
		check(err)

		plaintext += string(tmp)

	}
	return plaintext, nil
}

func main() {
	// Read command line arguments.
	args := os.Args[1:]

	if len(args) == 0 {
		fmt.Println(DESCRIPTION)
		fmt.Println()
		os.Exit(0)
	}

	command := args[0]

	if command == "help" {
		args = args[1:]

		if len(args) == 0 {
			fmt.Println(DESCRIPTION)
			fmt.Println()

		} else if len(args) == 1 {
			// topic := args[0]
			panic("help topic not yet implemented")

		} else {
			fmt.Println("Usage:\n\n    go-vigenere.exe help <command>")
			fmt.Println()
			os.Exit(0)
		}

	} else if command == "encrypt" {
		args = args[1:]
		if len(args) != 2 {
			fmt.Printf("Expected 2 arguments but got %v instead.\n\n", len(args))
			fmt.Println("Run 'go-vigenere.exe help encrypt' to learn more.")
			os.Exit(0)
		}

		key, plaintext_file := args[0], args[1] // TODO path strings

		// Read plaintext from file.
		data, err := os.ReadFile(plaintext_file)
		check(err)
		plaintext := string(data)

		var vigenere Vigenere
		err = vigenere.generate("0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz,. ")
		check(err)

		ciphertext, err := vigenere.encrypt(plaintext, key)
		check(err)

		fmt.Println(ciphertext)

	} else if command == "decrypt" {
		args = args[1:]
		if len(args) != 2 {
			fmt.Printf("Expected 2 arguments but got %v instead.\n\n", len(args))
			fmt.Println("Run 'go-vigenere.exe help decrypt' to learn more.")
			os.Exit(0)
		}

		key, ciphertext_file := args[0], args[1] // TODO path strings

		// Read plaintext from file.
		data, err := os.ReadFile(ciphertext_file)
		check(err)
		ciphertext := string(data)

		var vigenere Vigenere
		err = vigenere.generate("0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz,. ")
		check(err)

		plaintext, err := vigenere.decrypt(ciphertext, key)
		check(err)

		fmt.Println(plaintext)

	} else {
		fmt.Printf("Unknown command provided: %v\n", command)
		fmt.Println("Run 'go-vigenere.exe help' for usage.")
	}
}

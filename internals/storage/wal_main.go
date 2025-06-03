package storage

import "fmt"

func writeToWal() error {
	isSpaceLeft, err := checkingSpace()
	if err != nil {
		fmt.Printf("failed to check space %v",err)
	}

	if isSpaceLeft{
		
	}

	return ,nil

}

func checkingSpace(): bool {
	
}
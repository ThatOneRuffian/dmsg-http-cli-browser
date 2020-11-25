package dmsggui

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

func deleteServerWizard() {
DeletePrompt:
	fmt.Print("Which server do you want to delete? (Enter C to Cancel): ")
	consoleInputWhichServer := bufio.NewReader(os.Stdin)
	userDelete, _ := consoleInputWhichServer.ReadString('\n')
	userDelete = strings.ToUpper(stripIllegalChars(userDelete))

	switch userDelete {

	case "C":
		goto ExitLoop

	default:
		userInt, err := strconv.Atoi(userDelete)
		if err != nil {
			break
		}

	ConfirmDelete:
		if userInt >= 1 && userInt <= len(SavedServers) {
			deleteConfirmPrompt := fmt.Sprintf("Are you sure you want to delete (Y/N)? { %s | %s }", SavedServers[userInt-1][0], SavedServers[userInt-1][1])
			fmt.Println(deleteConfirmPrompt)
			deleteConfirmInput := bufio.NewReader(os.Stdin)

			deleteConfirm, _ := deleteConfirmInput.ReadString('\n')
			deleteConfirm = strings.ToUpper(stripIllegalChars(deleteConfirm))

			switch deleteConfirm {

			case "Y":

				deleteServerIndex(userInt)
				goto ExitLoop

			case "N":
				goto ExitLoop

			default:
				goto ConfirmDelete
			}

		} else {
			break
		}
	}
	goto DeletePrompt
ExitLoop:
}

func addServer() string {
	invalidChar := ";"
	keyLength := 66
	consoleInput := bufio.NewReader(os.Stdin)
	fmt.Print("Please enter the public key for the dmsg-http server you want to add (C to Cancel): ")

PubKeyInput:
	publicKey, _ := consoleInput.ReadString('\n')
	publicKey = stripIllegalChars(publicKey)
	if strings.ToUpper(publicKey) == "C" {
		publicKey = ""
		goto Exit
	}
	if strings.Contains(publicKey, invalidChar) {
		fmt.Printf("Pulic key cannot contain '%s'.", invalidChar)
		fmt.Print("\nPlease enter public key again (C to Cancel): ")
		goto PubKeyInput
	}
	if len(publicKey) == keyLength {
	FriendlyName:
		fmt.Print("Add a friendly name to this public key (default: [public_key] ): ")
		friendlyName, _ := consoleInput.ReadString('\n')
		friendlyName = stripIllegalChars(friendlyName)
		if len(friendlyName) == 0 {
			appendToConfig(publicKey, publicKey)
		} else if strings.Contains(friendlyName, invalidChar) {
			fmt.Printf("Friendly name cannot contain '%s'\n", invalidChar)
			goto FriendlyName
		}
		appendToConfig(friendlyName, publicKey)
		fmt.Println("Entry cached.")

	} else {
		fmt.Printf("Provided key has length of %d. Expected length of %d.", len(publicKey), keyLength)
		fmt.Print("Invalid key length please enter public key again (C to Cancel): ")
		goto PubKeyInput
	}
Exit:
	return publicKey
}

func FirstRunWizard() {
	fmt.Println("It looks like this is your frist time running the dmsg-http CLI browser.")
	serverPublicKey := addServer()
	browseNow(serverPublicKey)
	LoadCache()
}

func browseNow(serverPublicKey string) {
	consoleInput := bufio.NewReader(os.Stdin)
	if len(serverPublicKey) > 0 {
	Browse:
		fmt.Print("Would you like to browse this server now? (Y/N): ")
		userAnswer, _ := consoleInput.ReadString('\n')

		switch formattedInput := strings.ToUpper(stripIllegalChars(userAnswer)); formattedInput {
		case "Y":
			refreshServerIndex(serverPublicKey, true)
			ServerIndexMenuHandler(serverPublicKey)
		case "N":
			// exit loop and continue...
		default:
			goto Browse
		}
	}
}

// =========== String formatting functions ===========

func stripIllegalChars(userInput string) string {
	InputBytes := []byte(userInput)
	const byteLowRange byte = 33
	const byteHighRange byte = 126
	returnString := ""

	// strip all illegal chars
	for charIndex := range InputBytes {
		byteValue := InputBytes[charIndex]
		if byteValue >= byteLowRange && byteValue <= byteHighRange {
			returnString += string(byteValue)
		}
	}

	return returnString
}

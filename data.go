package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
)

func saveData(content Content, credentials *Credentials) {
	bytes, err := json.Marshal(content)

	if err != nil {
		fmt.Fprintln(os.Stderr, "Unable to marshal json")
		os.Exit(1)
	}

	if credentials != nil && credentials.derivedKey != nil {
		cipherBytes, err := encrypt(credentials.derivedKey, string(bytes))
		bytes = append(credentials.salt, cipherBytes...)

		if err != nil {
			fmt.Fprintln(os.Stderr, "Unable to encrypt data")
			os.Exit(1)
		}
	}

	ioutil.WriteFile(dataFile, bytes, 0644)
}

func loadData(credentials *Credentials) Content {
	var content Content

	jsonData, err := ioutil.ReadFile(dataFile)

	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to read data file : '%s'\n", dataFile)
		os.Exit(1)
	}

	if credentials != nil && credentials.derivedKey != nil {
		stringData, err := decrypt(credentials.derivedKey, string(jsonData[32:]))
		jsonData = []byte(stringData)

		if err != nil {
			fmt.Fprintln(os.Stderr, "Unable to decrypt data")
			os.Exit(1)
		}
	}

	err = json.Unmarshal(jsonData, &content)

	if err != nil {
		fmt.Fprintln(os.Stderr, "Unable to parse JSON data")
		os.Exit(1)
	}

	return content
}

func readSaltFromDataFile() []byte {
	fs, _ := os.Open(dataFile)
	buff := make([]byte, 32)
	fs.Read(buff)

	return buff
}

// Quick and *dirty* way to find out if a file is encrypted. Must find a better way to detect encryption.
func isDataFileEncrypted() bool {
	var content Content

	jsonData, err := ioutil.ReadFile(dataFile)

	if err != nil {
		return false
	}

	err = json.Unmarshal(jsonData, &content)

	if err != nil {
		return true
	}

	return false
}

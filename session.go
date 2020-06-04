package phoenix

import (
	"fmt"
	"log"
	"os"

	"github.com/gorilla/sessions"
)

const (
	keyFileName string = "cookie.key"
	keyLength   int    = 32
)

func createStoreBuilder() Builder {
	store := buildSessionStore()
	return func() sessions.Store {
		return store
	}
}

func buildSessionStore() sessions.Store {
	key, err := getStoredKey()
	if err != nil {
		key = generateRandomBytes(32)
		saveSessionKey(key)
	}
	return sessions.NewCookieStore(key)
}

func getStoredKey() ([]byte, error) {
	keyFile, err := os.Open(fmt.Sprintf("./%s", keyFileName))
	if err != nil {
		log.Println("Failed opening cookie.key file")
		return nil, err
	}
	defer keyFile.Close()
	buffer := make([]byte, keyLength)
	readLen, err := keyFile.Read(buffer)
	if err != nil || readLen != keyLength {
		log.Println("Invalid key in cookie.key")
		return nil, err
	}
	return buffer, nil
}

func saveSessionKey(key []byte) {
	path := fmt.Sprintf("./%s", keyFileName)
	keyFile, err := os.Create(path)
	if err != nil {
		log.Println("Failed creating cookie.key file")
		return
	}
	defer keyFile.Close()
	writeLen, err := keyFile.Write(key)
	if err != nil || writeLen != keyLength {
		log.Println("Error while writing key inside cookie.key. The file will be removed.")
		os.Remove(path)
		return
	}
}

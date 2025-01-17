package password_hasher

import (
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"
)

type PasswordHasher interface {
	HashPassword(password string, salt string, iteration int64) (string, error)
}

var passwordHashers = map[HashAlgorithmType]PasswordHasher{
	HashAlgorithmPBK2DF: &PBK2DFSHA256PasswordHasher{},
}

func getPasswordHasher(algorithm HashAlgorithmType) PasswordHasher {
	passwordHasher, ok := passwordHashers[algorithm]
	if ok {
		return passwordHasher
	}
	return nil
}

func MakeHashedPassword(rawPassword string, algorithm HashAlgorithmType, salt string, iteration int64) (string, error) {
	hasher := getPasswordHasher(algorithm)
	if hasher == nil {
		return "", errors.New("hash algorithm undefined")
	}
	return hasher.HashPassword(rawPassword, salt, iteration)
}

func MakeHashedPasswordStore(rawPassword string) (string, error) {
	salt := os.Getenv("HASH_SALT")
	hash, err := MakeHashedPassword(rawPassword, HashAlgorithmPBK2DF, salt, PasswordHashIteration)
	if err != nil {
		return "", err
	}
	store := fmt.Sprintf("%s$%d$%s$%s", HashAlgorithmPBK2DF, PasswordHashIteration, salt, hash)
	return store, nil
}

func ValidatePassword(password, passwordStore string) (bool, error) {
	val := strings.Split(passwordStore, "$")
	if len(val) != 4 {
		return false, nil
	}

	iterations, err := strconv.ParseInt(val[1], 10, 64)
	if err != nil {
		return false, err
	}

	expectedHash := val[3]

	actualHash, err := MakeHashedPassword(password, val[0], val[2], iterations)

	if err != nil {
		return false, err
	}

	if actualHash == expectedHash {
		return true, nil
	}

	return false, nil
}

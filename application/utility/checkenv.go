package utility

import (
	"fmt"
	"os"
)

func CheckEnv() {
	msgs := []string{}
	db_addr := os.Getenv("DB_ADDR")
	if db_addr == "" {
		msgs = append(msgs, "DB_ADDR env is required, but is not set.")
	}

	token_secret := os.Getenv("TOKEN_SECRET")
	if token_secret == "" {
		msgs = append(msgs, "TOKEN_SECRET env is required but is not set.")
	}

	if len(msgs) == 0 {
		return
	}

	for _, v := range msgs {
		fmt.Println(v)
	}
	panic("Required ENVs are not set.")
}

package main

import (
	"encoding/csv"
	"fmt"
	"log"
	"os"
)

func main() {
	var stdin string
	w := csv.NewWriter(os.Stdout)
	w.Write([]string{
		"host",
		"error",
		"expires_on",
	})
	w.Flush()

	for {
		n, err := fmt.Scan(&stdin)
		if n == 0 {
			break
		}
		if err != nil {
			panic(err)
		}
		r, err := SSLCheck(stdin)
		if err != nil {
			log.Printf("[error]%s --> %v\n", stdin, err)
			w.Write([]string{
				stdin,
				fmt.Sprint(err),
				"",
			})
		} else {
			w.Write([]string{
				fmt.Sprintf("%s:%s", r.Host, r.Port),
				fmt.Sprint(r.Error),
				fmt.Sprint(r.ExpiresOn),
			})
			w.Flush()
		}
	}
}

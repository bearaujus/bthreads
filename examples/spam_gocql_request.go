package examples

import (
	"fmt"
	"log"
	"time"

	"github.com/Bearaujus/bthreads"
	"github.com/gocql/gocql"
)

func SpamGocqlRequest(session *gocql.Session, tableName string) {
	// Initialize bthread instance
	bt, err := bthreads.NewInstance(&bthreads.Config{
		// Add delay between goroutines before it revoked by endless loop
		GoroutinesDelay: 60 * time.Millisecond,
		// Put your defined funcs into X goroutine
		FuncGoroutinesCount: 3,
	})
	if err != nil {
		log.Fatal(err)
		return
	}

	// Add funcs
	bt.AddFunc(func() bool {
		var count int

		if err := session.Query(fmt.Sprintf(`SELECT COUNT(*) FROM %v`, tableName)).Scan(
			&count,
		); err != nil {
			return false
		}
		return true
	})

	// Start bthread instance
	bt.Start()
}

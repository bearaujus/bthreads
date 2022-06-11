package examples

import (
	"log"
	"math/rand"

	"github.com/Bearaujus/bthreads"
)

func TestYourLuck() {
	// Initialize bthread instance
	bt, err := bthreads.NewInstance(&bthreads.Config{
		// Set instance name
		Name: "Probability to get number less than 7 from random number 1 - 1000",
		// Show simple monitor
		HideWorkersData: true,
	})
	if err != nil {
		log.Fatal(err)
		return
	}

	var seed int64
	min, max := 1, 1000

	// Add func
	bt.AddFunc(func() bool {

		rand.Seed(seed)
		randomNumber := rand.Intn(max-min) + min
		seed++
		return randomNumber < 7
	})

	// Start bthread instance
	bt.Start()
}

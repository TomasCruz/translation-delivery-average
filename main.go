package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/TomasCruz/translation-delivery-average/configuration"
	"github.com/TomasCruz/translation-delivery-average/database"
	"github.com/TomasCruz/translation-delivery-average/service"
)

func main() {
	flag.Usage = func() {
		fmt.Fprintln(os.Stderr, os.Args[0], "has mandatory argument 'input_file',")
		fmt.Fprintln(os.Stderr, " and optional argument 'window_size' (10 by default, has to be a positive integer).")
	}

	var (
		windowSize int
		inputFile  string
	)

	flag.IntVar(&windowSize, "window_size", 10, "window size, a positive integer")
	flag.StringVar(&inputFile, "input_file", "", "input file in JSON format")
	flag.Parse()

	if inputFile == "" {
		flag.Usage()
		os.Exit(1)
	}

	if windowSize <= 0 {
		flag.Usage()
		os.Exit(1)
	}

	// In a real world scenario, I would have a perpetually running worker app,
	// storing and processing incoming queue messages, i.e. events. By processing I don't mean the average calculations,
	// and in this case no processing is required as described below.
	//
	// If there is a need to continously monitor translation average times,
	// another job app (or just a goroutine of the worker app) would query the DB every minute,
	// calculate the average for the minute that just expired, and send that info somewhere. Alternativelly to continous monitoring,
	// there would be a translation average job app, for example running overnight to calculate translation averages of minutes of the day.
	// In the described approaches, there would be an event_status field to keep track of the processing of events (CREATED, PROCESSED...),
	// event fields created_at, updated_at would be added, and if there's a need to keep the event history,
	// event_history table would be added to store changes done on the event
	//
	// The concept of this challenge is different, input is fixed and used during the app execution only.
	// So, input doesn't really have to be stored, however I opted for storing it to make it all appear more realistic.
	// For simplicity purposes, events are just stored on receipt, and are only being processed afterwards to calculate averages,
	// so there's no need for processing to deal with storing temporary sums of durations and word counts,
	// concurrency issues that might arise with multiple instances of the app running etc.
	// Approach taken is a sound one, considering the averages are only being calculated
	// for the events that happened in the past (even if "the past" refers to previous minute),
	// i.e. calculating averages is in it's nature about past performance statistics
	//
	// translation_id, as generated on event creation is a bit tricky. There was probably a UUID or ULID generated (for uniqueness purposes),
	// and then it got encoded to 20 chars, likely belonging to an alphabet being a subset of 0-9a-f. Unfortunatelly, I've got attracted to the problem,
	// hence wasting some time on making my own encoding scheme for it...
	// So, I dropped it, decided not to bother with UUID/ULID at all, and just use random 10 bytes encoded to a string.
	// That gives me 256^10 == (2^8)^10 == (2^10)^8 > 1000^8 == 10^24 combinations, so it is fairly random with a low collision probability.
	//
	// Another thing is what is being translated to what - there would likely be an important KPI of translation average of "en" to "fr", ge to "en" etc.
	// That would require indexing of source and destination languages of the events. That could be done with GIN indexing of JSONB payload field,
	// or (as with the approach I've taken) duplicating playload info and indexing it. I've just done the overall average, as required

	// populate configuration
	config := setupFromEnvVars()

	// init DB
	databaseInterface, err := database.InitializeDatabase(config.DbURL)
	if err != nil {
		log.Fatalf("failed to initialize database: %s", err)
	}

	// init service
	/*svc :=*/
	service.NewService(databaseInterface)

	//processInput()
	//averages := calculateAverages()
	// presentAverages()
	fmt.Println("done")
}

func setupFromEnvVars() (config configuration.Config) {
	config.DbURL = readAndCheckEnvVar("DELIVERY_DB_URL")
	return
}

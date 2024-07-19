package server

import (
	"cookie_supply_management/internal/seeds"
	"cookie_supply_management/pkg/base/base_seed"
	"log"

	"github.com/elliotchance/orderedmap"
)

func Runner(sources []string) {
	log.Println("Seeding...")

	router := orderedmap.NewOrderedMap()
	router.Set("user", seeds.NewUserSeed().Seed)

	if len(sources) != 0 {
		for _, source := range sources {
			seed, exists := router.Get(source)
			if !(exists) {
				log.Printf("No source registered with the name: %s", source)
				continue
			}
			runSeed(seed)
		}
	} else {
		for el := router.Front(); el != nil; el = el.Next() {
			runSeed(el.Value)
		}
	}

	log.Println("Finish seeding")
}

func runSeed(fn interface{}) {
	call := fn.(func() (base_seed.Summary, error))
	call()
}

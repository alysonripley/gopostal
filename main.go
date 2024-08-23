package main

import (
	"fmt"

	neardupe "github.com/alyripley/gopostal/neardupe"
)

func main() {

    address1_labels := []string{"house_number", "road", "unit", "city", "state", "postcode"}
    address1_values := []string{"123", "Main St", "#3", "Anytown", "CA", "12345"}

    address2_labels := []string{"house_number", "road", "unit", "city", "state", "postcode"}
    address2_values := []string{"123", "Main Street", "Unit 3", "Anytown", "California", "12345"}

    options := neardupe.NearDupeHashOptions{}
    options.WithName = false
    options.WithAddress = true
    options.WithUnit = true
    options.WithCityOrEquivalent = true
    options.WithSmallContainingBoundaries = false
    options.WithPostalCode = true
    // options.WithLatlon = true
    // options.Latitude = 43.916847
    // options.Longitude = -69.977149
    // options.GeohashPrecision = 6
    options.NameAndAddressKeys = false
    options.NameOnlyKeys = true
    options.AddressOnlyKeys = true

    // neardupehash := neardupe.HashAddress()
    neardupehash1 := neardupe.NearDupe(address1_labels, address1_values, options)
    fmt.Println(neardupehash1)

    neardupehash2 := neardupe.NearDupe(address2_labels, address2_values, options)
    fmt.Println(neardupehash2)

	// Create a map to store the hashes from list1
	hashMap := make(map[string]bool)
	for _, hash := range neardupehash1 {
		hashMap[hash] = true
	}

	// Check for matches in list2
	matchFound := false
	for _, hash := range neardupehash2 {
		if _, found := hashMap[hash]; found {
			matchFound = true
			break
		}
	}
	// Output the result
	if matchFound {
		fmt.Println("Addresses are the same")
	} else {
		fmt.Println("Addresses are unique")
	}
}


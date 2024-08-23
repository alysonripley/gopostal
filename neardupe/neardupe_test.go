package postal

import (
	"reflect"
	"testing"
)

func TestNearDupeHashes(t *testing.T) {
    testCases := []struct {
        name           string
        parsedAddress  map[string]string
        options        NearDupeHashOptions
        expectedHashes []string
    }{
        {
            name: "Full address with unit",
            parsedAddress: map[string]string{
                "house_number": "23",
                "road":         "School St",
                "unit":         "Apt 3",
                "city":         "Brunswick",
                "state":        "ME",
                "postcode":     "04011",
            },
            options: NearDupeHashOptions{
                WithName:             false,
                WithAddress:          true,
                WithUnit:             true,
                WithCityOrEquivalent: true,
                WithPostalCode:       true,
                WithLatlon:           false,
                AddressOnlyKeys:      true,
            },
            expectedHashes: []string{
                "auct|school saint|23|3|brunswick",
                "auct|school street|23|3|brunswick",
                "auct|school|23|3|brunswick",
                "aupc|school saint|23|3|04011",
                "aupc|school street|23|3|04011",
                "aupc|school|23|3|04011",
            },
        },
        {
            name: "Address without unit",
            parsedAddress: map[string]string{
                "house_number": "42",
                "road":         "Main St",
                "city":         "Portland",
                "state":        "OR",
                "postcode":     "97201",
            },
            options: NearDupeHashOptions{
                WithName:             false,
                WithAddress:          true,
                WithUnit:             false,
                WithCityOrEquivalent: true,
                WithPostalCode:       true,
                WithLatlon:           false,
                AddressOnlyKeys:      true,
            },
            expectedHashes: []string{
                "act|main saint|42|portland",
                "act|main street|42|portland",
                "act|main|42|portland",
                "apc|main saint|42|97201",
                "apc|main street|42|97201",
                "apc|main|42|97201",
            },
        },
        {
            name: "Address with name",
            parsedAddress: map[string]string{
                "house":         "Central Park",
                "house_number": "1",
                "road":         "Park Ave",
                "city":         "New York",
                "state":        "NY",
                "postcode":     "10022",
            },
            options: NearDupeHashOptions{
                WithName:             true,
                WithAddress:          true,
                WithUnit:             false,
                WithCityOrEquivalent: true,
                WithPostalCode:       true,
                WithLatlon:           false,
                NameAndAddressKeys:   true,
            },
            expectedHashes: []string{
                "nact|SNTR|park avenue|1|new york", 
                "nact|SNTR|park|1|new york", 
                "nact|NTRL|park avenue|1|new york", 
                "nact|NTRL|park|1|new york", 
                "nact|cent|park avenue|1|new york", 
                "nact|cent|park|1|new york", 
                "nact|entr|park avenue|1|new york", 
                "nact|entr|park|1|new york", 
                "nact|ntra|park avenue|1|new york", 
                "nact|ntra|park|1|new york", 
                "nact|tral|park avenue|1|new york", 
                "nact|tral|park|1|new york", 
                "nact|KP|park avenue|1|new york", 
                "nact|KP|park|1|new york",
                "napc|SNTR|park avenue|1|10022", 
                "napc|SNTR|park|1|10022", 
                "napc|NTRL|park avenue|1|10022", 
                "napc|NTRL|park|1|10022", 
                "napc|cent|park avenue|1|10022", 
                "napc|cent|park|1|10022", 
                "napc|entr|park avenue|1|10022", 
                "napc|entr|park|1|10022", 
                "napc|ntra|park avenue|1|10022", 
                "napc|ntra|park|1|10022", 
                "napc|tral|park avenue|1|10022", 
                "napc|tral|park|1|10022", 
                "napc|KP|park avenue|1|10022", 
                "napc|KP|park|1|10022",
            },
        },
        {
            name: "Address with geohash",
            parsedAddress: map[string]string{
                "house_number": "350",
                "road":         "5th Ave",
                "city":         "New York",
                "state":        "NY",
                "postcode":     "10118",
            },
            options: NearDupeHashOptions{
                WithName:             false,
                WithAddress:          true,
                WithUnit:             false,
                WithCityOrEquivalent: true,
                WithPostalCode:       false,
                WithLatlon:           true,
                Latitude:             40.7484,
                Longitude:            -73.9857,
                GeohashPrecision:     6,
                AddressOnlyKeys:      true,
            },
            expectedHashes: []string{
                "agh|5th avenue|350|dr5ru6",
                "agh|5th avenue|350|dr5ru4", 
                "agh|5th avenue|350|dr5rud", 
                "agh|5th avenue|350|dr5ru3", 
                "agh|5th avenue|350|dr5ru1", 
                "agh|5th avenue|350|dr5ru9", 
                "agh|5th avenue|350|dr5ru7", 
                "agh|5th avenue|350|dr5ru5", 
                "agh|5th avenue|350|dr5rue", 
                "agh|5 avenue|350|dr5ru6", 
                "agh|5 avenue|350|dr5ru4", 
                "agh|5 avenue|350|dr5rud", 
                "agh|5 avenue|350|dr5ru3", 
                "agh|5 avenue|350|dr5ru1", 
                "agh|5 avenue|350|dr5ru9", 
                "agh|5 avenue|350|dr5ru7", 
                "agh|5 avenue|350|dr5ru5", 
                "agh|5 avenue|350|dr5rue", 
                "agh|5th|350|dr5ru6", 
                "agh|5th|350|dr5ru4", 
                "agh|5th|350|dr5rud", 
                "agh|5th|350|dr5ru3", 
                "agh|5th|350|dr5ru1", 
                "agh|5th|350|dr5ru9", 
                "agh|5th|350|dr5ru7", 
                "agh|5th|350|dr5ru5", 
                "agh|5th|350|dr5rue", 
                "agh|5|350|dr5ru6", 
                "agh|5|350|dr5ru4", 
                "agh|5|350|dr5rud", 
                "agh|5|350|dr5ru3", 
                "agh|5|350|dr5ru1", 
                "agh|5|350|dr5ru9", 
                "agh|5|350|dr5ru7", 
                "agh|5|350|dr5ru5", 
                "agh|5|350|dr5rue", 
                "act|5th avenue|350|new york", 
                "act|5 avenue|350|new york", 
                "act|5th|350|new york", 
                "act|5|350|new york",
            },
        },
    }

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			labels := make([]string, 0, len(tc.parsedAddress))
			values := make([]string, 0, len(tc.parsedAddress))
			for key, value := range tc.parsedAddress {
				labels = append(labels, key)
				values = append(values, value)
			}
	
			hashes := NearDupe(labels, values, tc.options)
			if hashes == nil {
				t.Fatalf("NearDupeHashes returned empty: %v", hashes)
			}
	
			if !reflect.DeepEqual(hashes, tc.expectedHashes) {
				t.Errorf("NearDupeHashes returned unexpected hashes for %s:\n", tc.name)
				t.Errorf("Got %d hashes, Want %d hashes\n", len(hashes), len(tc.expectedHashes))
				for i := 0; i < len(hashes) || i < len(tc.expectedHashes); i++ {
					if i < len(hashes) && i < len(tc.expectedHashes) {
						if hashes[i] != tc.expectedHashes[i] {
							t.Errorf("Hash %d:\nGot:  %s\nWant: %s", i, hashes[i], tc.expectedHashes[i])
						}
					} else if i < len(hashes) {
						t.Errorf("Extra hash %d: %s", i, hashes[i])
					} else {
						t.Errorf("Missing expected hash %d: %s", i, tc.expectedHashes[i])
					}
				}
			}
		})
	}
}

func TestPlaceLanguages(t *testing.T) {
	testCases := []struct {
		name             string
		parsedAddress    map[string]string
		expectedLanguages []string
	}{
		{
			name: "English address",
			parsedAddress: map[string]string{
				"house_number": "123",
				"road":         "Main St",
				"city":         "New York",
				"state":        "NY",
				"country":      "USA",
			},
			expectedLanguages: []string{"en"},
		},
		{
			name: "French address",
			parsedAddress: map[string]string{
				"house_number": "15",
				"road":         "Rue de la Paix",
				"city":         "Paris",
				"country":      "France",
			},
			expectedLanguages: []string{"fr"},
		},
		{
			name: "German address",
			parsedAddress: map[string]string{
				"house_number": "1",
				"road":         "Unter den Linden",
				"city":         "Berlin",
				"country":      "Germany",
			},
			expectedLanguages: []string{"de"},
		},
		{
			name: "Multilingual address (Brussels)",
			parsedAddress: map[string]string{
				"house_number": "10",
				"road":         "Rue de la Loi",
				"city":         "Bruxelles",
				"country":      "Belgium",
			},
			expectedLanguages: []string{"fr"},
		},
		{
			name: "Japanese address with Japanese characters",
			parsedAddress: map[string]string{
				"suburb":  "渋谷区",
				"city":    "東京都",
				"country": "日本",
			},
			expectedLanguages: []string{"ja"},
		},
		{
			name: "Japanese address with transliteration",
			parsedAddress: map[string]string{
				"suburb":  "Shibuya",
				"city":    "Tokyo",
				"country": "Japan",
			},
			expectedLanguages: []string{"en"},
		},
		{
			name: "Russian address",
			parsedAddress: map[string]string{
				"house_number": "1",
				"road":         "Тверская улица",
				"city":         "Москва",
				"country":      "Россия",
			},
			expectedLanguages: []string{"ru"},
		},
		{
			name: "Chinese address",
			parsedAddress: map[string]string{
				"road":    "北京市东城区东长安街",
				"city":    "北京市",
				"country": "中国",
			},
			expectedLanguages: []string{"zh"},
		},
		{
			name: "Korean address",
			parsedAddress: map[string]string{
				"road":    "세종로",
				"city":    "서울",
				"country": "대한민국",
			},
			expectedLanguages: []string{"ko"},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			labels := make([]string, 0, len(tc.parsedAddress))
			values := make([]string, 0, len(tc.parsedAddress))
			for key, value := range tc.parsedAddress {
				labels = append(labels, key)
				values = append(values, value)
			}

			languages := PlaceLanguages(labels, values)
			if languages == nil {
				t.Fatalf("PlaceLanguages return empty: %v", languages)
			}

			if !reflect.DeepEqual(languages, tc.expectedLanguages) {
				t.Errorf("PlaceLanguages returned unexpected languages.\nGot:  %v\nWant: %v", languages, tc.expectedLanguages)
			}
		})
	}
}

func TestNearDupeHashesLanguages(t *testing.T) {
	testCases := []struct {
		name           string
		parsedAddress  map[string]string
		options        NearDupeHashOptions
		languages      []string
		expectedHashes []string
	}{
		{
			name: "English address",
			parsedAddress: map[string]string{
				"house_number": "123",
				"road":         "Main St",
				"city":         "New York",
				"state":        "NY",
				"postcode":     "10001",
			},
			options: NearDupeHashOptions{
				WithName:             false,
				WithAddress:          true,
				WithUnit:             false,
				WithCityOrEquivalent: true,
				WithPostalCode:       true,
				WithLatlon:           false,
				AddressOnlyKeys:      true,
			},
			languages: []string{"en"},
			expectedHashes: []string{
				"act|main saint|123|new york",
				"act|main street|123|new york",
				"act|main|123|new york",
				"apc|main saint|123|10001",
				"apc|main street|123|10001",
				"apc|main|123|10001",
			},
		},
		{
			name: "French address",
			parsedAddress: map[string]string{
				"house_number": "15",
				"road":         "Rue de la Paix",
				"city":         "Paris",
				"postcode":     "75002",
			},
			options: NearDupeHashOptions{
				WithName:             false,
				WithAddress:          true,
				WithUnit:             false,
				WithCityOrEquivalent: true,
				WithPostalCode:       true,
				WithLatlon:           false,
				AddressOnlyKeys:      true,
			},
			languages: []string{"fr"},
			expectedHashes: []string{
                "act|rue de la paix|15|paris", 
                "act|paix|15|paris", 
                "apc|rue de la paix|15|75002", 
                "apc|paix|15|75002",
			},
		},
		{
			name: "Multilingual address (Brussels)",
			parsedAddress: map[string]string{
				"house_number": "10",
				"road":         "Rue de la Loi",
				"city":         "Bruxelles",
				"postcode":     "1000",
			},
			options: NearDupeHashOptions{
				WithName:             false,
				WithAddress:          true,
				WithUnit:             false,
				WithCityOrEquivalent: true,
				WithPostalCode:       true,
				WithLatlon:           false,
				AddressOnlyKeys:      true,
			},
			languages: []string{"fr"},
			expectedHashes: []string{
                "act|rue de la loi|10|bruxelles", 
                "act|loi|10|bruxelles", 
                "apc|rue de la loi|10|1000", 
                "apc|loi|10|1000", 
			},
		},
		{
			name: "Japanese address with transliteration",
			parsedAddress: map[string]string{
				"house_number": "1",
				"road":         "丁目",
				"suburb":       "渋谷",
				"city":         "東京",
				"postcode":     "150-0042",
			},
			options: NearDupeHashOptions{
				WithName:             false,
				WithAddress:          true,
				WithUnit:             false,
				WithCityOrEquivalent: true,
				WithPostalCode:       true,
				WithLatlon:           false,
				AddressOnlyKeys:      true,
			},
			languages: []string{"ja"},
			expectedHashes: []string{
				"act|丁目|1|東京",
				"act|丁目|1|dongjing",
				"act|丁目|1|渋谷", 
				"act|丁目|1|segu",
				"act|dingmu|1|東京",
				"act|dingmu|1|dongjing",
				"act|dingmu|1|渋谷",
				"act|dingmu|1|segu",
				"apc|丁目|1|150-0042",
				"apc|dingmu|1|150-0042",
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			labels := make([]string, 0, len(tc.parsedAddress))
			values := make([]string, 0, len(tc.parsedAddress))
			for key, value := range tc.parsedAddress {
				labels = append(labels, key)
				values = append(values, value)
			}
	
			hashes := NearDupeLanguages(labels, values, tc.options, tc.languages)
			if hashes == nil {
				t.Fatalf("NearDupeHashesLanguages returned empty: %v", hashes)
			}
	
			if !reflect.DeepEqual(hashes, tc.expectedHashes) {
				t.Errorf("NearDupeHashesLanguages returned unexpected hashes for %s:\n", tc.name)
				t.Errorf("Got %d hashes, Want %d hashes\n", len(hashes), len(tc.expectedHashes))
				for i := 0; i < len(hashes) || i < len(tc.expectedHashes); i++ {
					if i < len(hashes) && i < len(tc.expectedHashes) {
						if hashes[i] != tc.expectedHashes[i] {
							t.Errorf("Hash %d:\nGot:  %s\nWant: %s", i, hashes[i], tc.expectedHashes[i])
						}
					} else if i < len(hashes) {
						t.Errorf("Extra hash %d: %s", i, hashes[i])
					} else {
						t.Errorf("Missing expected hash %d: %s", i, tc.expectedHashes[i])
					}
				}
			}
		})
	}
}


func TestNearDupeNameHashes(t *testing.T) {
	testCases := []struct {
		name           string
		input          string
		options        NormalizeOptions
		expectedHashes []string
	}{
		{
			name:  "Single word with default options: Atlantic",
			input: "Atlantic",
			options: GetDefaultNormalizeOptions(),
			expectedHashes: []string{
				"ATLN",
				"TLNT",
				"LNTK",
				"atla",
				"tlan",
				"lant",
				"anti",
				"ntic",
			},
		},
		{
			name:  "Roman Numeral: IV",
			input: "IV",
			options: GetDefaultNormalizeOptions(),
			expectedHashes: []string{
				"AF",
				"iv",
				"4",
			},
		},
		{
			name:  "Name with space: New York",
			input: "New York",
			options: GetDefaultNormalizeOptions(),
			expectedHashes: []string{
				"N",
				"NF",
				"new",
				"ARK",
				"york",
				"NRK",
			},
		},
		{
			name:  "Ordinal Number: 6th",
			input: "6th",
			options: GetDefaultNormalizeOptions(),
			expectedHashes: []string{
				"0",
				"T",
				"6th",
				"6",
				"th",
			},
		},
		{
			name:  "Spelled Numbers: Six",
			input: "Six",
			options: GetDefaultNormalizeOptions(),
			expectedHashes: []string{
				"6",
			},
		},
		{
			name: "Name with accents and lowercase option: Café",
			input: "Café",
			options: NormalizeOptions{
				Lowercase:    true,
				StripAccents: false,
			},
			expectedHashes: []string{
				"KF",
				"café", 
				"cafe",
			},
		},
		{
			name: "Name with numbers: 7Eleven",
			input: "7Eleven",
			options: GetDefaultNormalizeOptions(),
			expectedHashes: []string{
				"LFN", 
				"7ele",
				"elev",
				"leve",
				"even",
				"7-el",
				"-ele",
			},
		},
		{
			name: "Name with apostrophe and delete apostrophes option True: McDonald's",
			input: "McDonald's",
			options: NormalizeOptions{
				DeleteApostrophes: true,
			},
			expectedHashes: []string{
				"MKTN", 
				"KTNL",
				"TNLT", 
				"NLTS", 
				"McDo", 
				"cDon", 
				"Dona", 
				"onal", 
				"nald", 
				"alds",
			},
		},
		{
			name: "Name with apostrophe and delete apostrophes option False: McDonald's",
			input: "McDonald's",
			options: NormalizeOptions{
				DeleteApostrophes: false,
			},
			expectedHashes: []string{
				"MKTN", 
				"KTNL",
				"TNLT", 
				"NLTS", 
				"McDo", 
				"cDon", 
				"Dona", 
				"onal", 
				"nald", 
				"ald'",
				"ld's",
			},
		},
		// The Options here don't seem to do anything. Will require further research.
		// May not be an issue with implementation of cgo, but of original logic in near_dupe.c
		// {
		// 	name: "Non-English name with transliteration",
		// 	input: "München",
		// 	options: NormalizeOptions{
		// 		Transliterate: false,
		// 	},
		// 	expectedHashes: []string{
		// 		"MNXN", "MNKN",
		// 		"Muen", "uenc", "ench", "nche",
		// 		"chen", "Munc", "unch", "Münc", "ünch",
		// 	},
		// },
		// {
		// 	name: "Name with hyphen and delete hyphens option: Coca-Cola",
		// 	input: "Coca-Cola",
		// 	options: NormalizeOptions{
		// 		DeleteWordHyphens: false,
		// 	},
		// 	expectedHashes: []string{
		// 		"KKKL",
		// 		"Coca", 
		// 		"oca-", 
		// 		"ca-C", 
		// 		"a-Co", 
		// 		"-Col", 
		// 		"Cola", 
		// 		"K",
		// 	},
		// },
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			hashes := NearDupeNameOptions(tc.input, tc.options)
			if hashes == nil {
				t.Fatalf("NearDupeNameHashes returned empty: %v", hashes)
			}

			if !reflect.DeepEqual(hashes, tc.expectedHashes) {
				t.Errorf("NearDupeNameHashes returned unexpected hashes.\nGot:  %v\nWant: %v", hashes, tc.expectedHashes)
			}
		})
	}
}
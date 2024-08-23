// Package postal provides Go bindings for the libpostal C library,
// offering functionality for parsing, expanding, and generating
// near-dupe hashes for postal addresses.
package postal

/*
#cgo pkg-config: libpostal
#include <libpostal/libpostal.h>
#include <stdlib.h>

*/
import "C"

import (
	"log"
	"sync"
	"unicode/utf8"
	"unsafe"
)

var mu sync.Mutex

func init() {
    if (!bool(C.libpostal_setup()) || !bool(C.libpostal_setup_language_classifier())) {
        log.Fatal("Could not load libpostal")
    }
}

// NormalizeOptions represents the options allowed for name normalization for NearDupeNameHashes input.
// Corresponds to the C libpostal_normalize_options_t struct.
type NormalizeOptions struct {
    Languages []string
    AddressComponents uint16
    LatinAscii bool
    Transliterate bool
    StripAccents bool
    Decompose bool
    Lowercase bool
    TrimString bool
    ReplaceWordHyphens bool
    DeleteWordHyphens bool
    ReplaceNumericHyphens bool
    DeleteNumericHyphens bool
    SplitAlphaFromNumeric bool
    DeleteFinalPeriods bool
    DeleteAcronymPeriods bool
    DropEnglishPossessives bool
    DeleteApostrophes bool
    ExpandNumex bool
    RomanNumerals bool
}

// NearDupeHashOptions represents the options allowed for near-dupe hashing.
// It corresponds to the C libpostal_near_dupe_hash_options_t struct.
type NearDupeHashOptions struct {
    WithName bool
    WithAddress bool
    WithUnit bool
    WithCityOrEquivalent bool
    WithSmallContainingBoundaries bool
    WithPostalCode bool
    WithLatlon bool
    Latitude float64
    Longitude float64
    GeohashPrecision uint32
    NameAndAddressKeys bool
    NameOnlyKeys bool
    AddressOnlyKeys bool
}

// Fetch the default Libpostal C options for Near Dupe Hash and Normalization
var cDefaultOptions = C.libpostal_get_default_options()

// #define DEFAULT_NEAR_DUPE_GEOHASH_PRECISION 6
// static libpostal_near_dupe_hash_options_t LIBPOSTAL_NEAR_DUPE_HASH_DEFAULT_OPTIONS = {
//     .with_name = true,
//     .with_address = true,
//     .with_unit = false,
//     .with_city_or_equivalent = true,
//     .with_small_containing_boundaries = true,
//     .with_postal_code = true,
//     .with_latlon = false,
//     .latitude = 0.0,
//     .longitude = 0.0,
//     .geohash_precision = DEFAULT_NEAR_DUPE_GEOHASH_PRECISION,
//     .name_and_address_keys = true,
//     .name_only_keys = false,
//     .address_only_keys = false
// };
var cHashDefaultOptions = C.libpostal_get_near_dupe_hash_default_options()

// GetDefaultNormalizeOptions returns the default options for name normalization.
// Initializes a NormalizeOptions struct with the default values from libpostal.
//
// Returns:
//   - NormalizeOptions: A struct containing the default normalization options.
func GetDefaultNormalizeOptions() NormalizeOptions {
	return NormalizeOptions{
		Languages: nil,
		AddressComponents: uint16(cDefaultOptions.address_components),
		LatinAscii: bool(cDefaultOptions.latin_ascii),
		Transliterate: bool(cDefaultOptions.transliterate),
		StripAccents: bool(cDefaultOptions.strip_accents),
		Decompose: bool(cDefaultOptions.decompose),
		Lowercase: bool(cDefaultOptions.lowercase),
		TrimString: bool(cDefaultOptions.trim_string),
		ReplaceWordHyphens: bool(cDefaultOptions.replace_word_hyphens),
		DeleteWordHyphens: bool(cDefaultOptions.delete_word_hyphens),
		ReplaceNumericHyphens: bool(cDefaultOptions.replace_numeric_hyphens),
		DeleteNumericHyphens: bool(cDefaultOptions.delete_numeric_hyphens),
		SplitAlphaFromNumeric: bool(cDefaultOptions.split_alpha_from_numeric),
		DeleteFinalPeriods: bool(cDefaultOptions.delete_final_periods),
		DeleteAcronymPeriods: bool(cDefaultOptions.delete_acronym_periods),
		DropEnglishPossessives: bool(cDefaultOptions.drop_english_possessives),
		DeleteApostrophes: bool(cDefaultOptions.delete_apostrophes),
		ExpandNumex: bool(cDefaultOptions.expand_numex),
		RomanNumerals: bool(cDefaultOptions.roman_numerals),
	}
}

// GetDefaultNearDupeHashOptions returns the default options for near-dupe hashing.
// Initializes a NearDupeHashOptions struct with the default values from libpostal.
//
// Returns:
//   - NearDupeHashOptions: A struct containing the default near-dupe hash options.
func GetDefaultNearDupeHashOptions() NearDupeHashOptions {
    return NearDupeHashOptions{
        WithName: bool(cHashDefaultOptions.with_name),
        WithAddress: bool(cHashDefaultOptions.with_address),
        WithUnit: bool(cHashDefaultOptions.with_unit),
        WithCityOrEquivalent: bool(cHashDefaultOptions.with_city_or_equivalent),
        WithSmallContainingBoundaries: bool(cHashDefaultOptions.with_small_containing_boundaries),
        WithPostalCode: bool(cHashDefaultOptions.with_postal_code),
        WithLatlon: bool(cHashDefaultOptions.with_latlon),
        Latitude: float64(cHashDefaultOptions.latitude),
        Longitude: float64(cHashDefaultOptions.longitude),
        GeohashPrecision: uint32(cHashDefaultOptions.geohash_precision),
        NameAndAddressKeys: bool(cHashDefaultOptions.name_and_address_keys),
        NameOnlyKeys: bool(cHashDefaultOptions.name_only_keys),
        AddressOnlyKeys: bool(cHashDefaultOptions.address_only_keys),
    }
}

var libpostalDefaultOptions = GetDefaultNormalizeOptions()
var libpostalDefaultHashOptions = GetDefaultNearDupeHashOptions()

// NearDupeNameOptions generates near-dupe hashes for a given name
// using the specified normalization options.
//
// Parameters:
//   - name: A string representing the name to generate hashes for.
//   - options: NormalizeOptions specifying the normalization configuration.
//
// Returns:
//   - []string: A slice of strings containing the generated near-dupe hashes.
//     Returns nil if the input name is not a valid UTF-8 string.
func NearDupeNameOptions(name string, options NormalizeOptions) []string {
    if !utf8.ValidString(name) {
        return nil
    }

	mu.Lock()
	defer mu.Unlock()

	cName := C.CString(name)
	defer C.free(unsafe.Pointer(cName))

    var char_ptr *C.char
    ptr_size := unsafe.Sizeof(char_ptr)

    cOptions := C.libpostal_get_default_options()
    if options.Languages != nil {
        cLanguages := C.calloc(C.size_t(len(options.Languages)), C.size_t(ptr_size))
        cLanguagesPtr := (*[1<<30](*C.char))(unsafe.Pointer(cLanguages))

        defer C.free(unsafe.Pointer(cLanguages))

        for i := 0; i < len(options.Languages); i++ {
            cLang := C.CString(options.Languages[i])
            defer C.free(unsafe.Pointer(cLang))
            cLanguagesPtr[i] = cLang
        }

        cOptions.languages = (**C.char)(cLanguages)
        cOptions.num_languages = C.size_t(len(options.Languages))
    } else {
        cOptions.num_languages = 0
    }

    cOptions.address_components = C.uint16_t(options.AddressComponents)
    cOptions.latin_ascii = C.bool(options.LatinAscii)
    cOptions.transliterate = C.bool(options.Transliterate)
    cOptions.strip_accents = C.bool(options.StripAccents)
    cOptions.decompose = C.bool(options.Decompose)
    cOptions.lowercase = C.bool(options.Lowercase)
    cOptions.trim_string = C.bool(options.TrimString)
    cOptions.replace_word_hyphens = C.bool(options.ReplaceWordHyphens)
    cOptions.delete_word_hyphens = C.bool(options.DeleteWordHyphens)
    cOptions.replace_numeric_hyphens = C.bool(options.ReplaceNumericHyphens)
    cOptions.delete_numeric_hyphens = C.bool(options.DeleteNumericHyphens)
    cOptions.split_alpha_from_numeric = C.bool(options.SplitAlphaFromNumeric)
    cOptions.delete_final_periods = C.bool(options.DeleteFinalPeriods)
    cOptions.delete_acronym_periods = C.bool(options.DeleteAcronymPeriods)
    cOptions.drop_english_possessives = C.bool(options.DropEnglishPossessives)
    cOptions.delete_apostrophes = C.bool(options.DeleteApostrophes)
    cOptions.expand_numex = C.bool(options.ExpandNumex)
    cOptions.roman_numerals = C.bool(options.RomanNumerals)

	var cNumHashes = C.size_t(0)

	cHashes := C.libpostal_near_dupe_name_hashes(cName, cOptions, &cNumHashes)
	defer C.free(unsafe.Pointer(cHashes))

	return cStringArrayToStringSlice(cHashes, cNumHashes)
}

// NearDupeNames generates near-dupe hashes for a given name using default options.
//
// Parameters:
//   - name: A string representing the name to generate hashes for.
//
// Returns:
//   - []string: A slice of strings containing the generated near-dupe hashes.
//     Returns nil if the input name is not a valid UTF-8 string.
func NearDupeNames(name string) ([]string) {
	return NearDupeNameOptions(name, libpostalDefaultOptions)
}

// NearDupeHashesWithOptions generates near-dupe hashes for the given components
// with custom options and optional languages.
//
// Parameters:
//   - labels: A slice of strings representing the labels of address components.
//   - values: A slice of strings representing the values of address components.
//   - options: NearDupeHashOptions specifying the hashing configuration.
//   - languages: A slice of strings representing 2-letter ISO language codes to consider.
//
// Returns:
//   - []string: A slice of strings containing the generated near-dupe hashes.
//     Returns nil if the input slices have different lengths or are empty.
func NearDupeOptions(labels []string, values []string, options NearDupeHashOptions, languages []string) []string {
    if len(labels) != len(values) {
        return nil
    }

    mu.Lock()
    defer mu.Unlock()

    numComponents := len(labels)
    if numComponents == 0 {
        return nil
    }

    cLabels := make([]*C.char, numComponents)
    cValues := make([]*C.char, numComponents)

    for i := 0; i < numComponents; i++ {
        cLabels[i] = C.CString(labels[i])
        cValues[i] = C.CString(values[i])
        defer C.free(unsafe.Pointer(cLabels[i]))
        defer C.free(unsafe.Pointer(cValues[i]))
    }

    cOptions := C.libpostal_get_near_dupe_hash_default_options()
    cOptions.with_name = C.bool(options.WithName)
    cOptions.with_address = C.bool(options.WithAddress)
    cOptions.with_unit = C.bool(options.WithUnit)
    cOptions.with_city_or_equivalent = C.bool(options.WithCityOrEquivalent)
    cOptions.with_small_containing_boundaries = C.bool(options.WithSmallContainingBoundaries)
    cOptions.with_postal_code = C.bool(options.WithPostalCode)
    cOptions.with_latlon = C.bool(options.WithLatlon)
    cOptions.latitude = C.double(options.Latitude)
    cOptions.longitude = C.double(options.Longitude)
    cOptions.geohash_precision = C.uint32_t(options.GeohashPrecision)
    cOptions.name_and_address_keys = C.bool(options.NameAndAddressKeys)
    cOptions.name_only_keys = C.bool(options.NameOnlyKeys)
    cOptions.address_only_keys = C.bool(options.AddressOnlyKeys)


    var cNumHashes C.size_t
    var cHashes **C.char

    if len(languages) > 0 {
        cLanguages := make([]*C.char, len(languages))
        for i, lang := range languages {
            cLanguages[i] = C.CString(lang)
            defer C.free(unsafe.Pointer(cLanguages[i]))
        }

        cHashes = C.libpostal_near_dupe_hashes_languages(
            C.size_t(numComponents),
            (**C.char)(unsafe.Pointer(&cLabels[0])),
            (**C.char)(unsafe.Pointer(&cValues[0])),
            cOptions,
            C.size_t(len(languages)),
            (**C.char)(unsafe.Pointer(&cLanguages[0])),
            &cNumHashes,
        )
    } else {
        cHashes = C.libpostal_near_dupe_hashes(
            C.size_t(numComponents),
            (**C.char)(unsafe.Pointer(&cLabels[0])),
            (**C.char)(unsafe.Pointer(&cValues[0])),
            cOptions,
            &cNumHashes,
        )
    }
    defer C.free(unsafe.Pointer(cHashes))

    return cStringArrayToStringSlice(cHashes, cNumHashes)
}

// NearDupe generates near-dupe hashes for the given components using default options.
//
// Parameters:
//   - labels: A slice of strings representing the labels of address components.
//   - values: A slice of strings representing the values of address components.
//   - options: NearDupeHashOptions specifying the hashing configuration.
//
// Returns:
//   - []string: A slice of strings containing the generated near-dupe hashes.
//     Returns nil if the input slices have different lengths or are empty.
func NearDupe(labels []string, values []string, options NearDupeHashOptions) []string {
    return NearDupeOptions(labels, values, options, nil)
}

// NearDupe generates near-dupe hashes for the given components using default options.
//
// Parameters:
//   - labels: A slice of strings representing the labels of address components.
//   - values: A slice of strings representing the values of address components.
//
// Returns:
//   - []string: A slice of strings containing the generated near-dupe hashes.
//     Returns nil if the input slices have different lengths or are empty.
func NearDupeDefaultOptions(labels []string, values []string) []string {
    return NearDupeOptions(labels, values, libpostalDefaultHashOptions, nil)
}

// NearDupeLanguages generates near-dupe hashes for the given components
// with specified languages using default options.
//
// Parameters:
//   - labels: A slice of strings representing the labels of address components.
//   - values: A slice of strings representing the values of address components.
//   - options: NearDupeHashOptions specifying the hashing configuration.
//   - languages: A slice of strings representing language codes to consider.
//
// Returns:
//   - []string: A slice of strings containing the generated near-dupe hashes.
//     Returns nil if the input slices have different lengths or are empty.
func NearDupeLanguages(labels []string, values []string, options NearDupeHashOptions, languages []string) []string {
    return NearDupeOptions(labels, values, options, languages)
}

// PlaceLanguages returns the languages for the given address components.
//
// Parameters:
//   - labels: A slice of strings representing the labels of address components.
//   - values: A slice of strings representing the values of address components.
//
// Returns:
//   - []string: A slice of strings containing the detected languages.
//     Returns nil if the input slices have different lengths or are empty.
func PlaceLanguages(labels []string, values []string) []string {
    if len(labels) != len(values) {
        return nil
    }

    mu.Lock()
    defer mu.Unlock()

    numComponents := len(labels)
    if numComponents == 0 {
        return nil
    }

	cLabels := make([]*C.char, numComponents)
	cValues := make([]*C.char, numComponents)

    for i := 0; i < numComponents; i++ {
        cLabels[i] = C.CString(labels[i])
        cValues[i] = C.CString(values[i])
        defer C.free(unsafe.Pointer(cLabels[i]))
        defer C.free(unsafe.Pointer(cValues[i]))
    }

	var cNumLanguages = C.size_t(0)

	cLanguages := C.libpostal_place_languages(
		C.size_t(len(labels)),
		(**C.char)(unsafe.Pointer(&cLabels[0])),
		(**C.char)(unsafe.Pointer(&cValues[0])),
		&cNumLanguages,
	)
	defer C.free(unsafe.Pointer(cLanguages))

	return cStringArrayToStringSlice(cLanguages, cNumLanguages)
}

// cStringArrayToStringSlice converts a C array of strings to a Go slice of strings.
//
// This function is an internal helper used to bridge between C and Go data structures.
// It takes a pointer to a C array of strings and its size, then creates a new Go
// slice and populates it with the strings from the C array.
//
// Parameters:
//   - cArray: A pointer to a C array of strings (type **C.char).
//   - arraySize: The size of the C array (type C.size_t).
//
// Returns:
//   - []string: A Go slice containing the strings from the C array.
//
// Note: This function assumes that the C array is properly null-terminated and
// that the arraySize accurately reflects the number of strings in the array.
// Callers are responsible for freeing the original C array after using this function.
func cStringArrayToStringSlice(cArray **C.char, arraySize C.size_t) []string {
    slice := make([]string, int(arraySize))
    cArrayPtr := (*[1<<30](*C.char))(unsafe.Pointer(cArray))
    
    var i uint64
    for i = 0; i < uint64(arraySize); i++ {
        slice[i] = C.GoString(cArrayPtr[i])
    }
    return slice
}

// Frees resources allocated by libpostal
func NearDupeTeardown() {
    mu.Lock()
    defer mu.Unlock()
    C.libpostal_teardown()
    C.libpostal_teardown_language_classifier()
}

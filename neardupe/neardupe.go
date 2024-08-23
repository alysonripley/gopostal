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

var cDefaultOptions = C.libpostal_get_default_options()
var cHashDefaultOptions = C.libpostal_get_near_dupe_hash_default_options()

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

func NearDupeNames(name string) ([]string) {
	return NearDupeNameOptions(name, libpostalDefaultOptions)
}

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

func NearDupe(labels []string, values []string, options NearDupeHashOptions) []string {
    return NearDupeOptions(labels, values, options, nil)
}

func NearDupeDefaultOptions(labels []string, values []string) []string {
    return NearDupeOptions(labels, values, libpostalDefaultHashOptions, nil)
}

func NearDupeLanguages(labels []string, values []string, options NearDupeHashOptions, languages []string) []string {
    return NearDupeOptions(labels, values, options, languages)
}

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

func cStringArrayToStringSlice(cArray **C.char, arraySize C.size_t) []string {
    slice := make([]string, int(arraySize))
    cArrayPtr := (*[1<<30](*C.char))(unsafe.Pointer(cArray))
    
    var i uint64
    for i = 0; i < uint64(arraySize); i++ {
        slice[i] = C.GoString(cArrayPtr[i])
    }
    return slice
}

func NearDupeTeardown() {
    mu.Lock()
    defer mu.Unlock()
    C.libpostal_teardown()
    C.libpostal_teardown_language_classifier()
}

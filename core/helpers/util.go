package helpers

import (
	"fmt"
	"log"
	"slices"
	"sort"
	"strings"
	"time"
	"unicode"
	    "golang.org/x/text/language"
    "golang.org/x/text/message"
)

func IsEmpty(v string) bool {
	ret := false
	if len(v) == 0 {
		ret = true
	}
	return ret
}

func IsEmptyLog(varCheck, varText string, isRequired bool) {
	if IsEmpty(varCheck) {
		if isRequired {
			log.Fatalf("Entry %s are Required %s", varText, ICON_CANCEL_RED)
		} else {
			log.Printf(varText+" %s", ICON_CANCEL)
		}
	} else {
		log.Printf(varText+" %s", ICON_OK)
	}
}

func IsNotListLog(varCheck string, dataMaster []string, varText string, isRequired bool) {
	err := slices.Contains(dataMaster, varCheck)
	if err {
		log.Printf(varText+" %s", ICON_OK)
	} else {
		if isRequired {
			log.Fatalf("Entry %s are Required %s", varText, ICON_CANCEL_RED)
		} else {
			log.Printf(varText+" %s", ICON_CANCEL)
		}
	}
}

func IsError(err error, varText string, isRequired bool) {
	if err != nil {
		if isRequired {
			log.Fatalf("%s are Required %s, error %v", varText, ICON_CANCEL_RED, err)
		} else {
			log.Printf(varText+" %s", ICON_CANCEL)
		}
	} else {
		log.Printf(varText+" %s", ICON_OK)
	}
}

func ToCamelCase(s string) string {
	parts := strings.Split(s, "_")
	for i := range parts {
		parts[i] = strings.Title(parts[i])
	}
	return strings.Join(parts, "")
}

func ToSnakeCase(s string) string {
	var result []rune
	for i, r := range s {
		if unicode.IsUpper(r) && i > 0 {
			result = append(result, '_')
		}
		result = append(result, unicode.ToLower(r))
	}
	return string(result)
}

func FormatDataAsTable(data []map[string]interface{}) string {
	if len(data) == 0 {
		return "No data found."
	}
	// Collect all keys
	keys := make(map[string]bool)
	var header []string
	for _, row := range data {
		for k := range row {
			if !keys[k] {
				keys[k] = true
				header = append(header, k)
			}
		}
	}
	sort.Strings(header)

	var b strings.Builder
	// Header
	for _, h := range header {
		b.WriteString("| " + h + " ")
	}
	b.WriteString("|\n")
	// Separator
	for range header {
		b.WriteString("| --- ")
	}
	b.WriteString("|\n")
	// Rows
	for _, row := range data {
		for _, h := range header {
			val := ""
			if v, ok := row[h]; ok {
				val = fmt.Sprintf("%v", v)
			}
			// Sanitize newlines
			val = strings.ReplaceAll(val, "\n", " ")
			b.WriteString("| " + val + " ")
		}
		b.WriteString("|\n")
	}
	return b.String()
}

func FormatThousands(n int64) string {
    p := message.NewPrinter(language.Indonesian)
    return p.Sprintf("%d", n)
}

func FormatRupiah(n int64) string {
    return "Rp " + FormatThousands(n)
}

func FormatDateIndonesianWithFormat(dateStr string, format string) string {
	// Indonesian month names
	monthNames := map[int]string{
		1: "Januari", 2: "Februari", 3: "Maret", 4: "April",
		5: "Mei", 6: "Juni", 7: "Juli", 8: "Agustus",
		9: "September", 10: "Oktober", 11: "November", 12: "Desember",
	}
	
	// Indonesian day names
	dayNames := map[int]string{
		0: "Minggu", 1: "Senin", 2: "Selasa", 3: "Rabu",
		4: "Kamis", 5: "Jumat", 6: "Sabtu",
	}
	
	// Try multiple date formats
	formats := []string{
		"2006-01-02T15:04:05Z07:00", // ISO8601 with timezone
		"2006-01-02T15:04:05Z",      // ISO8601 UTC
		"2006-01-02 15:04:05",       // MySQL datetime
		"2006-01-02",                // Date only
		"02/01/2006",                // DD/MM/YYYY
		"01/02/2006",                // MM/DD/YYYY
		"2006/01/02",                // YYYY/MM/DD
	}
	
	var parsedDate time.Time
	var err error
	
	for _, f := range formats {
		parsedDate, err = time.Parse(f, dateStr)
		if err == nil {
			break
		}
	}
	
	// If parsing failed, return original string
	if err != nil {
		return dateStr
	}
	
	day := parsedDate.Day()
	monthName := monthNames[int(parsedDate.Month())]
	year := parsedDate.Year()
	
	// Format based on requested format
	switch format {
	case "long":
		// "Senin, 6 Januari 2026"
		dayName := dayNames[int(parsedDate.Weekday())]
		return fmt.Sprintf("%s, %d %s %d", dayName, day, monthName, year)
	case "short":
		// "6 Januari 2026"
		return fmt.Sprintf("%d %s %d", day, monthName, year)
	case "numeric":
		// "06/01/2026"
		return fmt.Sprintf("%02d/%02d/%d", day, int(parsedDate.Month()), year)
	default:
		// Default to long format
		dayName := dayNames[int(parsedDate.Weekday())]
		return fmt.Sprintf("%s, %d %s %d", dayName, day, monthName, year)
	}
}

// FormatDateIndonesian formats date with default "long" format for backward compatibility
func FormatDateIndonesian(dateStr string) string {
	return FormatDateIndonesianWithFormat(dateStr, "long")
}
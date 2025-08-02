package parser

import (
	"regexp"
	"strings"

	"lwnra-devo-api/models"
)

// ParseDevotional parses a Facebook post message into a Devotional struct
func ParseDevotional(msg string) models.Devotional {
	var devo models.Devotional
	lines := strings.Split(msg, "\n")

	// Find values between curly braces {} or after known section headers
	devo.Reading = findReadingAfterPrefix(lines, "Read")
	devo.Date = normalizeDate(findBraceThatLooksLikeDate(lines))
	devo.Version = findBibleVersion(lines, devo.Reading)
	devo.Passage = grabPassageAfterVersion(lines, devo.Reading, devo.Version)
	devo.ReflectionQs = getLinesBetween(lines, "REFLECTION QUESTIONS", "")
	devo.Title = findFirstBracketSection(lines, "WHEN", "WATCHING") // crude, works w/ "WHEN NO ONE IS WATCHING"
	if devo.Title == "" {
		devo.Title = findCapsSectionAfter(lines, "REFLECTION QUESTIONS")
	}
	devo.Author = findAfterTitle(lines, devo.Title)
	devo.Body = grabSection(lines, devo.Author, "PRAYER")
	devo.Prayer = grabPrayerSection(lines) // Use dedicated function for prayer

	// Clean up
	devo.Passage = strings.Trim(devo.Passage, "{} \n")
	devo.Body = strings.Trim(devo.Body, "{} \n")
	devo.Prayer = strings.Trim(devo.Prayer, "{} \n")

	// Reflection questions: only keep up to, but not including, lines starting with '{' or 'WHEN'
	if len(devo.ReflectionQs) > 0 {
		var qs []string
		for _, l := range devo.ReflectionQs {
			if strings.HasPrefix(l, "{") || strings.HasPrefix(l, "WHEN ") {
				break
			}
			if l != "" {
				qs = append(qs, strings.Trim(l, "{} "))
			}
		}
		devo.ReflectionQs = qs
	}

	return devo
}

func findBraceAfterPrefix(lines []string, prefix string) string {
	for i, l := range lines {
		if strings.HasPrefix(l, prefix) && strings.Contains(l, "{") && strings.Contains(l, "}") {
			return extractCurly(l)
		}
		// Sometimes it's on the next line
		if strings.HasPrefix(l, prefix) && i+1 < len(lines) && strings.HasPrefix(strings.TrimSpace(lines[i+1]), "{") {
			return extractCurly(lines[i+1])
		}
	}
	return ""
}

// findReadingAfterPrefix finds the reading reference after "Read" line, handling both braced and unbraced formats
func findReadingAfterPrefix(lines []string, prefix string) string {
	for i, l := range lines {
		// Check if line starts with the prefix (e.g., "Read")
		if strings.HasPrefix(l, prefix) {
			// First, check if reading is on the same line (e.g., "Read Matthew 6:16-18")
			if len(l) > len(prefix) {
				reading := strings.TrimSpace(l[len(prefix):])
				// Check if it contains braces
				if strings.Contains(reading, "{") && strings.Contains(reading, "}") {
					return extractCurly(reading)
				}
				// If no braces, return the text after "Read"
				if reading != "" {
					return reading
				}
			}

			// If not on same line, check the next line
			if i+1 < len(lines) {
				nextLine := strings.TrimSpace(lines[i+1])
				if nextLine != "" {
					// Check if it's a braced format
					if strings.HasPrefix(nextLine, "{") && strings.HasSuffix(nextLine, "}") {
						return extractCurly(nextLine)
					}
					// Check if it looks like a Bible reference (before the date)
					if !isDateLine(nextLine) {
						return nextLine
					}
				}
			}
		}
	}
	return ""
}

func findBraceThatLooksLikeDate(lines []string) string {
	// Look for date in format: Month Day, Year (with or without curly braces)
	reBraced := regexp.MustCompile(`\{(January|February|March|April|May|June|July|August|September|October|November|December) [0-9]{1,2}, [0-9]{4}\}`)
	reUnbraced := regexp.MustCompile(`^(January|February|March|April|May|June|July|August|September|October|November|December) [0-9]{1,2}, [0-9]{4}$`)

	// First, try to find it after "Read" line for better accuracy
	foundRead := false
	for _, l := range lines {
		l = strings.TrimSpace(l)
		if !foundRead {
			if strings.HasPrefix(l, "Read") {
				foundRead = true
			}
			continue
		}

		// After "Read" line, look for a valid date line
		// First try braced format
		m := reBraced.FindStringSubmatch(l)
		if len(m) > 0 {
			return strings.Trim(m[0], "{}")
		}

		// Then try unbraced format (like "August 2, 2025")
		m = reUnbraced.FindStringSubmatch(l)
		if len(m) > 0 {
			return m[0]
		}

		// Stop searching if we hit a non-empty line that's not a date after "Read"
		if foundRead && l != "" && !isDateLine(l) {
			break
		}
	}

	// If not found after "Read", search entire content for any date pattern
	for _, l := range lines {
		l = strings.TrimSpace(l)

		// Try braced format first
		m := reBraced.FindStringSubmatch(l)
		if len(m) > 0 {
			return strings.Trim(m[0], "{}")
		}

		// Then try unbraced format
		m = reUnbraced.FindStringSubmatch(l)
		if len(m) > 0 {
			return m[0]
		}
	}

	return ""
}

func grabSection(lines []string, start string, until string) string {
	startIdx := -1
	untilIdx := len(lines)
	for i, l := range lines {
		// For finding the start section
		if strings.Contains(strings.ToUpper(l), strings.ToUpper(start)) && startIdx == -1 {
			startIdx = i + 1 // exclude header line
		}
		// For finding the until section - be strict about "PRAYER" being all caps
		if until != "" && startIdx != -1 && i > startIdx {
			if until == "PRAYER" {
				// Strict matching: must contain exactly "PRAYER" in all caps
				if strings.Contains(l, "PRAYER") {
					untilIdx = i
					break
				}
			} else {
				// For other sections, use case-insensitive matching
				if strings.Contains(strings.ToUpper(l), strings.ToUpper(until)) {
					untilIdx = i
					break
				}
			}
		}
	}
	if startIdx >= 0 && startIdx < untilIdx {
		return strings.TrimSpace(strings.Join(lines[startIdx:untilIdx], "\n"))
	}
	return ""
}

func getLinesBetween(lines []string, marker string, stopAt string) []string {
	var found bool
	var out []string
	for i, l := range lines {
		if found {
			if stopAt != "" && strings.Contains(strings.ToUpper(l), strings.ToUpper(stopAt)) {
				break
			}
			if strings.TrimSpace(l) != "" {
				out = append(out, strings.TrimSpace(l))
			}
		}
		if strings.Contains(strings.ToUpper(l), strings.ToUpper(marker)) {
			found = true
		}
		// end if hit all-caps line that is probably title
		if found && i > 0 && isLikelyTitle(l) {
			break
		}
	}
	return out
}

func findFirstBracketSection(lines []string, mustContain, mustContain2 string) string {
	for _, l := range lines {
		if strings.Contains(strings.ToUpper(l), strings.ToUpper(mustContain)) && strings.Contains(strings.ToUpper(l), strings.ToUpper(mustContain2)) {
			return strings.Trim(l, "{} ")
		}
	}
	for _, l := range lines {
		if isLikelyTitle(l) {
			return strings.Trim(l, "{} ")
		}
	}
	return ""
}

func isLikelyTitle(line string) bool {
	line = strings.TrimSpace(line)
	return len(line) > 0 && strings.ToUpper(line) == line && len(line) > 5 && !strings.HasPrefix(line, "PRAYER") && !strings.Contains(line, "REFLECTION")
}

// Try to find author, assuming it's the line after title
func findAfterTitle(lines []string, title string) string {
	for i, l := range lines {
		if strings.Trim(l, "{} ") == title && i+1 < len(lines) && len(strings.TrimSpace(lines[i+1])) > 0 {
			return strings.Trim(lines[i+1], "{} ")
		}
	}
	return ""
}

// Extract first {...} group from string
func extractCurly(s string) string {
	i1 := strings.Index(s, "{")
	i2 := strings.Index(s, "}")
	if i1 >= 0 && i2 > i1 {
		return strings.TrimSpace(s[i1+1 : i2])
	}
	return ""
}

func findCapsSectionAfter(lines []string, after string) string {
	found := false
	for _, l := range lines {
		if found && isLikelyTitle(l) {
			return strings.Trim(l, "{} ")
		}
		if strings.Contains(strings.ToUpper(l), strings.ToUpper(after)) {
			found = true
		}
	}
	return ""
}

// grabPrayerSection specifically looks for "PRAYER" in all caps and captures everything after it
func grabPrayerSection(lines []string) string {
	startIdx := -1
	for i, l := range lines {
		// Look for exact "PRAYER" in all caps (can be surrounded by other text)
		if strings.Contains(l, "PRAYER") {
			startIdx = i + 1 // start from the line after "PRAYER"
			break
		}
	}

	if startIdx >= 0 && startIdx < len(lines) {
		// Capture everything from after "PRAYER" to the end
		return strings.TrimSpace(strings.Join(lines[startIdx:], "\n"))
	}
	return ""
}

// normalizeDate ensures the date is in consistent "Month Day, Year" format
func normalizeDate(dateStr string) string {
	if dateStr == "" {
		return ""
	}

	// The regex already ensures the format is "Month Day, Year"
	// This function can be extended later if we need to handle other date formats
	dateStr = strings.TrimSpace(dateStr)

	// Ensure proper spacing around comma (in case of variations like "August 2,2025")
	re := regexp.MustCompile(`([A-Za-z]+) ([0-9]{1,2}),?\s*([0-9]{4})`)
	matches := re.FindStringSubmatch(dateStr)
	if len(matches) == 4 {
		return matches[1] + " " + matches[2] + ", " + matches[3]
	}

	return dateStr
}

// isDateLine checks if a line looks like a date (with or without braces)
func isDateLine(line string) bool {
	line = strings.TrimSpace(line)
	// Check for braced date format
	reBraced := regexp.MustCompile(`^\{(January|February|March|April|May|June|July|August|September|October|November|December) [0-9]{1,2}, [0-9]{4}\}$`)
	if reBraced.MatchString(line) {
		return true
	}

	// Check for unbraced date format
	reUnbraced := regexp.MustCompile(`^(January|February|March|April|May|June|July|August|September|October|November|December) [0-9]{1,2}, [0-9]{4}$`)
	return reUnbraced.MatchString(line)
}

// findBibleVersion finds the Bible version (like NIV, ESV, NASB, AMPLIFIED, MESSAGE) after the reading reference
func findBibleVersion(lines []string, reading string) string {
	if reading == "" {
		return ""
	}

	// Look for lines that contain the reading reference followed by a version
	for _, l := range lines {
		l = strings.TrimSpace(l)
		// Skip empty lines and lines that are just the date
		if l == "" || isDateLine(l) {
			continue
		}

		// Look for lines that start with the reading reference
		if strings.HasPrefix(l, reading) {
			// Extract everything after the reading reference
			afterReading := strings.TrimSpace(l[len(reading):])
			if afterReading != "" {
				// Bible versions pattern - can be 2+ uppercase letters or common longer names
				versionRegex := regexp.MustCompile(`^([A-Z]{2,}|AMPLIFIED|MESSAGE|PHILLIPS|CONTEMPORARY|LIVING|PASSION|VOICE)\b`)
				matches := versionRegex.FindStringSubmatch(afterReading)
				if len(matches) > 1 {
					return matches[1]
				}
			}
		}

		// Also check for patterns like "Matthew 6:16-18 NIV" anywhere in the line
		if strings.Contains(l, reading) {
			// Find the position after the reading reference
			idx := strings.Index(l, reading) + len(reading)
			if idx < len(l) {
				afterReading := strings.TrimSpace(l[idx:])
				versionRegex := regexp.MustCompile(`^([A-Z]{2,}|AMPLIFIED|MESSAGE|PHILLIPS|CONTEMPORARY|LIVING|PASSION|VOICE)\b`)
				matches := versionRegex.FindStringSubmatch(afterReading)
				if len(matches) > 1 {
					return matches[1]
				}
			}
		}
	}

	return ""
}

// grabPassageAfterVersion extracts the passage text that comes after the version line
func grabPassageAfterVersion(lines []string, reading string, version string) string {
	if reading == "" {
		return ""
	}

	startIdx := -1
	for i, l := range lines {
		l = strings.TrimSpace(l)
		// Look for the line that contains the reading and version
		if strings.Contains(l, reading) && (version == "" || strings.Contains(l, version)) {
			startIdx = i + 1 // Start from the line after the version line
			break
		}
	}

	// If we found the version line, extract everything until "REFLECTION QUESTIONS"
	if startIdx >= 0 && startIdx < len(lines) {
		endIdx := len(lines)
		// Look for "REFLECTION QUESTIONS" to know where to stop
		for i := startIdx; i < len(lines); i++ {
			if strings.Contains(strings.ToUpper(lines[i]), "REFLECTION QUESTIONS") {
				endIdx = i
				break
			}
		}

		if startIdx < endIdx {
			return strings.TrimSpace(strings.Join(lines[startIdx:endIdx], "\n"))
		}
	}

	return ""
}

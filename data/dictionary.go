package data

import (
	"sort"
)

type Dictionary struct {
	Entries []string
}

func (dict *Dictionary) Add(word string) {
	dict.Entries = append(dict.Entries, word)
}

// Contains determines if the given word exists in the Dictionary object's entries
func (dict *Dictionary) Contains(word string) bool {
	for _, entry := range dict.Entries {
		if entry == word {
			return true
		}
	}

	return false
}

// Length calls the len function on the Dictionary's entries
func (dict *Dictionary) Length() int {
	return len(dict.Entries)
}

// NearestMatch Finds the smallest leveshtein distance alphabetically from the given word
func (dict *Dictionary) NearestMatch(word string) string {
	min := -1
	levenshteinMap := make(map[string]int)
	for _, entry := range dict.Entries {
		distance := levenshteinDistance(word, entry)
		if min == -1 || distance < min {
			min = distance
		}
		levenshteinMap[entry] = distance
	}

	var closestMatches []string
	for entry, distance := range levenshteinMap {
		if distance == min {
			closestMatches = append(closestMatches, entry)
		}
	}

	// If there are multiple possible matches, return the first alphabetically
	if len(closestMatches) > 1 {
		sort.Strings(closestMatches)
	}

	return closestMatches[0]
}

// Remove removes the given string if possible, returns a boolean true or false depending on success
func (dict *Dictionary) Remove(item string) bool {
	var newEntries []string
	found := false
	for _, entry := range dict.Entries {
		if entry == item {
			found = true
		} else {
			newEntries = append(newEntries, entry)
		}
	}

	if !found {
		return false
	}

	dict.Entries = newEntries
	return true
}

// Sort calls the sort strings function on the dictionary's entries
func (dict *Dictionary) Sort() {
	sort.Strings(dict.Entries)
}

// ToSlice converts the dictionary into an alphabetical list with the number of entries at the top
func (dict *Dictionary) ToSlice() []string {
	dict.Sort()
	return dict.Entries
}

/**
 * Interpretation of the full matrix approach to levenshtein distance calculation found at
 * https://en.wikipedia.org/wiki/Levenshtein_distance
 */
func levenshteinDistance(stringA string, stringB string) int {
	lenA := len(stringA)
	lenB := len(stringB)

	// Going from left to right for comparison in a double nested loop, or 2d matrix. Additionally, treat strings as
	// starting at 1 for matrix math to make sense
	row := make([]int, lenA+1)

	// Because the strings start at 1 in this version of levenshtein, the first row and first column values will be
	// the distance they are from the first element. e.g. [0, 0] will be 0, [1, 0] and [0, 1] will be 1
	for i := 1; i <= lenA; i++ {
		row[i] = i
	}

	// for loop variables represent row/column, outer loop (index name b for string b) increments column and
	// inner loop (index name a for string a) increments row. Because the value we want is the last row value in the
	// last column, don't need to keep track of previous columns
	for b := 1; b <= lenB; b++ {
		// initial value for first item in column will be distance into loop represented by [0,b], preset 'previous'
		previous := b
		for a := 1; a < lenA; a++ {
			current := row[a-1]

			// Find minimum of insertion, deletion, or substitution, then set that as the minimum for this column value
			// if the character at each index are different
			if stringB[b-1] != stringA[a-1] {
				current = min(min(row[a-1]+1, previous+1), row[a]+1)
			}

			row[a-1] = previous
			previous = current
		}
		row[lenA] = previous
	}

	return row[lenA]
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

package generator

import (
	"math/rand"
	"strings"
	"time"
)

func GenerateNickname() (res string) {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))

	vowels := []string{"a", "e", "i", "o", "u", "y"}

	initialConsonants := []string{
		"b", "c", "d", "f", "g", "h", "j", "k", "l", "m", "n", "p", "qu", "r",
		"s", "t", "v", "w", "x", "z", "ch", "cr", "dr", "fr", "gr", "kl", "ph",
		"pr", "st", "tr", "sh", "th", "wh",
	}

	midConsonants := []string{
		"b", "c", "d", "f", "g", "h", "j", "k", "l", "m", "n", "p", "r",
		"s", "t", "v", "z", "br", "cr", "dr", "gr", "kr", "pr", "tr",
		"ss", "tt", "ll", "mm", "nn", "th", "sh", "ph", "st", "sk", "sp",
	}

	finalConsonants := []string{
		"b", "c", "d", "f", "g", "h", "k", "l", "m", "n", "p", "r",
		"s", "t", "x", "z", "sh", "th", "ph", "ng", "ch", "gh", "ck",
	}

	suffixes := []string{
		"a", "o", "ia", "io", "is", "os", "us", "um", "ix", "ox", "ax",
		"an", "en", "in", "on", "ar", "er", "ir", "or", "ur", "az", "ez",
		"iz", "oz", "uz",
	}

	generatePseudoroot := func(minLength, maxLength int) string {
		rootLength := minLength
		if maxLength > minLength {
			rootLength = minLength + r.Intn(maxLength-minLength+1)
		}

		var root strings.Builder

		if r.Float64() > 0.2 {
			root.WriteString(initialConsonants[r.Intn(len(initialConsonants))])
		} else {
			root.WriteString(vowels[r.Intn(len(vowels))])
		}

		for root.Len() < rootLength {
			lastChar := root.String()[root.Len()-1]
			isVowel := false

			for _, v := range vowels {
				if string(lastChar) == v {
					isVowel = true
					break
				}
			}

			if isVowel {
				root.WriteString(midConsonants[r.Intn(len(midConsonants))])
			} else {
				root.WriteString(vowels[r.Intn(len(vowels))])
			}
		}

		if root.Len() > rootLength {
			return root.String()[:rootLength]
		}

		lastChar := root.String()[root.Len()-1]
		isVowel := false

		for _, v := range vowels {
			if string(lastChar) == v {
				isVowel = true
				break
			}
		}

		if !isVowel && r.Float64() > 0.5 {
			root.WriteString(vowels[r.Intn(len(vowels))])
		}

		return root.String()
	}

	length := 10

	strategy := r.Intn(4) + 1

	nickname := ""

	switch strategy {
	case 1:
		rootLength := length - 2
		nickname = generatePseudoroot(3, rootLength)

		if len(nickname) <= length-2 {
			var suitableSuffixes []string
			for _, suffix := range suffixes {
				if len(suffix) <= length-len(nickname) {
					suitableSuffixes = append(suitableSuffixes, suffix)
				}
			}

			if len(suitableSuffixes) > 0 {
				nickname += suitableSuffixes[r.Intn(len(suitableSuffixes))]
			}
		}

	case 2:
		firstPartLength := 3 + r.Intn(length-5)
		firstPart := generatePseudoroot(3, firstPartLength)

		secondPartLength := length - len(firstPart)
		if secondPartLength >= 2 {
			secondPart := generatePseudoroot(2, secondPartLength)
			nickname = firstPart + secondPart
		} else {
			nickname = firstPart
			letters := "abcdefghijklmnopqrstuvwxyz"
			for i := 0; i < secondPartLength; i++ {
				nickname += string(letters[r.Intn(len(letters))])
			}
		}

	case 3:
		rootLength := length - 1
		nickname = generatePseudoroot(3, rootLength)

		if len(nickname) < length {
			if r.Float64() > 0.6 {
				nickname += string('0' + rune(r.Intn(10)))
			} else {
				letters := "abcdefghijklmnopqrstuvwxyz"
				nickname += string(letters[r.Intn(len(letters))])
			}
		}

	case 4:
		nickname = generatePseudoroot(length, length)
	}

	if len(nickname) > length {
		nickname = nickname[:length]
	}

	for len(nickname) < length {
		lastChar := nickname[len(nickname)-1]
		isVowel := false

		for _, v := range vowels {
			if string(lastChar) == v {
				isVowel = true
				break
			}
		}

		if isVowel {
			nickname += finalConsonants[r.Intn(len(finalConsonants))]
		} else {
			nickname += vowels[r.Intn(len(vowels))]
		}

		if len(nickname) > length {
			nickname = nickname[:length]
		}
	}

	return strings.ToLower(nickname)
}

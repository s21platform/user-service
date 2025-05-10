package generator

import (
	"math/rand"
	"strings"
	"time"
	"unicode"

	"github.com/brianvoe/gofakeit/v6"
)

// набор согласных для создания более читаемых комбинаций
var consonants = []string{"b", "c", "d", "f", "g", "h", "j", "k", "l", "m", "n", "p", "q", "r", "s", "t", "v", "w", "x", "z"}

// набор гласных для создания более читаемых комбинаций
var vowels = []string{"a", "e", "i", "o", "u", "y"}

// интересные краткие слова для генерации
var interestingWords = []string{
	"seed", "moon", "star", "sky", "rain", "sun", "air", "fire", "ice",
	"flow", "wave", "wind", "leaf", "tree", "rock", "cave", "hill", "sand",
	"dawn", "dusk", "echo", "cake", "soda", "brew", "milk", "wine", "mint",
	"sage", "rose", "lily", "pine", "oak", "time", "bell", "song", "tune",
}

// популярные имена разных культур
var popularNames = []string{
	"luna", "alex", "mila", "dora", "zora", "tara", "lena", "maya", "elsa",
	"nora", "lisa", "emma", "yara", "aria", "mira", "sara", "cora", "lola",
	"kira", "nova", "ava", "nia", "geo", "leo", "kai", "mio", "neo", "ari",
}

// генерирует один случайный слог
func generateSyllable() string {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))

	// чередование согласной и гласной
	return consonants[r.Intn(len(consonants))] + vowels[r.Intn(len(vowels))]
}

// генерирует простое читаемое слово
func generateReadableWord(minLength, maxLength int) string {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))

	// определяем случайную длину слова в заданных пределах
	wordLength := minLength
	if maxLength > minLength {
		wordLength = minLength + r.Intn(maxLength-minLength+1)
	}

	var word strings.Builder

	// начинаем с согласной или гласной с равной вероятностью
	if r.Intn(2) == 0 {
		word.WriteString(consonants[r.Intn(len(consonants))])
	} else {
		word.WriteString(vowels[r.Intn(len(vowels))])
	}

	// добавляем слоги, чередуя гласные и согласные
	for word.Len() < wordLength {
		lastChar := rune(word.String()[word.Len()-1])
		if isVowel(lastChar) {
			word.WriteString(consonants[r.Intn(len(consonants))])
		} else {
			word.WriteString(vowels[r.Intn(len(vowels))])
		}
	}

	// обрезаем слово до нужной длины
	if word.Len() > wordLength {
		return word.String()[:wordLength]
	}

	return word.String()
}

// проверяет, является ли символ гласной
func isVowel(r rune) bool {
	r = unicode.ToLower(r)
	for _, v := range vowels {
		if string(r) == v {
			return true
		}
	}
	return false
}

// GenerateNickname создает читаемый никнейм длиной ровно 10 символов
func GenerateNickname() string {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	gofakeit.Seed(0)

	// выбираем стратегию генерации
	strategy := r.Intn(5)

	var nickname string

	switch strategy {
	case 0:
		// короткое имя + читаемое слово
		name := strings.ToLower(popularNames[r.Intn(len(popularNames))])
		remainingLength := 10 - len(name)
		if remainingLength > 0 {
			secondPart := generateReadableWord(remainingLength, remainingLength)
			nickname = name + secondPart
		} else {
			nickname = name[:10]
		}

	case 1:
		// известное короткое слово + слоги
		word := interestingWords[r.Intn(len(interestingWords))]
		remainingLength := 10 - len(word)
		suffix := ""
		for i := 0; i < remainingLength; i++ {
			if i%2 == 0 {
				suffix += consonants[r.Intn(len(consonants))]
			} else {
				suffix += vowels[r.Intn(len(vowels))]
			}
		}
		nickname = word + suffix

	case 2:
		// случайное читаемое слово + одна буква или цифра
		mainWord := generateReadableWord(8, 9)
		suffix := ""
		if len(mainWord) < 10 {
			remainingLength := 10 - len(mainWord)
			for i := 0; i < remainingLength; i++ {
				// добавляем одну букву или цифру
				suffix += string('a' + rune(r.Intn(26)))
			}
		}
		nickname = mainWord + suffix

	case 3:
		// два коротких читаемых слова
		firstWord := generateReadableWord(4, 6)
		remainingLength := 10 - len(firstWord)
		secondWord := generateReadableWord(remainingLength, remainingLength)
		nickname = firstWord + secondWord

	case 4:
		// полностью случайное но читаемое слово
		nickname = generateReadableWord(10, 10)
	}

	// обеспечиваем ровно 10 символов
	if len(nickname) > 10 {
		nickname = nickname[:10]
	} else if len(nickname) < 10 {
		// добавляем случайные буквы, если не хватает
		for len(nickname) < 10 {
			nickname += consonants[r.Intn(len(consonants))]
		}
	}

	return strings.ToLower(nickname)
}

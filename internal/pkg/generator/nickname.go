package generator

import (
	"strings"

	"github.com/brianvoe/gofakeit/v6"
)

func GenerateNickname() (res string) {
	defer func() {
		res = strings.ToLower(res)
	}()
	gofakeit.Seed(0)

	name := gofakeit.FirstName()
	if len(name) >= 10 {
		res = name[:10]
		return
	}

	ost := 10 - len(name)
	surname := gofakeit.LastName()
	if len(surname) >= ost {
		res = name + surname[:ost]
		return
	}

	ost -= len(surname)
	color := gofakeit.Color()
	if len(color) >= ost {
		res = name + surname + color[:ost]
		return
	}

	dinner := gofakeit.Dinner()
	res = (name + surname + color + dinner)[0:10]
	return
}

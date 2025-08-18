package generator

import (
	"strings"

	"github.com/brianvoe/gofakeit/v6"
)

const l = 10

func GenerateNickname() (res string) {
	defer func() {
		res = strings.ToLower(res)
	}()
	gofakeit.Seed(0)
	name := gofakeit.FirstName()
	if len(name) >= l {
		res = name[:l]
		return
	}
	ost := l - len(name)
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
	res = (name + surname + color + dinner)[0:l]
	return
}

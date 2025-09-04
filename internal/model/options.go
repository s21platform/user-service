package model

type Attribute int

func (a Attribute) Int64() int64 {
	return int64(a)
}

const (
	AttributeNickname_1 = Attribute(1)
	AttributeName_2     = Attribute(2)
	AttributeSurname_3  = Attribute(3)
	AttributeBirthday_4 = Attribute(4)
	AttributeCity_5     = Attribute(5)
	AttributeTelegram_6 = Attribute(6)
)

var PersonalityForm = []Attribute{
	AttributeName_2,
	AttributeSurname_3,
	AttributeBirthday_4,
}

type AttributeMetaMap map[int64]AttributeMeta

type AttributeMeta struct {
	Label string
	Type  string
}

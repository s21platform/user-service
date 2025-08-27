package rest

import (
	api "github.com/s21platform/user-service/internal/generated"
	"github.com/s21platform/user-service/internal/model"
	"github.com/samber/lo"
)

// mapPersonalityToProfileItems преобразует данные о личности в профильные элементы
func mapPersonalityToProfileItems(personality model.Personality, options model.AttributeMetaMap) []api.ProfileItem {
	items := make([]api.ProfileItem, 0)

	// Name
	if meta, ok := options[model.AttributeName_2.Int64()]; ok && personality.Name != nil {
		items = append(items, api.ProfileItem{
			Title: meta.Label,
			Type:  meta.Type,
			Value: personality.Name,
		})
	}

	// Surname
	if meta, ok := options[model.AttributeSurname_3.Int64()]; ok && personality.Surname != nil {
		items = append(items, api.ProfileItem{
			Title: meta.Label,
			Type:  meta.Type,
			Value: personality.Surname,
		})
	}

	// Birthday
	if meta, ok := options[model.AttributeBirthday_4.Int64()]; ok && personality.Birthdate != nil {
		birthdateStr := lo.ToPtr(personality.Birthdate.String())
		items = append(items, api.ProfileItem{
			Title: meta.Label,
			Type:  meta.Type,
			Value: birthdateStr,
		})
	}

	return items
}

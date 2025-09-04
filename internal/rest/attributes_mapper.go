package rest

import (
	openapi_types "github.com/oapi-codegen/runtime/types"
	api "github.com/s21platform/user-service/internal/generated"
	"github.com/s21platform/user-service/internal/model"
)

// mapUserAttributesToAttributeItems преобразует данные атрибутов пользователя в AttributeItem
// Всегда возвращает структуру для каждого запрошенного атрибута, даже если значение отсутствует
func mapUserAttributesToAttributeItems(userAttributes model.UserAttributes, options model.AttributeMetaMap, attributeIds []int64) []api.AttributeItem {
	var result []api.AttributeItem

	for _, attrId := range attributeIds {
		meta, exists := options[attrId]
		if !exists {
			continue
		}

		item := api.AttributeItem{
			AttributeId: attrId,
			Title:       meta.Label,
			Type:        meta.Type,
		}

		// Маппим значения в зависимости от типа атрибута
		// Всегда создаем структуру, даже если значение nil
		switch model.Attribute(attrId) {
		case model.AttributeName_2:
			item.ValueString = userAttributes.Name
		case model.AttributeSurname_3:
			item.ValueString = userAttributes.Surname
		case model.AttributeBirthday_4:
			if userAttributes.Birthdate != nil {
				// Преобразуем time.Time в openapi_types.Date
				date := openapi_types.Date{Time: *userAttributes.Birthdate}
				item.ValueDate = &date
			}
		case model.AttributeCity_5:
			item.ValueInt = userAttributes.CityId
		case model.AttributeTelegram_6:
			item.ValueString = userAttributes.Telegram
		}

		result = append(result, item)
	}

	return result
}

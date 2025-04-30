package model

import (
	model "github.com/s21platform/optionhub-lib"
)

const (
	Attribute_name      = model.Attribute_Name
	Attribute_OS        = model.Attribute_OS
	Attribute_Birthdate = model.Attribute_Birthdate
	Attribute_Telegram  = model.Attribute_Telegram
	Attribute_Git       = model.Attribute_Git
)

var AttributeSetters = map[int64]func(attr model.AttributeValue, model *ProfileData){
	Attribute_name: func(attr model.AttributeValue, model *ProfileData) {
		if attr.ValueString != nil {
			model.Name = *attr.ValueString
		}
	},
	Attribute_OS: func(attr model.AttributeValue, model *ProfileData) {
		if attr.ValueInt != nil {
			model.OsId = *attr.ValueInt
		}
	},
	Attribute_Birthdate: func(attr model.AttributeValue, model *ProfileData) {
		if attr.ValueDate != nil {
			model.Birthdate = attr.ValueDate
		}
	},
	Attribute_Telegram: func(attr model.AttributeValue, model *ProfileData) {
		if attr.ValueString != nil {
			model.Telegram = *attr.ValueString
		}
	},
	Attribute_Git: func(attr model.AttributeValue, model *ProfileData) {
		if attr.ValueString != nil {
			model.Git = *attr.ValueString
		}
	},
}

package functions

import (
	"unicode/utf8"

	"github.com/asaskevich/govalidator"
	"go.uber.org/zap"
)

// :""`
// " valid:"tenderDescription"`
// " valid:"tenderType"`
// id:"tenderStatus"`
// nId" valid:"organizationId"`
// name" valid:"username"`

func Validate(data any) (bool, error) {
	return govalidator.ValidateStruct(data)
}

func InitDtoValidator(logger *zap.Logger) {
	govalidator.SetFieldsRequiredByDefault(true)

	govalidator.TagMap["name"] = func(name string) bool {
		lenName := utf8.RuneCountInString(name)
		return lenName > 0 && lenName <= 100
	}

	govalidator.TagMap["description"] = func(description string) bool {
		lenDescription := utf8.RuneCountInString(description)
		return lenDescription > 0 && lenDescription <= 500
	}

	govalidator.TagMap["username"] = func(username string) bool {
		usernameLen := utf8.RuneCountInString(username)
		return usernameLen > 0 && usernameLen <= 50
	}

	govalidator.TagMap["firstName"] = func(firstName string) bool {
		firstNameLen := utf8.RuneCountInString(firstName)
		return firstNameLen > 0 && firstNameLen <= 50
	}

	govalidator.TagMap["lastName"] = func(lastName string) bool {
		lastNameLen := utf8.RuneCountInString(lastName)
		return lastNameLen > 0 && lastNameLen <= 50
	}
	logger.Info("Custom tags created")
}

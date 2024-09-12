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

	govalidator.TagMap["serviceType"] = func(serviceType string) bool {
		lenTenderType := utf8.RuneCountInString(serviceType)
		return lenTenderType > 0 && lenTenderType <= 12
	}

	govalidator.TagMap["status"] = func(status string) bool {
		statusLen := utf8.RuneCountInString(status)
		return statusLen > 0 && statusLen <= 9
	}

	govalidator.TagMap["organizationId"] = func(organizationId string) bool {
		organizationIdLen := utf8.RuneCountInString(organizationId)
		return organizationIdLen > 0
	}

	govalidator.TagMap["username"] = func(username string) bool {
		usernameLen := utf8.RuneCountInString(username)
		return usernameLen > 0 && usernameLen <= 50
	}

	govalidator.TagMap["tenderId"] = func(tenderId string) bool {
		lenTenderId := utf8.RuneCountInString(tenderId)
		return lenTenderId > 0
	}

	logger.Info("Custom tags created")
}

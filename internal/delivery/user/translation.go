package user

import (
	"tender-workspace/internal/entity"
	"tender-workspace/internal/entity/dto"
	"tender-workspace/internal/utils/functions"
)

func newUserOutput(uData *entity.Employee) *dto.UserOutPut {
	return &dto.UserOutPut{
		ID:        uData.ID,
		Username:  uData.Username,
		FirstName: uData.FirstName,
		LastName:  uData.LastName,
		CreatedAt: functions.FormatTime(uData.CreatedAt),
	}
}

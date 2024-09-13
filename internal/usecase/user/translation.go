package user

import (
	ent "tender-workspace/internal/entity"
	"tender-workspace/internal/entity/dto"
)

func newUser(data *dto.UserInput) *ent.Employee {
	return &ent.Employee{
		Username:  data.Username,
		FirstName: data.FirstName,
		LastName:  data.LastName,
	}
}

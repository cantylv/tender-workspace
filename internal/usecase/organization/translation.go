package organization

import (
	ent "tender-workspace/internal/entity"
	"tender-workspace/internal/entity/dto"
)

func newOrganization(data *dto.OrganizationInput) *ent.Organization {
	return &ent.Organization{
		Name:        data.Name,
		Description: data.Description,
		Type:        data.Type,
	}
}

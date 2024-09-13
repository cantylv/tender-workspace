package organization

import (
	ent "tender-workspace/internal/entity"
	"tender-workspace/internal/entity/dto"
	f "tender-workspace/internal/utils/functions"
)

func newArrayOrgOutput(orgs []*ent.Organization) []*dto.OrganizationOutput {
	res := make([]*dto.OrganizationOutput, 0, len(orgs))
	for _, org := range orgs {
		orgOutput := newOrganizationOutput(org)
		res = append(res, orgOutput)
	}
	return res
}

func newOrganizationOutput(org *ent.Organization) *dto.OrganizationOutput {
	return &dto.OrganizationOutput{
		ID:          org.ID,
		Name:        org.Name,
		Description: org.Description,
		Type:        org.Type,
		CreatedAt:   f.FormatTime(org.CreatedAt),
	}
}

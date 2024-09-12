package tender

import (
	ent "tender-workspace/internal/entity"
	"tender-workspace/internal/entity/dto"
	tqp "tender-workspace/internal/entity/dto/queries/tenders"
	t "tender-workspace/internal/repo/tender"
)

func newTender(user *ent.Employee, tenderInput *dto.TenderInput) *ent.Tender {
	return &ent.Tender{
		Name:           tenderInput.Description,
		Description:    tenderInput.Description,
		Type:           tenderInput.Type,
		Status:         tenderInput.Status,
		Version:        1,
		OrganizationID: tenderInput.OrganizationID,
		CreatorID:      user.ID,
	}
}

func newOrganizationResponsible(organizationId, userId int) *ent.OrganizationResponsible {
	return &ent.OrganizationResponsible{
		OrganizationID: organizationId,
		UserID:         userId,
	}
}

func newUserTenderProps(params *tqp.ListUserTenders, user *ent.Employee) *t.UserTendersProps {
	return &t.UserTendersProps{
		Limit:  params.Limit,
		Offset: params.Offset,
		UserID: user.ID,
	}
}

func newUpdateTenderData(updateData *dto.TenderUpdateDataInput) *ent.UpdateTenderData {
	return &ent.UpdateTenderData{
		Name:        updateData.Name,
		Description: updateData.Description,
		Type:        updateData.ServiceType,
	}
}

func newUpdateTenderProps(updateProps *tqp.TenderUpdate, user *ent.Employee) *t.UpdateTenderProps {
	return &t.UpdateTenderProps{
		TenderID: updateProps.TenderID,
		UserID:   user.ID,
	}
}

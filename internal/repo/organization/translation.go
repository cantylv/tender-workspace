package organization

import (
	"database/sql"
	ent "tender-workspace/internal/entity"
	"time"
)

type organizationDB struct {
	ID          int
	Name        string
	Description sql.NullString
	Type        string
	CreatedAt   time.Time
}

func getArrayOrganizationFromDB(rows []*organizationDB) []*ent.Organization {
	orgs := make([]*ent.Organization, 0, len(rows))
	for _, row := range rows {
		orgs = append(orgs, getOrganizationFromDB(row))
	}
	return orgs
}

func getOrganizationFromDB(row *organizationDB) *ent.Organization {
	orgDescription := ""
	if row.Description.Valid {
		orgDescription = row.Description.String
	}
	return &ent.Organization{
		ID:          row.ID,
		Name:        row.Name,
		Description: orgDescription,
		Type:        row.Type,
		CreatedAt:   row.CreatedAt,
	}
}

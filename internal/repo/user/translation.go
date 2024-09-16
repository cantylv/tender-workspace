package user

import (
	"database/sql"
	ent "tender-workspace/internal/entity"
	"time"
)

type userDB struct {
	ID        int
	Username  string
	FirstName sql.NullString
	LastName  sql.NullString
	CreatedAt time.Time
}

func getArrayUserFromDB(rows []*userDB) []*ent.Employee {
	users := make([]*ent.Employee, 0, len(rows))
	for _, row := range rows {
		users = append(users, getUserFromDB(row))
	}
	return users
}

func getUserFromDB(row *userDB) *ent.Employee {
	userFirstName := ""
	if row.FirstName.Valid {
		userFirstName = row.FirstName.String
	}
	userLastName := ""
	if row.LastName.Valid {
		userLastName = row.LastName.String
	}
	return &ent.Employee{
		ID:        row.ID,
		Username:  row.Username,
		FirstName: userFirstName,
		LastName:  userLastName,
		CreatedAt: row.CreatedAt,
	}
}

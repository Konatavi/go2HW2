package store

import (
	"fmt"
	"log"

	"github.com/Konatavi/go2HW2/internal/app/models"
)

type UsersautoRepository struct {
	store *Store
}

var (
	tableUser string = "usersauto"
)

//Create user in database
func (ur *UsersautoRepository) Create(u *models.Usersauto) (*models.Usersauto, error) {
	query := fmt.Sprintf("INSERT INTO %s (username, password) VALUES ($1, $2) RETURNING id", tableUser)
	if err := ur.store.db.QueryRow(
		query,
		u.Username,
		u.Password,
	).Scan(&u.ID); err != nil {
		return nil, err
	}
	return u, nil
}

//Find by Username
func (ur *UsersautoRepository) FindByUsername(username string) (*models.Usersauto, bool, error) {
	usersauto, err := ur.SelectAll()
	var founded bool
	if err != nil {
		return nil, founded, err
	}
	var usersautoFinded *models.Usersauto
	for _, u := range usersauto {
		if u.Username == username {
			usersautoFinded = u
			founded = true
			break
		}
	}
	return usersautoFinded, founded, nil
}

//Select All
func (ur *UsersautoRepository) SelectAll() ([]*models.Usersauto, error) {
	query := fmt.Sprintf("SELECT * FROM %s", tableUser)
	rows, err := ur.store.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	usersauto := make([]*models.Usersauto, 0)
	for rows.Next() {
		u := models.Usersauto{}
		err := rows.Scan(&u.ID, &u.Username, &u.Password)
		if err != nil {
			log.Println(err)
			continue
		}
		usersauto = append(usersauto, &u)
	}
	return usersauto, nil

}

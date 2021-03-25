package store

import (
	"fmt"
	"log"

	"github.com/Konatavi/go2HW2/internal/app/models"
)

type AutomobilesRepository struct {
	store *Store
}

var (
	tableAutomobiles string = "automobiles"
)

//For Post request
func (ar *AutomobilesRepository) Create(a *models.Automobiles) (*models.Automobiles, error) {
	query := fmt.Sprintf("INSERT INTO %s (mark, maxspeed, distance, handler, stock) VALUES ($1, $2, $3, $4, $5) RETURNING id", tableAutomobiles)
	if err := ar.store.db.QueryRow(query, a.Mark, a.Maxspeed, a.Distance, a.Handler, a.Stock).Scan(&a.ID); err != nil {
		return nil, err
	}
	return a, nil
}

//For DELETE request
func (ar *AutomobilesRepository) DeleteByMark(mark string) (*models.Automobiles, error) {
	article, ok, err := ar.FindAutomobileByMark(mark)
	if err != nil {
		return nil, err
	}
	if ok {
		query := fmt.Sprintf("delete from %s where mark=$1", tableAutomobiles)
		_, err = ar.store.db.Exec(query, mark)
		if err != nil {
			return nil, err
		}
	}

	return article, nil
}

//Helper for find by mask and GET request
func (ar *AutomobilesRepository) FindAutomobileByMark(mark string) (*models.Automobiles, bool, error) {
	automobiles, err := ar.SelectAll()
	founded := false
	if err != nil {
		return nil, founded, err
	}
	var automobileFinded *models.Automobiles
	for _, a := range automobiles {
		if a.Mark == mark {
			automobileFinded = a
			founded = true
		}
	}

	return automobileFinded, founded, nil

}

//Get all request and helper for FindAutomobileByMark
func (ar *AutomobilesRepository) SelectAll() ([]*models.Automobiles, error) {
	query := fmt.Sprintf("SELECT * FROM %s", tableAutomobiles)
	rows, err := ar.store.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	automobiles := make([]*models.Automobiles, 0)
	for rows.Next() {
		a := models.Automobiles{}
		err := rows.Scan(&a.ID, &a.Mark, &a.Maxspeed, &a.Distance, &a.Handler, &a.Stock)
		if err != nil {
			log.Println(err)
			continue
		}
		automobiles = append(automobiles, &a)
	}
	return automobiles, nil
}

//For UPDATE request
func (ar *AutomobilesRepository) UpdateByMark(mark string, newAuto *models.Automobiles) (*models.Automobiles, error) {
	oldAuto, ok, err := ar.FindAutomobileByMark(mark)
	if err != nil {
		return nil, err
	}
	if ok {
		query := fmt.Sprintf("update %s set maxspeed = $1, distance = $2, handler = $3, stock = $4 where mark=$5", tableAutomobiles)
		_, err = ar.store.db.Exec(query, newAuto.Maxspeed, newAuto.Distance, newAuto.Handler, newAuto.Stock, mark)
		if err != nil {
			return nil, err
		}
	}
	newAuto.ID = oldAuto.ID

	return newAuto, nil
}

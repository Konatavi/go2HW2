package store

import (
	"database/sql"
	"log"

	_ "github.com/lib/pq"
)

//Instance of store
type Store struct {
	config                *Config
	db                    *sql.DB
	usersautoRepository   *UsersautoRepository
	automobilesRepository *AutomobilesRepository
}

// Constructor for store
func New(config *Config) *Store {
	return &Store{
		config: config,
	}
}

//Open store method
func (s *Store) Open() error {
	db, err := sql.Open("postgres", s.config.DatabaseURL)
	if err != nil {
		return err
	}
	//Проверим, что все ок. Реально соединение тут не создается. Соединение только при первом вызове
	//db.Ping() // Пустой SELECT *
	if err := db.Ping(); err != nil {
		return err
	}
	s.db = db
	log.Println("Connection to db successfully")
	return nil
}

//Close store method
func (s *Store) Close() {
	s.db.Close()
}

//Public for UsersautoRepositoryRepo
func (s *Store) Usersauto() *UsersautoRepository {
	if s.usersautoRepository != nil {
		return s.usersautoRepository
	}
	s.usersautoRepository = &UsersautoRepository{
		store: s,
	}
	return s.usersautoRepository
}

//Public for ArticleRepo
func (s *Store) Automobiles() *AutomobilesRepository {
	if s.automobilesRepository != nil {
		return s.automobilesRepository
	}
	s.automobilesRepository = &AutomobilesRepository{
		store: s,
	}
	return s.automobilesRepository
}

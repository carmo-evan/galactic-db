package db

import (
	"encoding/json"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"log"
)

type Armament struct {
	Title string `json:"title" db:"Title"`
	SpaceshipId int `json:"-" db:"SpaceshipId"`
	Qty int `json:"qty" db:"Qty"`
}

type SpaceshipInfo struct {
	Id int `json:"id" db:"Id"`
	Name string `json:"name" db:"Name"`
	Status string `json:"status" db:"Status"`
}

type SpaceshipDTO struct {
	Id int `json:"id" db:"Id"`
	Name string `json:"name" db:"Name"`
	Status string `json:"status" db:"Status"`
	Class string `json:"class" db:"Class"`
	Image string `json:"image" db:"Image"`
	Crew int `json:"crew" db:"Crew"`
	Value int `json:"value" db:"Value""`
	ArmamentJSON string `json:"-" db:"Armaments"`
}

type Spaceship struct {
	Id int `json:"id" db:"Id"`
	Name string `json:"name" db:"Name"`
	Status string `json:"status" db:"Status"`
	Class string `json:"class" db:"Class"`
	Image string `json:"image" db:"Image"`
	Crew int `json:"crew" db:"Crew"`
	Value int `json:"value" db:"Value""`
	Armament []Armament `json:armament`
}

type SpaceshipFilters struct {
	Name string `json:name`
	Class string `json:class`
	Status string `json:status`
}

type Store interface {
	GetSpaceships(limit int, offset int, filters SpaceshipFilters) ([]SpaceshipInfo, error)
	InsertSpaceship(spaceship SpaceshipDTO) (int, error)
	UpdateSpaceship(spaceship SpaceshipDTO) error
	DeleteSpaceship(id int) error
	GetSpaceship(id int) (SpaceshipDTO, error)
}

type SQLStore struct {
	db *sqlx.DB
}

func CreateSpaceshipDTO(spaceship Spaceship) (SpaceshipDTO, error) {
	dto := SpaceshipDTO{
		Id:           spaceship.Id,
		Name:         spaceship.Name,
		Status:       spaceship.Status,
		Class:        spaceship.Class,
		Image:        spaceship.Image,
		Crew:         spaceship.Crew,
		Value:        spaceship.Value,
		ArmamentJSON: "",
	}
	armamentsJson, err := json.Marshal(spaceship.Armament)
	if err != nil {
		return dto, err
	}
	dto.ArmamentJSON = string(armamentsJson)
	return dto, nil
}

func ParseSpaceshipDTO(dto SpaceshipDTO) (Spaceship, error) {
	sp := Spaceship{
		Id:       dto.Id,
		Name:     dto.Name,
		Status:   dto.Status,
		Class:    dto.Class,
		Image:    dto.Image,
		Crew:     dto.Crew,
		Value:    dto.Value,
		Armament: nil,
	}
	var ar []Armament
	json.Unmarshal([]byte(dto.ArmamentJSON), &ar)
	sp.Armament = ar
	return sp, nil
}

func (s *SQLStore) InsertSpaceship(spaceship SpaceshipDTO) (int, error) {
	res, err := s.db.Exec("INSERT INTO spaceships (Name, Class, Image, Status, Crew, Value, Armaments) VALUES (?, ?, ?, ?, ?, ?, ?)",
		spaceship.Name, spaceship.Class, spaceship.Image, spaceship.Status, spaceship.Crew, spaceship.Value, spaceship.ArmamentJSON,
	)
	if err != nil {
		return 0, err
	}
	id, err := res.LastInsertId()
	return int(id), err
}

func (s *SQLStore) InsertArmament(armament Armament) error {
	_, err := s.db.Exec("INSERT INTO armaments (Title, SpaceshipId, Qty) VALUES (?, ?, ?)",
		armament.Title, armament.SpaceshipId, armament.Qty,
	)
	return err
}

func (s *SQLStore) UpdateSpaceship(spaceship SpaceshipDTO) error {
	_, err := s.db.Exec("UPDATE spaceships SET Name = ?, Class = ?, Image = ?, Status = ?, Crew = ?, Value = ?, Armaments = ? WHERE Id = ?",
		spaceship.Name, spaceship.Class, spaceship.Image, spaceship.Status, spaceship.Crew, spaceship.Value, spaceship.ArmamentJSON, spaceship.Id,
	)
	return err
}


func (s *SQLStore) DeleteSpaceship(id int) error {
	_, err := s.db.Exec("DELETE FROM spaceships WHERE Id = ?", id)
	return err
}


func (s *SQLStore) GetSpaceships(limit int, offset int, filters SpaceshipFilters) ([]SpaceshipInfo, error) {
	log.Println(limit, offset, filters)
	spaceships := []SpaceshipInfo{}
	err := s.db.Select(&spaceships, "SELECT Id, Name, Status FROM spaceships LIMIT ? OFFSET ?", limit, offset)
	log.Println("TODO: Implement filters", filters)
	return spaceships, err
}

func (s *SQLStore) GetSpaceship(id int) (SpaceshipDTO, error) {
	spaceship := SpaceshipDTO{}
	log.Println(id)
	err := s.db.Get(&spaceship, "SELECT * FROM spaceships WHERE Id = ? LIMIT 1", id)
	return spaceship, err
}

func GetSQLStore() Store {
	dsn := "root:example@/galactic"
	db, err := sqlx.Connect("mysql", dsn)
	if err != nil {
		log.Fatalln(err)
	}
	s := SQLStore{db: db}
	return &s
}

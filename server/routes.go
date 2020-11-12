package server

import (
	"alchemy/galacticdb/db"
	"encoding/json"
	"fmt"
	"github.com/go-chi/chi"
	"github.com/go-chi/render"
	"log"
	"net/http"
	"strconv"
)

func getSpaceshipRoute(s db.Store) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		id, err := getSpaceshipId(w, r)
		if err != nil {
			log.Println(id)
			render.Render(w, r, ErrorRenderer(fmt.Errorf("id is required")))
			return
		}
		sp, err := s.GetSpaceship(id)
		if err != nil {
			log.Println(err)
			render.Render(w, r, ErrorRenderer(fmt.Errorf("could not find spaceship")))
			return
		}
		armaments, err := s.GetArmaments(id)
		if err != nil {
			log.Println(err)
		}
		spaceship := db.Spaceship{SpaceshipDTO: sp, Armament: armaments}
		render.JSON(w, r, spaceship)
	}
}

func deleteSpaceshipRoute(s db.Store) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		id, err := getSpaceshipId(w, r)
		if err != nil {
			log.Println(id)
			render.Render(w, r, ErrorRenderer(fmt.Errorf("id is required")))
			return
		}
		err = s.DeleteArmaments(id)
		if err != nil {
			log.Println(err)
			render.Render(w, r, ErrorRenderer(fmt.Errorf("could not delete armaments")))
			return
		}
		err = s.DeleteSpaceship(id)
		if err != nil {
			log.Println(err)
			render.Render(w, r, ErrorRenderer(fmt.Errorf("could not delete spaceship")))
			return
		}
		response := SuccessResponse{Success: true}
		render.JSON(w, r, response)
	}
}


func getSpaceshipId(w http.ResponseWriter, r *http.Request) (id int, err error){
	idParam := chi.URLParam(r, "id")
	if idParam == "" {
		err := fmt.Errorf("id is required")
		render.Render(w, r, ErrorRenderer(err))
		return 0, err
	}
	id, err = strconv.Atoi(idParam)
	if err != nil {
		render.Render(w, r, ErrorRenderer(fmt.Errorf("id must be a number")))
		return id, err
	}
	return id, nil
}



func createSpaceshipRoute(s db.Store) func (w http.ResponseWriter, r *http.Request) {
	return func (w http.ResponseWriter, r *http.Request) {
		var sp db.Spaceship
		err := json.NewDecoder(r.Body).Decode(&sp)
		if err != nil {
			render.Render(w, r, ErrorRenderer(fmt.Errorf("invalid body")))
			return
		}
		id, err := s.InsertSpaceship(sp.SpaceshipDTO)
		sp.Id = id
		if err != nil {
			log.Println(err)
			render.Render(w, r, ServerErrorRenderer(fmt.Errorf("could not save spaceships")))
			return
		}
		for _, ar := range sp.Armament {
			ar.SpaceshipId = sp.Id
			log.Println(ar)
			err = s.InsertArmament(ar)
			if err != nil {
				log.Println(err)
				render.Render(w, r, ServerErrorRenderer(fmt.Errorf("could not save armaments")))
				return
			}
		}

		response := SuccessResponse{
			Success: true,
		}
		render.JSON(w, r, response)
	}
}

func updateSpaceshipRoute(s db.Store) func (w http.ResponseWriter, r *http.Request) {
	return func (w http.ResponseWriter, r *http.Request) {
		var sp db.Spaceship
		err := json.NewDecoder(r.Body).Decode(&sp)
		if err != nil {
			render.Render(w, r, ErrorRenderer(fmt.Errorf("invalid body")))
			return
		}
		err = s.UpdateSpaceship(sp.SpaceshipDTO)
		if err != nil {
			log.Println(err)
			render.Render(w, r, ServerErrorRenderer(fmt.Errorf("could not save spaceships")))
			return
		}
		for _, ar := range sp.Armament {
			ar.SpaceshipId = sp.Id
			log.Println(ar)
			err = s.UpdateArmament(ar)
			if err != nil {
				log.Println(err)
				render.Render(w, r, ServerErrorRenderer(fmt.Errorf("could not save armaments")))
				return
			}
		}
		response := SuccessResponse{
			Success: true,
		}
		render.JSON(w, r, response)
	}
}

func getSpaceshipsRoute(s db.Store) func(w http.ResponseWriter, r *http.Request) {
 return func (w http.ResponseWriter, r *http.Request) {
	 pageParam := chi.URLParam(r, "page")
	 if pageParam == "" {
		 pageParam = "1"
	 }
	 page, err := strconv.Atoi(pageParam)
	 if err != nil {
		 log.Println(err)
		 render.Render(w, r, ErrorRenderer(fmt.Errorf("invalid page argument")))
		 return
	 }
	 filters := db.SpaceshipFilters{
		 Name:   "",
		 Class:  "",
		 Status: "",
	 }
	 log.Println(100, (page - 1) * 100, filters)
	 spaceships, err := s.GetSpaceships(100, (page - 1) * 100, filters)
	if err != nil {
		log.Println(err)
		render.Render(w, r, ServerErrorRenderer(fmt.Errorf("could not retrieve spaceships")))
		return
	}
	render.JSON(w, r, spaceships)
 }
}

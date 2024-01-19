package handler

import (
	"city2city/api/models"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
)

func (h Handler) City(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		h.CreateCity(w, r)
	case http.MethodGet:
		values := r.URL.Query()
		if _, ok := values["id"]; !ok {
			h.GetCityList(w, r)
		} else {
			h.GetCityByID(w, r)
		}
	case http.MethodPut:
		h.UpdateCity(w, r)
	case http.MethodDelete:
		h.DeleteCity(w, r)
	}
}

func (h Handler) CreateCity(w http.ResponseWriter, r *http.Request) {
	createdCity := models.CreateCity{}

	if err := json.NewDecoder(r.Body).Decode(&createdCity); err != nil {
		handleResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	id, err := h.storage.City().Create(createdCity)
	if err != nil {
		handleResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	city, err := h.storage.City().Get(id)
	if err != nil {
		handleResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	handleResponse(w, http.StatusCreated, city)

}

func (h Handler) GetCityByID(w http.ResponseWriter, r *http.Request) {
	values := r.URL.Query()

	if len(values["id"]) <= 0 {
		handleResponse(w, http.StatusBadRequest, errors.New("id is required"))
		return
	}

	id := values["id"][0]

	city, err := h.storage.City().Get(id)
	if err != nil {
		handleResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	handleResponse(w, http.StatusOK, city)
}

func (h Handler) GetCityList(w http.ResponseWriter, r *http.Request) {
	var (
		page  = 1
		limit = 10
		err   error
	)

	values := r.URL.Query()

	if len(values["page"]) > 0 {
		page, err = strconv.Atoi(values["page"][0])
		if err != nil {
			page = 1
		}
	}

	if len(values["limit"]) > 0 {
		limit, err = strconv.Atoi(values["limit"][0])
		if err != nil {
			limit = 10
		}
	}

	resp, err := h.storage.City().GetList(models.GetListRequest{
		Page:  page,
		Limit: limit,
	})

	handleResponse(w, http.StatusOK, resp)
}

func (h Handler) UpdateCity(w http.ResponseWriter, r *http.Request) {
	updatedCity := models.City{}

	if err := json.NewDecoder(r.Body).Decode(&updatedCity); err != nil {
		handleResponse(w, http.StatusBadRequest, err)
		return
	}

	id, err := h.storage.City().Update(updatedCity)
	if err != nil {
		handleResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	city, err := h.storage.City().Get(id)
	if err != nil {
		handleResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	handleResponse(w, http.StatusOK, city)
}

func (h Handler) DeleteCity(w http.ResponseWriter, r *http.Request) {
	values := r.URL.Query()
	if len(values["id"]) <= 0 {
		handleResponse(w, http.StatusBadRequest, errors.New("id is required"))
		return
	}

	id := values["id"][0]

	if err := h.storage.City().Delete(id); err != nil {
		handleResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	handleResponse(w, http.StatusOK, "city successfully deleted!")
}

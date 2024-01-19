package handler

import (
	"city2city/api/models"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
)

func (h Handler) Car(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		h.CreateCar(w, r)
	case http.MethodGet:
		values := r.URL.Query()
		if _, ok := values["id"]; !ok {
			h.GetCarList(w, r)
		} else {
			h.GetCarByID(w, r)
		}
	case http.MethodPut:
		h.UpdateCar(w, r)
	case http.MethodDelete:
		h.DeleteCar(w, r)
	case http.MethodPatch:
		values := r.URL.Query()
		if _, ok := values["status"]; ok {
			h.UpdateCarStatus(w, r)
		} else {
			h.UpdateCarRoute(w, r)
		}
	}
}

func (h Handler) CreateCar(w http.ResponseWriter, r *http.Request) {
	car := models.CreateCar{}

	if err := json.NewDecoder(r.Body).Decode(&car); err != nil {
		handleResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	id, err := h.storage.Car().Create(car)
	if err != nil {
		handleResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	createdCar, err := h.storage.Car().Get(id)
	if err != nil {
		handleResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	handleResponse(w, http.StatusOK, createdCar)
}

func (h Handler) GetCarByID(w http.ResponseWriter, r *http.Request) {
	values := r.URL.Query()
	if len(values["id"]) <= 0 {
		handleResponse(w, http.StatusBadRequest, errors.New("error is required"))
		return
	}
	id := values["id"][0]
	car, err := h.storage.Car().Get(id)
	if err != nil {
		handleResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	handleResponse(w, http.StatusOK, car)
}

func (h Handler) GetCarList(w http.ResponseWriter, r *http.Request) {
	var (
		page, limit = 1, 10
		err         error
	)
	values := r.URL.Query()
	if len(values["page"]) > 0 {
		page, err = strconv.Atoi(values["id"][0])
		if err != nil {
			page = 1
		}
	}

	if len(values["limit"]) > 0 {
		limit, err = strconv.Atoi(values["id"][0])
		if err != nil {
			limit = 10
		}
	}

	cars, err := h.storage.Car().GetList(models.GetListRequest{
		Page:  page,
		Limit: limit,
	})
	if err != nil {
		handleResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	handleResponse(w, http.StatusOK, cars)
}

func (h Handler) UpdateCar(w http.ResponseWriter, r *http.Request) {
	car := models.Car{}

	if err := json.NewDecoder(r.Body).Decode(&car); err != nil {
		handleResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	id, err := h.storage.Car().Update(car)
	if err != nil {
		handleResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	updatedCar, err := h.storage.Car().Get(id)
	if err != nil {
		handleResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	handleResponse(w, http.StatusOK, updatedCar)
}

func (h Handler) DeleteCar(w http.ResponseWriter, r *http.Request) {
	values := r.URL.Query()
	if len(values["id"]) <= 0 {
		handleResponse(w, http.StatusBadRequest, errors.New("id is required"))
		return
	}
	id := values["id"][0]
	if err := h.storage.Car().Delete(id); err != nil {
		handleResponse(w, http.StatusInternalServerError, err.Error())
		return
	}
	handleResponse(w, http.StatusOK, "car is deleted successfully!")
}

func (h Handler) UpdateCarRoute(w http.ResponseWriter, r *http.Request) {
	route := models.UpdateCarRoute{}

	if err := json.NewDecoder(r.Body).Decode(&route); err != nil {
		handleResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	if err := h.storage.Car().UpdateCarRoute(route); err != nil {
		handleResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	handleResponse(w, http.StatusOK, "car route updated!")
}

func (h Handler) UpdateCarStatus(w http.ResponseWriter, r *http.Request) {
	status := models.UpdateCarStatus{}

	if err := json.NewDecoder(r.Body).Decode(&status); err != nil {
		handleResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	if err := h.storage.Car().UpdateCarStatus(status); err != nil {
		handleResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	handleResponse(w, http.StatusOK, "car status updated!")
}

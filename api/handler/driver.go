package handler

import (
	"city2city/api/models"
	"city2city/check"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
)

func (h Handler) Driver(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		h.CreateDriver(w, r)
	case http.MethodGet:
		values := r.URL.Query()
		_, ok := values["id"]
		if !ok {
			h.GetDriverList(w, r)
		} else {
			h.GetDriverByID(w, r)
		}
	case http.MethodPut:
		h.UpdateDriver(w, r)
	case http.MethodDelete:
		h.DeleteDriver(w, r)
	}
}

func (h Handler) CreateDriver(w http.ResponseWriter, r *http.Request) {
	driver := models.CreateDriver{}

	if err := json.NewDecoder(r.Body).Decode(&driver); err != nil {
		handleResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	if !check.PhoneNumber(driver.Phone) {
		handleResponse(w, http.StatusBadRequest, errors.New("phone format is not correct"))
		return
	}

	id, err := h.storage.Driver().Create(driver)
	if err != nil {
		handleResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	createdDriver, err := h.storage.Driver().Get(id)
	if err != nil {
		handleResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	handleResponse(w, http.StatusCreated, createdDriver)
}

func (h Handler) GetDriverByID(w http.ResponseWriter, r *http.Request) {
	values := r.URL.Query()
	if len(values["id"]) <= 0 {
		handleResponse(w, http.StatusBadRequest, errors.New("id is required"))
		return
	}

	id := values["id"][0]
	driver, err := h.storage.Driver().Get(id)
	if err != nil {
		handleResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	handleResponse(w, http.StatusOK, driver)
}

func (h Handler) GetDriverList(w http.ResponseWriter, r *http.Request) {
	var (
		page, limit = 1, 10
		err         error
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

	drivers, err := h.storage.Driver().GetList(models.GetListRequest{
		Page:  page,
		Limit: limit,
	})
	if err != nil {
		handleResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	handleResponse(w, http.StatusOK, drivers)
}

func (h Handler) UpdateDriver(w http.ResponseWriter, r *http.Request) {
	driver := models.Driver{}

	if err := json.NewDecoder(r.Body).Decode(&driver); err != nil {
		handleResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	if !check.PhoneNumber(driver.Phone) {
		handleResponse(w, http.StatusBadRequest, errors.New("phone format is not correct"))
		return
	}

	id, err := h.storage.Driver().Update(driver)
	if err != nil {
		handleResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	updatedDriver, err := h.storage.Driver().Get(id)
	if err != nil {
		handleResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	handleResponse(w, http.StatusOK, updatedDriver)
}

func (h Handler) DeleteDriver(w http.ResponseWriter, r *http.Request) {
	values := r.URL.Query()
	if len(values["id"]) <= 0 {
		handleResponse(w, http.StatusBadRequest, errors.New("id is required"))
		return
	}

	id := values["id"][0]
	if err := h.storage.Driver().Delete(id); err != nil {
		handleResponse(w, http.StatusInternalServerError, err)
		return
	}

	handleResponse(w, http.StatusOK, "driver is successfully deleted!")
}

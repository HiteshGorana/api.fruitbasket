package main

import (
	"api.fruitbasket/internals/data"
	"errors"
	"fmt"
	"net/http"
)

func (app *application) createFoodHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		FruitName string  `json:"fruit_name"`
		Price     float64 `json:"price"`
	}

	err := app.readJSON(w, r, &input)
	if err != nil {
		fmt.Println(1)
		app.badRequestResponse(w, r, err)
		return
	}

	fruit := &data.Fruit{
		FruitName: input.FruitName,
		Price:     input.Price,
	}

	//v := validator.New()

	//if data.ValidateFruit(v, fruit); !v.Valid() {
	//	fmt.Println("HI")
	//	app.failedValidationResponse(w, r, v.Errors)
	//	return
	//}

	err = app.models.Fruit.Insert(fruit)

	if err != nil {
		fmt.Println(2)
		app.serverErrorResponse(w, r, err)
		return
	}

	headers := make(http.Header)
	headers.Set("Location", fmt.Sprintf("/v1/fruits/%d", fruit.ID))

	err = app.writeJSON(w, http.StatusCreated, envelope{"fruit": fruit}, headers)
	if err != nil {
		fmt.Println(3)
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) showFoodHandler(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIDParam(r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}

	fruit, err := app.models.Fruit.Get(id)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"fruit": fruit}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) deleteFruitHandler(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIDParam(r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}

	err = app.models.Fruit.Delete(id)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"message": "movie successfully deleted"}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

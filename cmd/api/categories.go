package main

import (
	"errors"
	"fmt"
	"github.com/KZhambyl/HistoricalFigures/internal/data"
	"github.com/KZhambyl/HistoricalFigures/internal/validator"
	"net/http"
)

func (app *application) createCategoryHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Name    string `json:"name"`
		Version int32  `json:"version"`
	}

	err := app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	category := &data.Category{
		Name:    input.Name,
		Version: input.Version,
	}

	v := validator.New()

	if data.ValidateCategory(v, category); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	err = app.models.Categories.Insert(category)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	headers := make(http.Header)
	headers.Set("Location", fmt.Sprintf("/v1/categories/%d", category.ID))

	err = app.writeJSON(w, http.StatusCreated, envelope{"category": category}, headers)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) showCategoryHandler(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIDParam(r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}

	category, err := app.models.Categories.Get(id)
	if err != nil {
		if err != nil {
			switch {
			case errors.Is(err, data.ErrRecordNotFound):
				app.notFoundResponse(w, r)
			default:
				app.serverErrorResponse(w, r, err)
			}
			return
		}
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"category": category}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) updateCategoryHandler(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIDParam(r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}

	category, err := app.models.Categories.Get(id)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	var input struct {
		Name *string `json:"name"`
	}

	err = app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	if input.Name != nil {
		category.Name = *input.Name
	}

	v := validator.New()
	if data.ValidateCategory(v, category); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	err = app.models.Categories.Update(category)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrEditConflict):
			app.editConflictResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"category": category}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) deleteCategoryHandler(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIDParam(r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}
	err = app.models.Categories.Delete(id)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}
	err = app.writeJSON(w, http.StatusOK, envelope{"message": "category successfully deleted"}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) listCategoriesHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Name string
		data.Filters
	}
	v := validator.New()
	qs := r.URL.Query()

	input.Name = app.readString(qs, "name", "")
	input.Filters.Page = app.readInt(qs, "page", 1, v)
	input.Filters.PageSize = app.readInt(qs, "page_size", 20, v)
	input.Filters.Sort = app.readString(qs, "sort", "id")

	input.Filters.SortSafelist = []string{"id", "name", "-id", "-name"}

	if data.ValidateFilters(v, input.Filters); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	categories, metadata, err := app.models.Categories.GetAll(input.Name, input.Filters)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"categories": categories, "metadata": metadata}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) showCategoryFiguresHandler(w http.ResponseWriter, r *http.Request) {

	id, err := app.readIDParam(r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}

	var input struct {
		Name        string
		YearsOfLife string
		data.Filters
	}
	v := validator.New()
	qs := r.URL.Query()

	input.Name = app.readString(qs, "name", "")
	input.YearsOfLife = app.readString(qs, "years_of_life", "")
	input.Filters.Page = app.readInt(qs, "page", 1, v)
	input.Filters.PageSize = app.readInt(qs, "page_size", 20, v)
	input.Filters.Sort = app.readString(qs, "sort", "id")

	input.Filters.SortSafelist = []string{"id", "name", "years_of_life", "-id", "-name", "-years_of_life"}

	if data.ValidateFilters(v, input.Filters); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	figures, metadata, err := app.models.Categories.GetCategoryFigures(input.Name, input.YearsOfLife, input.Filters, id)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"figures": figures, "metadata": metadata}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

package main

import (
	"expvar"
	"github.com/julienschmidt/httprouter"
	"net/http"
)

func (app *application) routes() http.Handler {
	router := httprouter.New()
	router.NotFound = http.HandlerFunc(app.notFoundResponse)
	router.MethodNotAllowed = http.HandlerFunc(app.methodNotAllowedResponse)

	router.HandlerFunc(http.MethodGet, "/v1/healthcheck", app.healthcheckHandler)

	router.HandlerFunc(http.MethodGet, "/v1/figures", app.requirePermission("figures:read", app.listFiguresHandler))
	router.HandlerFunc(http.MethodPost, "/v1/figures", app.requirePermission("figures:write", app.createFigureHandler))
	router.HandlerFunc(http.MethodGet, "/v1/figures/:id", app.requirePermission("figures:read", app.showFigureHandler))
	router.HandlerFunc(http.MethodPatch, "/v1/figures/:id", app.requirePermission("figures:write", app.updateFigureHandler))
	router.HandlerFunc(http.MethodDelete, "/v1/figures/:id", app.requirePermission("figures:write", app.deleteFigureHandler))

	router.HandlerFunc(http.MethodPost, "/v1/users", app.registerUserHandler)
	router.HandlerFunc(http.MethodPut, "/v1/users/activated", app.activateUserHandler)
	router.HandlerFunc(http.MethodPost, "/v1/tokens/authentication", app.createAuthenticationTokenHandler)

	router.HandlerFunc(http.MethodGet, "/v1/categories", app.requirePermission("figures:read", app.listCategoriesHandler))
	router.HandlerFunc(http.MethodPost, "/v1/categories", app.requirePermission("figures:write", app.createCategoryHandler))
	router.HandlerFunc(http.MethodGet, "/v1/categories/:id", app.requirePermission("figures:read", app.showCategoryHandler))
	router.HandlerFunc(http.MethodPatch, "/v1/categories/:id", app.requirePermission("figures:write", app.updateCategoryHandler))
	router.HandlerFunc(http.MethodDelete, "/v1/categories/:id", app.requirePermission("figures:write", app.deleteCategoryHandler))
	router.HandlerFunc(http.MethodGet, "/v1/categories/:id/figures", app.requirePermission("figures:read", app.showCategoryFiguresHandler))

	router.Handler(http.MethodGet, "/v1/metrics", expvar.Handler())

	return app.metrics(app.recoverPanic(app.enableCORS(app.rateLimit(app.authenticate(router)))))
}

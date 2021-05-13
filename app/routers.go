package app

import (
	"golive/data"
	"golive/envvar"
	"html/template"
	"net/http"

	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/mongo"
)

type appDB struct {
	client          *mongo.Client
	mDB             *mongo.Database
	tpl             *template.Template
	env             envvar.Env
	dataCollections dataCollection
}

//SetRouter registers all the URL path and the handle functions
func SetRouter(db *data.AppDB) *mux.Router {
	app := initAppDB(db)
	app.dataCollections = loadInitData(db)

	r := mux.NewRouter()
	app.tpl = initTemplate()

	fs := http.FileServer(http.Dir("./static"))
	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", fs))
	r.Handle("/favicon.ico", http.NotFoundHandler())

	getR := r.Methods("GET").Subrouter()

	getR.HandleFunc("/api/v1/location/singapore", app.getSingapore)
	getR.HandleFunc("/api/v1/location/points", app.getLocationPoints)
	getR.HandleFunc("/api/v1/location/info", app.getLocationInfo)
	getR.HandleFunc("/api/v1/location/address", app.getAddressInfo)
	getR.HandleFunc("/api/v1/filter", app.getFilterList)
	getR.HandleFunc("/api/v1/provider", app.getProviderLocation)
	getR.Use(middlewareRecovery)
	getR.Use(middlewareLogging)
	getR.Use(middlewareValidateContentType)
	getR.Use(middlewareAddContentType)

	getTR := r.Methods("GET").Subrouter()
	getTR.HandleFunc("/", app.getInitMap)
	getTR.HandleFunc("/login", app.getLogin)
	getTR.HandleFunc("/logout", app.logout)
	getTR.HandleFunc("/register", app.getRegister)
	getTR.HandleFunc("/vehicle", app.getVehicle)
	getTR.HandleFunc("/updateProfile", app.updateProfile)
	getTR.HandleFunc("/changePassword", app.getChangePwd)
	getTR.HandleFunc("/changeEmail", app.getChangeEmail)
	getTR.HandleFunc("/cancel", cancel)
	getTR.HandleFunc("/view", app.getView)
	getTR.Use(middlewareRecovery)
	getTR.Use(middlewareLogging)

	postR := r.Methods("POST").Subrouter()
	postR.HandleFunc("/login", app.postLogin)
	postR.HandleFunc("/register", app.postRegister)
	postR.HandleFunc("/vehicle", app.postVehicle)
	postR.HandleFunc("/changePassword", app.postChangePwd)
	postR.HandleFunc("/changeEmail", app.postChangeEmail)
	postR.HandleFunc("/view", app.postView)
	postR.Use(middlewareRecovery)
	postR.Use(middlewareLogging)
	return r
}

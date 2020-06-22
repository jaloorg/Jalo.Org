package routes

import (
	"encoding/json"
	"log"
	"net/http"

	"Ayudaap.org/models"
	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Obtiene el handler principal de la aplicacion
func GetHandler() http.Handler {
	r := mux.NewRouter()
	r.StrictSlash(true)

	r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		log.Printf("Peticion en %v desde %v", r.RequestURI, r.RemoteAddr)
		w.Write([]byte("Hola desde default"))
	})

	// TODO: Cambiar la version de API desde configuracion
	api := r.PathPrefix("/api/v1").Subrouter()
	api.StrictSlash(true)

	api.HandleFunc("/oID/", GetPrimitiveID).Methods("GET").Name("getOId")

	apiOrg := api.PathPrefix("/organizacion").Subrouter()
	apiOrg.HandleFunc("/", CreateOrganizacion).Methods("POST").Name("crearOrganizacion")
	apiOrg.HandleFunc("/", GetALlOrganizacionesReq).Methods("GET").Name("obtenerOrganizaciones")
	apiOrg.HandleFunc("/", UpsertOrganizacion).Methods("PUT").Name("modificarOrganizacion")
	apiOrg.HandleFunc("/{id}", GetOrganizacionById).Methods("GET").Name("getOrganizacionById")
	apiOrg.HandleFunc("/{id}", DeleteOrganizacion).Methods("DELETE").Name("borrarOrganizacion")

	return r
}

// GetError : This is helper function to prepare error model.
// If you want to export your function. You must to start upper case function name. Otherwise you won't see your function when you import that on other class.
func GetError(err error, w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json")

	var response = models.ErrorResponse{
		ErrorMessage: err.Error(),
		StatusCode:   http.StatusInternalServerError,
	}

	message, _ := json.Marshal(response)

	w.WriteHeader(response.StatusCode)
	w.Write(message)
}

// Obtiene un Id General
func GetPrimitiveID(w http.ResponseWriter, r *http.Request) {
	id := primitive.NewObjectID()

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(struct {
		Id string `json:"id,omitempty"`
	}{Id: id.Hex()})
}

package api

import (
	"fmt"
	"net/http"
	"strconv"

	eram "github.com/Onefootball/entity-rest-api/manager"
	"github.com/ant0ine/go-json-rest/rest"
)

const (
	Offset           = "0"
	Limit            = "10"
	OrderDir         = "ASC"
	StatusCodeHeader = "X-Status-Code"
	EntityIDHeader   = "X-Entity-ID"
	LocationHeader   = "Location"
)

type EntityRestAPI struct {
	em *eram.EntityDbManager
}

func NewEntityRestAPI(em *eram.EntityDbManager) *EntityRestAPI {
	return &EntityRestAPI{
		em,
	}
}

func (api *EntityRestAPI) GetAllEntities(w rest.ResponseWriter, r *rest.Request) {
	entity := r.PathParam("entity")
	qs := r.Request.URL.Query()

	limit, offset, orderBy, orderDir := qs.Get("_perPage"), qs.Get("_page"), qs.Get("_sortField"), qs.Get("_sortDir")

	qs.Del("_perPage")
	qs.Del("_page")
	qs.Del("_sortField")
	qs.Del("_sortDir")

	filterParams := make(map[string]string)

	// remaining GET parameters are used to filter the result
	for filterName, _ := range qs {
		filterParams[filterName] = qs.Get(filterName)
	}

	if offset == "" {
		offset = Offset
	}

	if limit == "" {
		limit = Limit
	}

	if orderBy == "" {
		orderBy = api.em.GetIdColumn(entity)
	}

	if orderDir == "" {
		orderDir = OrderDir
	}

	allResults, count, dbErr := api.em.GetEntities(entity, filterParams, limit, offset, orderBy, orderDir)

	if dbErr != nil {
		rest.Error(w, dbErr.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Access-Control-Expose-Headers", "X-Total-Count")
	w.Header().Set("X-Total-Count", fmt.Sprintf("%d", count))

	w.WriteJson(allResults)
}

func (api *EntityRestAPI) GetEntity(w rest.ResponseWriter, r *rest.Request) {
	id := r.PathParam("id")
	entity := r.PathParam("entity")
	result, err := api.em.GetEntity(entity, id)
	if err != nil {
		rest.Error(w, err.Error(), http.StatusInternalServerError)
		return
	} else if len(result) <= 0 {
		rest.Error(w, "Not Found", http.StatusNotFound)
		return
	}

	w.WriteJson(result)
}

func (api *EntityRestAPI) PostEntity(w rest.ResponseWriter, r *rest.Request) {
	w.Header().Add("Access-Control-Expose-Headers", StatusCodeHeader)
	w.Header().Add("Access-Control-Expose-Headers", EntityIDHeader)

	entity := r.PathParam("entity")
	postData := map[string]interface{}{}
	if err := r.DecodeJsonPayload(&postData); err != nil {
		w.Header().Set(StatusCodeHeader, fmt.Sprintf("%d", http.StatusInternalServerError))
		rest.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	newId, err := api.em.PostEntity(entity, postData)
	if err != nil {
		w.Header().Set(StatusCodeHeader, fmt.Sprintf("%d", http.StatusInternalServerError))
		rest.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	insertedEntity, err := api.em.GetEntity(entity, strconv.FormatInt(newId, 10))
	if err != nil {
		w.Header().Set(StatusCodeHeader, fmt.Sprintf("%d", http.StatusInternalServerError))
		rest.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set(LocationHeader, fmt.Sprintf("%s/%d", entity, newId))
	w.Header().Set(StatusCodeHeader, fmt.Sprintf("%d", http.StatusCreated))
	w.Header().Set(EntityIDHeader, fmt.Sprintf("%d", newId))

	w.WriteHeader(http.StatusCreated)
	w.WriteJson(insertedEntity)
}

func (api *EntityRestAPI) PutEntity(w rest.ResponseWriter, r *rest.Request) {
	id := r.PathParam("id")
	entity := r.PathParam("entity")
	updated := map[string]interface{}{}
	if err := r.DecodeJsonPayload(&updated); err != nil {
		rest.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	rowsAffected, updatedEntity, err := api.em.UpdateEntity(entity, id, updated)
	if err != nil {
		rest.Error(w, err.Error(), http.StatusInternalServerError)
		return
	} else if len(updatedEntity) <= 0 {
		rest.Error(w, "Not Found", http.StatusNotFound)
		return
	}

	if rowsAffected == 0 {
		w.WriteHeader(http.StatusNoContent)
	} else {
		w.WriteHeader(http.StatusOK)
	}

	w.WriteJson(updatedEntity)
}

func (api *EntityRestAPI) DeleteEntity(w rest.ResponseWriter, r *rest.Request) {
	id := r.PathParam("id")
	entity := r.PathParam("entity")
	rowsAffected, err := api.em.DeleteEntity(entity, id)
	if err != nil {
		rest.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if rowsAffected == 0 {
		w.WriteHeader(http.StatusNotFound)
	} else {
		w.WriteHeader(http.StatusOK)
	}
}

/*
 * Data Anonymization Server
 *
 * This is a data anonymization server. You can set the anonymization requirements for the different datasets individually, and upload data to them. The uploaded data is anonymized on the server and can be then downloaded.
 *
 * API version: 0.1-alpha
 * Generated by: Swagger Codegen (https://github.com/swagger-api/swagger-codegen.git)
 */

package swagger

import (
	"anondb"
	"net/http"

	"github.com/gorilla/mux"
)

func datasetsGet(w http.ResponseWriter, r *http.Request) {
	datasetList, err := anondb.ListDatasets()

	if err != nil {
		logDBError(err)
		respondWithError(w, http.StatusInternalServerError, err.Error())
	} else {
		respondWithJSON(w, http.StatusOK, createDatasetsResponse(datasetList))
	}
}

func datasetsNameDelete(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	if err := anondb.DropDataset(vars["name"]); err != nil {
		handleDBNotFound(err, w, http.StatusNotFound, "The dataset with the specified name was not found")
	} else {
		w.WriteHeader(http.StatusNoContent)
	}
}

func datasetsNameGet(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	dataset, err := anondb.GetDataset(vars["name"])
	if err != nil {
		handleDBNotFound(err, w, http.StatusNotFound, "The dataset with the specified name was not found")
	} else {
		respondWithJSON(w, http.StatusOK, createDatasetResponse(&dataset))
	}
}

func datasetsNamePut(w http.ResponseWriter, r *http.Request) {
	var request CreateDatasetRequest
	if !tryReadRequestBody(r, &request, w) {
		return
	}

	vars := mux.Vars(r)
	dataset := createDataset(vars["name"], &request)

	if err := dataset.Validate(); err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	if err := anondb.CreateDataset(&dataset); err != nil {
		handleDBDuplicate(err, w, http.StatusConflict, "A dataset with the specified name already exists.")
	} else {
		respondWithJSON(w, http.StatusCreated, createDatasetResponse(&dataset))
	}
}

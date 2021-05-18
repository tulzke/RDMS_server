package rest

import (
	"Rostelecom_Device_Management_System/database"
	"Rostelecom_Device_Management_System/security"
	"Rostelecom_Device_Management_System/structs"
	"encoding/json"
	"net/http"
)

func RegisteringWorkStation(w http.ResponseWriter, r *http.Request) {
	var recv structs.Workstation

	err := json.NewDecoder(r.Body).Decode(&recv)
	if err != nil {
		SendResponse(http.StatusBadRequest, &w, "Can't parse")
		return
	}

	ws, err := database.GetWorkstationByName(recv.Name)
	if err != nil && err.Error() == "Workstation not found" {
		SendResponse(http.StatusNoContent, &w, "Workstation not found")
		return
	}

	if ws.Serial != "" {
		SendResponse(http.StatusForbidden, &w, "Workstation with this serial number is already exists")
		return
	}

	err = database.RegisterWS(&recv)
	if err != nil {
		SendResponse(http.StatusInternalServerError, &w, "Internal error")
		return
	}

	publicWS := structs.PublicWorkstation{}.FromWs(recv)

	SendResponse(http.StatusOK, &w, publicWS)
}

func GetWorkstations(w http.ResponseWriter, r *http.Request) {
	if status, _ := security.JwtAuth(r); !status {
		SendResponse(http.StatusUnauthorized, &w, "Unauthorized")
		return
	}
	result := database.GetWorkstations()
	SendResponse(http.StatusOK, &w, result)
}
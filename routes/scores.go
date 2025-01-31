package routes

import (
	"encoding/json"
	"fmt"
	"net/http"
	"platform/db"
	"platform/middleware"
	"strconv"
	"strings"
)

type GraphPoint struct {
	X string `json:"x"`
	Y int    `json:"y"`
}

type ScoresData struct {
	Data
	Users []db.UserScore
}

func submit(ctx *middleware.Ctx) {
	challID := ctx.FormValue("challID")
	flag := strings.TrimSpace(ctx.FormValue("flag"))

	if ctx.User == nil {
		ctx.AddFlash("You must be logged in to submit flags")
		ctx.Error("Unauthorized", http.StatusUnauthorized)
		return
	}

	chalID, err := strconv.Atoi(challID)
	if err != nil {
		ctx.InternalError(fmt.Errorf("error converting challID to int: %v", err))
		return
	}

	status, err := db.SubmitFlag(ctx.User, chalID, flag)
	if err != nil {
		ctx.InternalError(fmt.Errorf("error submitting flag: %v", err))
		return
	}

	header := http.StatusOK
	switch status {
	case db.StatusCorrectFlag:
		header = http.StatusAccepted
	case db.StatusAlreadySolved:
		header = http.StatusConflict
		ctx.AddFlash("Challenge already solved", "warning")
	case db.StatusWrongFlag:
		header = http.StatusNotAcceptable
	}

	ctx.WriteHeader(header)
}

func scores(ctx *middleware.Ctx) {
	tmpl := getTemplate(ctx, "scores")
	if tmpl == nil {
		return
	}

	users, err := db.GetUsersScores()
	if err != nil {
		ctx.InternalError(fmt.Errorf("error getting users: %v", err))
		return
	}

	data := ScoresData{Data: Data{}}
	data.Users = users
	data.User = ctx.User

	executeTemplate(ctx, tmpl, &data)
}

func graphData(ctx *middleware.Ctx) {
	data, err := db.GetGraphData()
	if err != nil {
		ctx.InternalError(fmt.Errorf("error getting graph data: %v", err))
		return
	}

	jsonData := make(map[string][]GraphPoint)
	for _, d := range data {
		if len(jsonData[d.User]) == 0 {
			point := GraphPoint{
				X: data[0].Timestamp.Format("2006-01-02 15:04:05"),
				Y: 0,
			}
			jsonData[d.User] = append(jsonData[d.User], point)
		}
		point := GraphPoint{
			X: d.Timestamp.Format("2006-01-02 15:04:05"),
			Y: jsonData[d.User][len(jsonData[d.User])-1].Y + d.Points,
		}
		jsonData[d.User] = append(jsonData[d.User], point)
	}

	response, err := json.Marshal(jsonData)
	if err != nil {
		ctx.InternalError(fmt.Errorf("error marshaling json: %v", err))
		return
	}
	ctx.Write(response)
}

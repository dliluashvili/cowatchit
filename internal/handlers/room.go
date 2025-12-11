package handlers

import (
	"fmt"
	"net/http"

	"github.com/dliluashvili/cowatchit/internal/dtos"
	"github.com/dliluashvili/cowatchit/internal/helpers"
	"github.com/dliluashvili/cowatchit/internal/models"
	"github.com/dliluashvili/cowatchit/internal/params"
	"github.com/dliluashvili/cowatchit/internal/services"
	"github.com/dliluashvili/cowatchit/internal/shared/constants"
	"github.com/dliluashvili/cowatchit/internal/templates"
	"github.com/go-chi/chi"
	"github.com/google/uuid"
)

type Roomhandler struct {
	roomService *services.RoomService
}

func NewRoomHandler(rs *services.RoomService) *Roomhandler {
	return &Roomhandler{
		roomService: rs,
	}
}

func (rh *Roomhandler) HandleRoomsPage(w http.ResponseWriter, r *http.Request) {
	session := r.Context().Value(constants.SessionContextKey).(*models.Session)

	userID := session.User.ID

	filter := r.URL.Query().Get("filter")
	keyword := r.URL.Query().Get("keyword")
	myStr := r.URL.Query().Get("my")

	my := helpers.StringToBool(myStr)

	findRoomDto := &dtos.FindRoomDto{
		Filter:     &filter,
		Keyword:    &keyword,
		My:         my,
		AuthUserID: &userID,
	}

	rooms, err := rh.roomService.Find(findRoomDto)

	if err != nil {

	}

	if filter == "" {
		filter = "all"
	}

	roomsParams := &params.RoomsParams{
		Filter:  filter,
		Keyword: keyword,
		My:      myStr,
	}

	if r.Header.Get("HX-Request") == "true" {
		templates.Rooms(rooms, roomsParams).Render(r.Context(), w)
		return
	}

	templates.RoomsPage(rooms, roomsParams).Render(r.Context(), w)
}

func (rh *Roomhandler) HandleCreateRoomPage(w http.ResponseWriter, r *http.Request) {
	if r.Header.Get("HX-Request") == "true" {
		templates.CreateRoom().Render(r.Context(), w)
	} else {
		templates.CreateRoomPage().Render(r.Context(), w)
	}
}

func (rh *Roomhandler) Create(w http.ResponseWriter, r *http.Request) {
	validated := r.Context().Value(constants.ValidatedContextKey).(*dtos.CreateRoomDto)

	session := r.Context().Value(constants.SessionContextKey).(*models.Session)

	userID := session.User.ID

	createRoomServiceDto := &dtos.CreateRoomServiceDto{
		HostID:        userID,
		CreateRoomDto: validated,
	}

	_, err := rh.roomService.Create(r.Context(), createRoomServiceDto)

	if err != nil {
		fmt.Println("roomHandler@Create", err)
		helpers.SendJson(w, &helpers.Response{
			Data: map[string]bool{
				"success": false,
			},
			Message: "Unable to create potential pair",
			Status:  http.StatusBadRequest,
		})

		return
	}

	helpers.SendJson(w, &helpers.Response{
		Data: map[string]bool{
			"success": true,
		},
		Message: "all good",
		Status:  http.StatusOK,
	})
}

func (rh *Roomhandler) HandleRoomJoinPage(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")

	ID, err := uuid.Parse(idStr)

	if err != nil {
		fmt.Println("err", err)
		helpers.SendJson(w, &helpers.Response{
			Data:    nil,
			Message: "bad request",
			Status:  http.StatusBadRequest,
		})

		return
	}

	room, err := rh.roomService.FindOne(ID)

	if err != nil {
		fmt.Println("err", err)

		helpers.SendJson(w, &helpers.Response{
			Data:    nil,
			Message: "bad request",
			Status:  http.StatusBadRequest,
		})

		return
	}

	if r.Header.Get("HX-Request") == "true" {
		templates.JoinRoom(room).Render(r.Context(), w)
		return
	}

	templates.JoinRoomPage(room).Render(r.Context(), w)
}

func (rh *Roomhandler) HandleRoomPage(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")

	ID, err := uuid.Parse(idStr)

	if err != nil {
		fmt.Println("err", err)
		helpers.SendJson(w, &helpers.Response{
			Data:    nil,
			Message: "bad request",
			Status:  http.StatusBadRequest,
		})
		return
	}

	exists, err := rh.roomService.Exists(ID)

	if err != nil || !exists {
		fmt.Println("err", err)
		helpers.SendJson(w, &helpers.Response{
			Data:    nil,
			Message: "bad request",
			Status:  http.StatusNotFound,
		})

		return
	}

	if r.Header.Get("HX-Request") == "true" {
		templates.JoiningView(ID.String()).Render(r.Context(), w)
		return
	}

	templates.RoomPage(idStr).Render(r.Context(), w)
}

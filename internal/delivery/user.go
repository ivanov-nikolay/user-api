package delivery

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/ivanov-nikolay/user-api/internal/dto"
	"github.com/ivanov-nikolay/user-api/internal/filters"
	"github.com/ivanov-nikolay/user-api/internal/usecase"
	"go.uber.org/zap"
)

type UserHandler struct {
	u      usecase.UserUseCase
	logger *zap.SugaredLogger
}

func New(u usecase.UserUseCase, logger *zap.SugaredLogger) *UserHandler {
	return &UserHandler{
		u:      u,
		logger: logger,
	}
}

func (uh *UserHandler) CreateUserHandler(w http.ResponseWriter, r *http.Request) {
	userCreateDTO := &dto.UserCreate{}
	rBody, err := io.ReadAll(r.Body)
	if err != nil {
		errText := fmt.Sprintf(`{"message": "error in reading request body: %s"}`, err)
		uh.logger.Errorf(errText)
		writeResponse(uh.logger, w, []byte(errText), http.StatusBadRequest)
		return
	}
	err = json.Unmarshal(rBody, userCreateDTO)
	if err != nil {
		errText := fmt.Sprintf(`{"message": "error in decoding user: %s"}`, err)
		uh.logger.Errorf(errText)
		writeResponse(uh.logger, w, []byte(errText), http.StatusBadRequest)
		return
	}

	if validationErrors := userCreateDTO.Validate(); len(validationErrors) != 0 {
		var errorsJSON []byte
		errorsJSON, err = json.Marshal(validationErrors)
		if err != nil {
			errText := fmt.Sprintf(`{"message": "error in json decoding: %s"}`, err)
			uh.logger.Errorf(errText)
			writeResponse(uh.logger, w, []byte(errText), http.StatusInternalServerError)
			return
		}
		writeResponse(uh.logger, w, errorsJSON, http.StatusUnprocessableEntity)
		return
	}

	user := userCreateDTO.ConvertToUser()
	addedUser, err := uh.u.CreateUserUseCase(user)
	if err != nil {
		errText := fmt.Sprintf(`{"message": "internal server error"}`)
		uh.logger.Errorf(errText)
		writeResponse(uh.logger, w, []byte(errText), http.StatusInternalServerError)
		return
	}
	userJSON, err := json.Marshal(addedUser)
	if err != nil {
		errText := fmt.Sprintf(`{"message": "error in coding user: %s"}`, err)
		uh.logger.Errorf(errText)
		writeResponse(uh.logger, w, []byte(errText), http.StatusInternalServerError)
		return
	}
	writeResponse(uh.logger, w, userJSON, http.StatusOK)
}

func (uh *UserHandler) DeleteUserHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userID := vars["USER_ID"]
	userIDInt, err := strconv.ParseInt(userID, 10, 64)
	if err != nil {
		errText := fmt.Sprintf(`{"message": "bad format of user id: %s"}`, err)
		uh.logger.Errorf(errText)
		writeResponse(uh.logger, w, []byte(errText), http.StatusBadRequest)
		return
	}
	wasDeleted, err := uh.u.DeleteUserUseCase(userIDInt)
	if err != nil {
		errText := fmt.Sprintf(`{"message": "internal server error"}`)
		uh.logger.Errorf(errText)
		writeResponse(uh.logger, w, []byte(errText), http.StatusInternalServerError)
		return
	}
	if !wasDeleted {
		errText := fmt.Sprintf(`{"message": "user with ID %d is not found"}`, userIDInt)
		uh.logger.Errorf(errText)
		writeResponse(uh.logger, w, []byte(errText), http.StatusNotFound)
		return
	}
	result := `{"result": "success"}`
	writeResponse(uh.logger, w, []byte(result), http.StatusOK)
}

func (uh *UserHandler) UpdateUserHandler(w http.ResponseWriter, r *http.Request) {
	userUpdateDTO := &dto.UserUpdate{}
	rBody, err := io.ReadAll(r.Body)
	if err != nil {
		errText := fmt.Sprintf(`{"message": "error in reading request body: %s"}`, err)
		uh.logger.Errorf(errText)
		writeResponse(uh.logger, w, []byte(errText), http.StatusBadRequest)
		return
	}
	err = json.Unmarshal(rBody, userUpdateDTO)
	if err != nil {
		errText := fmt.Sprintf(`{"message": "error in decoding user: %s"}`, err)
		uh.logger.Errorf(errText)
		writeResponse(uh.logger, w, []byte(errText), http.StatusBadRequest)
		return
	}

	if validationErrors := userUpdateDTO.Validate(); len(validationErrors) != 0 {
		var errorsJSON []byte
		errorsJSON, err = json.Marshal(validationErrors)
		if err != nil {
			errText := fmt.Sprintf(`{"message": "error in json decoding: %s"}`, err)
			uh.logger.Errorf(errText)
			writeResponse(uh.logger, w, []byte(errText), http.StatusInternalServerError)
			return
		}
		writeResponse(uh.logger, w, errorsJSON, http.StatusUnprocessableEntity)
		return
	}

	user := userUpdateDTO.ConvertToUser()
	updatedUser, err := uh.u.UpdateUserUseCase(user)
	if err != nil {
		errText := fmt.Sprintf(`{"message": "internal server error"}`)
		uh.logger.Errorf(errText)
		writeResponse(uh.logger, w, []byte(errText), http.StatusInternalServerError)
		return
	}
	if updatedUser == nil {
		errText := fmt.Sprintf(`{"message": "user with ID %d is not found"}`, user.ID)
		uh.logger.Errorf(errText)
		writeResponse(uh.logger, w, []byte(errText), http.StatusNotFound)
		return
	}
	userJSON, err := json.Marshal(updatedUser)
	if err != nil {
		errText := fmt.Sprintf(`{"message": "error in coding user: %s"}`, err)
		uh.logger.Errorf(errText)
		writeResponse(uh.logger, w, []byte(errText), http.StatusInternalServerError)
		return
	}
	writeResponse(uh.logger, w, userJSON, http.StatusOK)
}

func (uh *UserHandler) GetUserByIDHandlerID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userID := vars["USER_ID"]
	userIDInt, err := strconv.ParseInt(userID, 10, 64)
	if err != nil {
		errText := fmt.Sprintf(`{"message": "bad format of user id: %s"}`, err)
		uh.logger.Errorf(errText)
		writeResponse(uh.logger, w, []byte(errText), http.StatusBadRequest)
		return
	}
	user, err := uh.u.GetUserByIDUseCase(userIDInt)
	if err != nil {
		errText := fmt.Sprintf(`{"message": "internal server error"}`)
		uh.logger.Errorf(errText)
		writeResponse(uh.logger, w, []byte(errText), http.StatusInternalServerError)
		return
	}
	if user == nil {
		errText := fmt.Sprintf(`{"message": "user with ID %d is not found"}`, userIDInt)
		uh.logger.Errorf(errText)
		writeResponse(uh.logger, w, []byte(errText), http.StatusNotFound)
		return
	}

	userJSON, err := json.Marshal(user)
	if err != nil {
		errText := fmt.Sprintf(`{"message": "error in coding user: %s"}`, err)
		uh.logger.Errorf(errText)
		writeResponse(uh.logger, w, []byte(errText), http.StatusInternalServerError)
		return
	}
	writeResponse(uh.logger, w, userJSON, http.StatusOK)
}

func (uh *UserHandler) SearchUsersHandler(w http.ResponseWriter, r *http.Request) {
	filter, err := parseFilterFromRequest(r)
	if err != nil {
		errText := fmt.Sprintf(`{"message": "bad filtering params: %s"}`, err)
		uh.logger.Errorf(errText)
		writeResponse(uh.logger, w, []byte(errText), http.StatusBadRequest)
		return
	}
	users, err := uh.u.SearchUsersUseCase(filter)
	if err != nil {
		errText := fmt.Sprintf(`{"message": "internal server error"}`)
		uh.logger.Errorf(errText)
		writeResponse(uh.logger, w, []byte(errText), http.StatusInternalServerError)
		return
	}
	if users == nil {
		errText := fmt.Sprintf(`{"message": "users are not found"}`)
		uh.logger.Errorf(errText)
		writeResponse(uh.logger, w, []byte(errText), http.StatusNotFound)
		return
	}
	userJSON, err := json.Marshal(users)
	if err != nil {
		errText := fmt.Sprintf(`{"message": "error in coding users: %s"}`, err)
		uh.logger.Errorf(errText)
		writeResponse(uh.logger, w, []byte(errText), http.StatusInternalServerError)
		return
	}
	writeResponse(uh.logger, w, userJSON, http.StatusOK)
}

func parseFilterFromRequest(r *http.Request) (filters.Filter, error) {
	var filter filters.Filter
	params := r.URL.Query()

	filter.Gender = params.Get("Gender")
	if filter.Gender != "" {
		if filter.Gender != "male" && filter.Gender != "female" {
			return filter, fmt.Errorf("gender myst be male or female")
		}
	}
	filter.Status = params.Get("Status")
	if filter.Status != "" {
		if filter.Status != "active" && filter.Status != "banned" && filter.Status != "deleted" {
			return filter, fmt.Errorf("unknown statuses")
		}
	}
	filter.FullName = params.Get("FullName")

	filter.SortAsk, _ = strconv.ParseBool(params.Get("SortAsk"))
	filter.SortDesc, _ = strconv.ParseBool(params.Get("SortDesc"))

	if filter.SortAsk && filter.SortDesc {
		return filter, fmt.Errorf("you can not sort ask and desc at the same time")
	}

	attributesToSort := params.Get("AttributesToSort")
	filter.AttributesToSort = attributesToSort

	if filter.SortAsk || filter.SortDesc {
		if filter.AttributesToSort == "" {
			return filter, fmt.Errorf("sorting param is not set")

		}
		switch filter.AttributesToSort {
		case "id":
			break
		case "name":
			break
		case "surname":
			break
		case "patronymic":
			break
		case "gender":
			break
		case "status":
			break
		case "birthday":
			break
		case "join_date":
			break
		default:
			return filter, fmt.Errorf("unknown sorting param")
		}
	}

	limitStr := params.Get("Limit")
	if limit, err := strconv.Atoi(limitStr); err == nil {
		filter.Limit = limit
	}

	offsetStr := params.Get("Offset")
	if offset, err := strconv.Atoi(offsetStr); err == nil {
		filter.Offset = offset
	}

	return filter, nil
}

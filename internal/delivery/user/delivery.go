package user

import (
	"database/sql"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	ent "tender-workspace/internal/entity"
	"tender-workspace/internal/entity/dto"
	ucUser "tender-workspace/internal/usecase/user"
	f "tender-workspace/internal/utils/functions"
	mc "tender-workspace/internal/utils/myconstants"
	"tender-workspace/internal/utils/myerrors"
	e "tender-workspace/internal/utils/myerrors"

	"github.com/gorilla/mux"
	"go.uber.org/zap"
)

type DeliveryLayer struct {
	ucUser ucUser.Usecase
	logger *zap.Logger
}

func NewDeliveryLayer(ucUser ucUser.Usecase, logger *zap.Logger) *DeliveryLayer {
	return &DeliveryLayer{
		ucUser: ucUser,
		logger: logger,
	}
}

func (d *DeliveryLayer) CreateUser(w http.ResponseWriter, r *http.Request) {
	requestId := r.Context().Value(mc.ContextKey(mc.RequestID)).(string)
	if r.Method != "POST" {
		d.logger.Info(e.ErrMethodNotAllowed.Error(), zap.String(mc.RequestID, requestId))
		propsError := f.NewResponseProps(w, ent.ResponseReason{Reason: e.ErrMethodNotAllowed.Error()}, http.StatusMethodNotAllowed, mc.ApplicationJson)
		f.Response(propsError)
		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		d.logger.Info(err.Error(), zap.String(mc.RequestID, requestId))
		propsError := f.NewResponseProps(w, ent.ResponseReason{Reason: e.ErrRequestBody.Error()}, http.StatusBadRequest, mc.ApplicationJson)
		f.Response(propsError)
		return
	}
	var userData dto.UserInput
	err = json.Unmarshal(body, &userData)
	if err != nil {
		d.logger.Info(err.Error(), zap.String(mc.RequestID, requestId))
		propsError := f.NewResponseProps(w, ent.ResponseReason{Reason: e.ErrRequestBody.Error()}, http.StatusBadRequest, mc.ApplicationJson)
		f.Response(propsError)
		return
	}
	isValid, err := f.Validate(userData)
	if err != nil || !isValid {
		d.logger.Info(err.Error(), zap.String(mc.RequestID, requestId))
		propsError := f.NewResponseProps(w, ent.ResponseReason{Reason: e.ErrRequestBody.Error()}, http.StatusBadRequest, mc.ApplicationJson)
		f.Response(propsError)
		return
	}

	user, err := d.ucUser.Create(r.Context(), &userData)
	if err != nil {
		d.logger.Info(err.Error(), zap.String(mc.RequestID, requestId))
		if errors.Is(err, myerrors.ErrUserAlreadyExist) {
			propsError := f.NewResponseProps(w, ent.ResponseError{Error: err.Error()}, http.StatusUnauthorized, mc.ApplicationJson)
			f.Response(propsError)
			return
		}
		propsError := f.NewResponseProps(w, ent.ResponseError{Error: e.ErrInternal.Error()}, http.StatusInternalServerError, mc.ApplicationJson)
		f.Response(propsError)
		return
	}

	userOutput := newUserOutput(user)
	responseData := f.NewResponseProps(w, userOutput, http.StatusOK, mc.ApplicationJson)
	f.Response(responseData)
}

func (d *DeliveryLayer) GetUser(w http.ResponseWriter, r *http.Request) {
	requestId := r.Context().Value(mc.ContextKey(mc.RequestID)).(string)
	if r.Method != "GET" {
		d.logger.Info(e.ErrMethodNotAllowed.Error(), zap.String(mc.RequestID, requestId))
		propsError := f.NewResponseProps(w, ent.ResponseReason{Reason: e.ErrMethodNotAllowed.Error()}, http.StatusMethodNotAllowed, mc.ApplicationJson)
		f.Response(propsError)
		return
	}
	username := mux.Vars(r)["username"]
	user, err := d.ucUser.GetData(r.Context(), username)
	if err != nil {
		d.logger.Info(err.Error(), zap.String(mc.RequestID, requestId))
		if errors.Is(err, sql.ErrNoRows) {
			propsError := f.NewResponseProps(w, ent.ResponseReason{Reason: e.ErrUserNotExist.Error()},
				http.StatusNotFound, mc.ApplicationJson)
			f.Response(propsError)
			return
		}
		propsError := f.NewResponseProps(w, ent.ResponseError{Error: e.ErrInternal.Error()},
			http.StatusInternalServerError, mc.ApplicationJson)
		f.Response(propsError)
		return
	}
	userOutput := newUserOutput(user)
	responseData := f.NewResponseProps(w, userOutput, http.StatusOK, mc.ApplicationJson)
	f.Response(responseData)
}

func (d *DeliveryLayer) GetUserOrganizationsIds(w http.ResponseWriter, r *http.Request) {
	requestId := r.Context().Value(mc.ContextKey(mc.RequestID)).(string)
	if r.Method != "GET" {
		d.logger.Info(e.ErrMethodNotAllowed.Error(), zap.String(mc.RequestID, requestId))
		propsError := f.NewResponseProps(w, ent.ResponseReason{Reason: e.ErrMethodNotAllowed.Error()}, http.StatusMethodNotAllowed, mc.ApplicationJson)
		f.Response(propsError)
		return
	}
	username := mux.Vars(r)["username"]
	ids, err := d.ucUser.GetUserOrganizations(r.Context(), username)
	if err != nil {
		d.logger.Info(err.Error(), zap.String(mc.RequestID, requestId))
		if errors.Is(err, e.ErrUserExist) {
			propsError := f.NewResponseProps(w, ent.ResponseReason{Reason: e.ErrUserNotExist.Error()},
				http.StatusUnauthorized, mc.ApplicationJson)
			f.Response(propsError)
			return
		}
		propsError := f.NewResponseProps(w, ent.ResponseError{Error: e.ErrInternal.Error()},
			http.StatusInternalServerError, mc.ApplicationJson)
		f.Response(propsError)
		return
	}
	if ids == nil {
		ids = make([]int, 0)
	}
	responseData := f.NewResponseProps(w, ids, http.StatusOK, mc.ApplicationJson)
	f.Response(responseData)
}

package organization

import (
	"database/sql"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"strconv"
	ent "tender-workspace/internal/entity"
	"tender-workspace/internal/entity/dto"
	oqp "tender-workspace/internal/entity/dto/queries/organizations"
	usecaseOrg "tender-workspace/internal/usecase/organization"
	f "tender-workspace/internal/utils/functions"
	mc "tender-workspace/internal/utils/myconstants"
	e "tender-workspace/internal/utils/myerrors"

	"github.com/gorilla/mux"
	"go.uber.org/zap"
)

type DeliveryLayer struct {
	ucOrganizaiton usecaseOrg.Usecase
	logger         *zap.Logger
}

func NewDeliveryLayer(ucOrg usecaseOrg.Usecase, logger *zap.Logger) *DeliveryLayer {
	return &DeliveryLayer{
		ucOrganizaiton: ucOrg,
		logger:         logger,
	}
}

func (d *DeliveryLayer) GetListOfOrganizations(w http.ResponseWriter, r *http.Request) {
	requestId := r.Context().Value(mc.ContextKey(mc.RequestID)).(string)
	if r.Method != "GET" {
		d.logger.Info(e.ErrMethodNotAllowed.Error(), zap.String(mc.RequestID, requestId))
		propsError := f.NewResponseProps(w, ent.ResponseReason{Reason: e.ErrMethodNotAllowed.Error()}, http.StatusMethodNotAllowed, mc.ApplicationJson)
		f.Response(propsError)
		return
	}

	queryParams := new(oqp.OrganizationList)
	err := queryParams.GetParameters(r)
	if err != nil {
		d.logger.Info(err.Error(), zap.String(mc.RequestID, requestId))
		if errors.Is(err, e.ErrQPLimit) ||
			errors.Is(err, e.ErrQPOffset) ||
			errors.Is(err, e.ErrQPOrgType) {
			propsError := f.NewResponseProps(w, ent.ResponseReason{Reason: err.Error()}, http.StatusBadRequest, mc.ApplicationJson)
			f.Response(propsError)
			return
		}
		propsError := f.NewResponseProps(w, ent.ResponseError{Error: e.ErrInternal.Error()},
			http.StatusInternalServerError, mc.ApplicationJson)
		f.Response(propsError)
		return
	}

	orgs, err := d.ucOrganizaiton.GetAll(r.Context(), queryParams)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		d.logger.Info(err.Error(), zap.String(mc.RequestID, requestId))
		propsError := f.NewResponseProps(w, ent.ResponseError{Error: e.ErrInternal.Error()},
			http.StatusInternalServerError, mc.ApplicationJson)
		f.Response(propsError)
		return
	}
	if orgs == nil {
		orgs = make([]*ent.Organization, 0)
	}

	orgsOutput := newArrayOrgOutput(orgs)
	responseData := f.NewResponseProps(w, orgsOutput, http.StatusOK, mc.ApplicationJson)
	f.Response(responseData)
}

func (d *DeliveryLayer) CreateNewOrganization(w http.ResponseWriter, r *http.Request) {
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
	var orgData dto.OrganizationInput
	err = json.Unmarshal(body, &orgData)
	if err != nil {
		d.logger.Info(err.Error(), zap.String(mc.RequestID, requestId))
		propsError := f.NewResponseProps(w, ent.ResponseReason{Reason: e.ErrRequestBody.Error()}, http.StatusBadRequest, mc.ApplicationJson)
		f.Response(propsError)
		return
	}
	isValid, err := f.Validate(orgData)
	if err != nil || !isValid {
		d.logger.Info(err.Error(), zap.String(mc.RequestID, requestId))
		propsError := f.NewResponseProps(w, ent.ResponseReason{Reason: e.ErrRequestBody.Error()}, http.StatusBadRequest, mc.ApplicationJson)
		f.Response(propsError)
		return
	}

	org, err := d.ucOrganizaiton.Create(r.Context(), &orgData)
	if err != nil {
		d.logger.Info(err.Error(), zap.String(mc.RequestID, requestId))
		if errors.Is(err, e.ErrQPOrgType) {
			propsError := f.NewResponseProps(w, ent.ResponseReason{Reason: err.Error()}, http.StatusBadRequest, mc.ApplicationJson)
			f.Response(propsError)
			return
		}
		propsError := f.NewResponseProps(w, ent.ResponseError{Error: e.ErrInternal.Error()}, http.StatusInternalServerError, mc.ApplicationJson)
		f.Response(propsError)
		return
	}

	orgOutput := newOrganizationOutput(org)
	responseData := f.NewResponseProps(w, orgOutput, http.StatusOK, mc.ApplicationJson)
	f.Response(responseData)
}

func (d *DeliveryLayer) UpdateOrganization(w http.ResponseWriter, r *http.Request) {
	requestId := r.Context().Value(mc.ContextKey(mc.RequestID)).(string)
	if r.Method != "PUT" {
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
	var orgData dto.OrganizationInput
	err = json.Unmarshal(body, &orgData)
	if err != nil {
		d.logger.Info(err.Error(), zap.String(mc.RequestID, requestId))
		propsError := f.NewResponseProps(w, ent.ResponseReason{Reason: e.ErrRequestBody.Error()}, http.StatusBadRequest, mc.ApplicationJson)
		f.Response(propsError)
		return
	}
	isValid, err := f.Validate(orgData)
	if err != nil || !isValid {
		d.logger.Info(err.Error(), zap.String(mc.RequestID, requestId))
		propsError := f.NewResponseProps(w, ent.ResponseReason{Reason: e.ErrRequestBody.Error()}, http.StatusBadRequest, mc.ApplicationJson)
		f.Response(propsError)
		return
	}
	organizationID, err := strconv.Atoi(mux.Vars(r)["organizationID"])
	if err != nil || organizationID < 1 {
		d.logger.Info(e.ErrRequestBody.Error(), zap.String(mc.RequestID, requestId))
		propsError := f.NewResponseProps(w, ent.ResponseReason{Reason: e.ErrRequestBody.Error()}, http.StatusBadRequest, mc.ApplicationJson)
		f.Response(propsError)
		return
	}
	org, err := d.ucOrganizaiton.Update(r.Context(), &orgData, organizationID)
	if err != nil {
		d.logger.Info(err.Error(), zap.String(mc.RequestID, requestId))
		if errors.Is(err, e.ErrQPOrgType) {
			propsError := f.NewResponseProps(w, ent.ResponseReason{Reason: err.Error()}, http.StatusBadRequest, mc.ApplicationJson)
			f.Response(propsError)
			return
		}
		propsError := f.NewResponseProps(w, ent.ResponseError{Error: e.ErrInternal.Error()}, http.StatusInternalServerError, mc.ApplicationJson)
		f.Response(propsError)
		return
	}

	orgOutput := newOrganizationOutput(org)
	responseData := f.NewResponseProps(w, orgOutput, http.StatusOK, mc.ApplicationJson)
	f.Response(responseData)
}

func (d *DeliveryLayer) MakeResponsible(w http.ResponseWriter, r *http.Request) {
	requestId := r.Context().Value(mc.ContextKey(mc.RequestID)).(string)
	if r.Method != "POST" {
		d.logger.Info(e.ErrMethodNotAllowed.Error(), zap.String(mc.RequestID, requestId))
		propsError := f.NewResponseProps(w, ent.ResponseReason{Reason: e.ErrMethodNotAllowed.Error()}, http.StatusMethodNotAllowed, mc.ApplicationJson)
		f.Response(propsError)
		return
	}

	organizationID, err := strconv.Atoi(mux.Vars(r)["organizationID"])
	if err != nil || organizationID < 1 {
		d.logger.Info(e.ErrRequestBody.Error(), zap.String(mc.RequestID, requestId))
		propsError := f.NewResponseProps(w, ent.ResponseReason{Reason: e.ErrRequestBody.Error()}, http.StatusBadRequest, mc.ApplicationJson)
		f.Response(propsError)
		return
	}
	username := mux.Vars(r)["username"]
	err = d.ucOrganizaiton.MakeResponsible(r.Context(), username, organizationID)
	if err != nil {
		d.logger.Info(err.Error(), zap.String(mc.RequestID, requestId))
		if errors.Is(err, e.ErrUserExist) {
			propsError := f.NewResponseProps(w, ent.ResponseReason{Reason: err.Error()}, http.StatusUnauthorized, mc.ApplicationJson)
			f.Response(propsError)
			return
		}
		if errors.Is(err, e.ErrOrganizationExist) || errors.Is(err, e.ErrUserAlreadyResponsible) {
			propsError := f.NewResponseProps(w, ent.ResponseReason{Reason: err.Error()}, http.StatusBadRequest, mc.ApplicationJson)
			f.Response(propsError)
			return
		}
		propsError := f.NewResponseProps(w, ent.ResponseError{Error: e.ErrInternal.Error()}, http.StatusInternalServerError, mc.ApplicationJson)
		f.Response(propsError)
		return
	}

	responseData := f.NewResponseProps(w, ent.ResponseDetail{Detail: "ok"}, http.StatusOK, mc.ApplicationJson)
	f.Response(responseData)
}

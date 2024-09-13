package tender

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	ent "tender-workspace/internal/entity"
	"tender-workspace/internal/entity/dto"
	tqp "tender-workspace/internal/entity/dto/queries/tenders"
	"tender-workspace/internal/usecase/tender"
	f "tender-workspace/internal/utils/functions"
	mc "tender-workspace/internal/utils/myconstants"
	e "tender-workspace/internal/utils/myerrors"

	"go.uber.org/zap"
)

type DeliveryLayer struct {
	ucTender tender.Usecase
	logger   *zap.Logger
}

func NewDeliveryLayer(ucTender tender.Usecase, logger *zap.Logger) *DeliveryLayer {
	return &DeliveryLayer{
		ucTender: ucTender,
		logger:   logger,
	}
}

func (d *DeliveryLayer) GetListOfTenders(w http.ResponseWriter, r *http.Request) {
	requestId := r.Context().Value(mc.ContextKey(mc.RequestID)).(string)
	if r.Method != "GET" {
		d.logger.Info(e.ErrMethodNotAllowed.Error(), zap.String(mc.RequestID, requestId))
		propsError := f.NewResponseProps(w, ent.ResponseReason{Reason: e.ErrMethodNotAllowed.Error()}, http.StatusMethodNotAllowed, mc.ApplicationJson)
		f.Response(propsError)
		return
	}
	queryParams := new(tqp.ListTenders)
	err := queryParams.GetParameters(r)
	if err != nil {
		d.logger.Info(err.Error(), zap.String(mc.RequestID, requestId))
		if errors.Is(err, e.ErrQPLimit) || errors.Is(err, e.ErrQPOffset) || errors.Is(err, e.ErrQPServiceType) {
			propsError := f.NewResponseProps(w, ent.ResponseReason{Reason: err.Error()}, http.StatusBadRequest, mc.ApplicationJson)
			f.Response(propsError)
			return
		}
		propsError := f.NewResponseProps(w, ent.ResponseError{Error: e.ErrInternal.Error()},
			http.StatusInternalServerError, mc.ApplicationJson)
		f.Response(propsError)
		return
	}

	tenders, err := d.ucTender.GetTenders(r.Context(), queryParams)
	if err != nil {
		d.logger.Info(err.Error(), zap.String(mc.RequestID, requestId))
		propsError := f.NewResponseProps(w, ent.ResponseError{Error: e.ErrInternal.Error()},
			http.StatusInternalServerError, mc.ApplicationJson)
		f.Response(propsError)
		return
	}
	if tenders == nil {
		tenders = make([]*ent.Tender, 0)
	}

	tendersOutput := dto.NewArrayTenderOutput(tenders)
	responseData := f.NewResponseProps(w, tendersOutput, http.StatusOK, mc.ApplicationJson)
	f.Response(responseData)
}

func (d *DeliveryLayer) CreateNewTender(w http.ResponseWriter, r *http.Request) {
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
	var tenderData dto.TenderInput
	err = json.Unmarshal(body, &tenderData)
	if err != nil {
		d.logger.Info(err.Error(), zap.String(mc.RequestID, requestId))
		propsError := f.NewResponseProps(w, ent.ResponseReason{Reason: e.ErrRequestBody.Error()}, http.StatusBadRequest, mc.ApplicationJson)
		f.Response(propsError)
		return
	}
	isValid, err := f.Validate(tenderData)
	if err != nil || !isValid {
		d.logger.Info(err.Error(), zap.String(mc.RequestID, requestId))
		propsError := f.NewResponseProps(w, ent.ResponseReason{Reason: e.ErrRequestBody.Error()}, http.StatusBadRequest, mc.ApplicationJson)
		f.Response(propsError)
		return
	}

	tender, err := d.ucTender.CreateTender(r.Context(), &tenderData)
	if err != nil {
		d.logger.Info(err.Error(), zap.String(mc.RequestID, requestId))
		if errors.Is(err, e.ErrBadStatusCreate) || errors.Is(err, e.ErrQPServiceType) || errors.Is(err, e.ErrOrganizationExist) {
			propsError := f.NewResponseProps(w, ent.ResponseReason{Reason: err.Error()}, http.StatusBadRequest, mc.ApplicationJson)
			f.Response(propsError)
			return
		}
		if errors.Is(err, e.ErrResponsibilty) || errors.Is(err, e.ErrUserExist) {
			propsError := f.NewResponseProps(w, ent.ResponseReason{Reason: err.Error()}, http.StatusUnauthorized, mc.ApplicationJson)
			f.Response(propsError)
			return
		}
		propsError := f.NewResponseProps(w, ent.ResponseError{Error: e.ErrInternal.Error()}, http.StatusInternalServerError, mc.ApplicationJson)
		f.Response(propsError)
		return
	}

	tenderOutput := dto.NewTenderOutput(tender)
	responseData := f.NewResponseProps(w, tenderOutput, http.StatusOK, mc.ApplicationJson)
	f.Response(responseData)
}

func (d *DeliveryLayer) GetUserTenders(w http.ResponseWriter, r *http.Request) {
	requestId := r.Context().Value(mc.ContextKey(mc.RequestID)).(string)
	if r.Method != "GET" {
		d.logger.Info(e.ErrMethodNotAllowed.Error(), zap.String(mc.RequestID, requestId))
		propsError := f.NewResponseProps(w, ent.ResponseReason{Reason: e.ErrMethodNotAllowed.Error()}, http.StatusMethodNotAllowed, mc.ApplicationJson)
		f.Response(propsError)
		return
	}

	queryParams := new(tqp.ListUserTenders)
	err := queryParams.GetParameters(r)
	if err != nil {
		d.logger.Info(err.Error(), zap.String(mc.RequestID, requestId))
		if errors.Is(err, e.ErrQPLimit) || errors.Is(err, e.ErrQPOffset) {
			propsError := f.NewResponseProps(w, ent.ResponseReason{Reason: err.Error()}, http.StatusBadRequest, mc.ApplicationJson)
			f.Response(propsError)
			return
		}
		if errors.Is(err, e.ErrBadPermission) {
			propsError := f.NewResponseProps(w, ent.ResponseReason{Reason: err.Error()}, http.StatusUnauthorized, mc.ApplicationJson)
			f.Response(propsError)
			return
		}
		propsError := f.NewResponseProps(w, ent.ResponseError{Error: e.ErrInternal.Error()},
			http.StatusInternalServerError, mc.ApplicationJson)
		f.Response(propsError)
		return
	}

	tenders, err := d.ucTender.GetUserTenders(r.Context(), queryParams)
	if err != nil {
		d.logger.Info(err.Error(), zap.String(mc.RequestID, requestId))
		if errors.Is(err, e.ErrUserIsNotResponsible) || errors.Is(err, e.ErrUserExist) {
			propsError := f.NewResponseProps(w, ent.ResponseReason{Reason: err.Error()}, http.StatusUnauthorized, mc.ApplicationJson)
			f.Response(propsError)
			return
		}
		if errors.Is(err, e.ErrBigInterval) {
			propsError := f.NewResponseProps(w, ent.ResponseReason{Reason: err.Error()}, http.StatusBadRequest, mc.ApplicationJson)
			f.Response(propsError)
			return
		}
		propsError := f.NewResponseProps(w, ent.ResponseError{Error: e.ErrInternal.Error()},
			http.StatusInternalServerError, mc.ApplicationJson)
		f.Response(propsError)
		return
	}

	tenderOutput := dto.NewArrayTenderOutput(tenders)
	responseData := f.NewResponseProps(w, tenderOutput, http.StatusOK, mc.ApplicationJson)
	f.Response(responseData)
}

func (d *DeliveryLayer) GetTenderStatus(w http.ResponseWriter, r *http.Request) {
	requestId := r.Context().Value(mc.ContextKey(mc.RequestID)).(string)
	if r.Method != "GET" {
		d.logger.Info(e.ErrMethodNotAllowed.Error(), zap.String(mc.RequestID, requestId))
		propsError := f.NewResponseProps(w, ent.ResponseReason{Reason: e.ErrMethodNotAllowed.Error()}, http.StatusMethodNotAllowed, mc.ApplicationJson)
		f.Response(propsError)
		return
	}

	queryParams := new(tqp.TenderStatus)
	err := queryParams.GetParameters(r)
	if err != nil {
		d.logger.Info(err.Error(), zap.String(mc.RequestID, requestId))
		if errors.Is(err, e.ErrExistTenderID) || errors.Is(err, e.ErrTenderID) {
			propsError := f.NewResponseProps(w, ent.ResponseReason{Reason: err.Error()}, http.StatusBadRequest, mc.ApplicationJson)
			f.Response(propsError)
			return
		}
		if errors.Is(err, e.ErrBadPermission) {
			propsError := f.NewResponseProps(w, ent.ResponseReason{Reason: err.Error()}, http.StatusUnauthorized, mc.ApplicationJson)
			f.Response(propsError)
			return
		}
		propsError := f.NewResponseProps(w, ent.ResponseError{Error: e.ErrInternal.Error()},
			http.StatusInternalServerError, mc.ApplicationJson)
		f.Response(propsError)
		return
	}

	status, err := d.ucTender.GetTenderStatus(r.Context(), queryParams)
	if err != nil {
		d.logger.Info(err.Error(), zap.String(mc.RequestID, requestId))
		if errors.Is(err, e.ErrUserExist) || errors.Is(err, e.ErrResponsibilty) {
			propsError := f.NewResponseProps(w, ent.ResponseReason{Reason: err.Error()}, http.StatusUnauthorized, mc.ApplicationJson)
			f.Response(propsError)
			return
		}
		if errors.Is(err, e.ErrNoTenders) {
			propsError := f.NewResponseProps(w, ent.ResponseReason{Reason: err.Error()}, http.StatusBadRequest, mc.ApplicationJson)
			f.Response(propsError)
			return
		}
		propsError := f.NewResponseProps(w, ent.ResponseError{Error: e.ErrInternal.Error()},
			http.StatusInternalServerError, mc.ApplicationJson)
		f.Response(propsError)
		return
	}

	responseData := f.NewResponseProps(w, dto.TenderStatus{Status: status}, http.StatusOK, mc.ApplicationJson)
	f.Response(responseData)
}

func (d *DeliveryLayer) UpdateTenderStatus(w http.ResponseWriter, r *http.Request) {
	requestId := r.Context().Value(mc.ContextKey(mc.RequestID)).(string)
	if r.Method != "PUT" {
		d.logger.Info(e.ErrMethodNotAllowed.Error(), zap.String(mc.RequestID, requestId))
		propsError := f.NewResponseProps(w, ent.ResponseReason{Reason: e.ErrMethodNotAllowed.Error()}, http.StatusMethodNotAllowed, mc.ApplicationJson)
		f.Response(propsError)
		return
	}

	queryParams := new(tqp.UpdateTenderStatus)
	err := queryParams.GetParameters(r)
	if err != nil {
		d.logger.Info(err.Error(), zap.String(mc.RequestID, requestId))
		if errors.Is(err, e.ErrExistTenderID) ||
			errors.Is(err, e.ErrTenderID) ||
			errors.Is(err, e.ErrTenderStatus) ||
			errors.Is(err, e.ErrQPChangeStatus) {
			propsError := f.NewResponseProps(w, ent.ResponseReason{Reason: err.Error()}, http.StatusBadRequest, mc.ApplicationJson)
			f.Response(propsError)
			return
		}
		if errors.Is(err, e.ErrBadPermission) {
			propsError := f.NewResponseProps(w, ent.ResponseReason{Reason: err.Error()}, http.StatusUnauthorized, mc.ApplicationJson)
			f.Response(propsError)
			return
		}
		propsError := f.NewResponseProps(w, ent.ResponseError{Error: e.ErrInternal.Error()},
			http.StatusInternalServerError, mc.ApplicationJson)
		f.Response(propsError)
		return
	}

	tender, err := d.ucTender.UpdateTenderStatus(r.Context(), queryParams)
	if err != nil {
		d.logger.Info(err.Error(), zap.String(mc.RequestID, requestId))
		if errors.Is(err, e.ErrUserExist) || errors.Is(err, e.ErrResponsibilty) {
			propsError := f.NewResponseProps(w, ent.ResponseReason{Reason: err.Error()}, http.StatusUnauthorized, mc.ApplicationJson)
			f.Response(propsError)
			return
		}
		if errors.Is(err, e.ErrNoTenders) {
			propsError := f.NewResponseProps(w, ent.ResponseReason{Reason: err.Error()}, http.StatusBadRequest, mc.ApplicationJson)
			f.Response(propsError)
			return
		}
		propsError := f.NewResponseProps(w, ent.ResponseError{Error: e.ErrInternal.Error()},
			http.StatusInternalServerError, mc.ApplicationJson)
		f.Response(propsError)
		return
	}

	tenderOutput := dto.NewTenderOutput(tender)
	responseData := f.NewResponseProps(w, tenderOutput, http.StatusOK, mc.ApplicationJson)
	f.Response(responseData)
}

func (d *DeliveryLayer) UpdateTender(w http.ResponseWriter, r *http.Request) {
	requestId := r.Context().Value(mc.ContextKey(mc.RequestID)).(string)
	if r.Method != "PATCH" {
		d.logger.Info(e.ErrMethodNotAllowed.Error(), zap.String(mc.RequestID, requestId))
		propsError := f.NewResponseProps(w, ent.ResponseReason{Reason: e.ErrMethodNotAllowed.Error()}, http.StatusMethodNotAllowed, mc.ApplicationJson)
		f.Response(propsError)
		return
	}

	queryParams := new(tqp.TenderUpdate)
	err := queryParams.GetParameters(r)
	if err != nil {
		d.logger.Info(err.Error(), zap.String(mc.RequestID, requestId))
		if errors.Is(err, e.ErrExistTenderID) || errors.Is(err, e.ErrTenderID) {
			propsError := f.NewResponseProps(w, ent.ResponseReason{Reason: err.Error()}, http.StatusBadRequest, mc.ApplicationJson)
			f.Response(propsError)
			return
		}
		if errors.Is(err, e.ErrBadPermission) {
			propsError := f.NewResponseProps(w, ent.ResponseReason{Reason: err.Error()}, http.StatusUnauthorized, mc.ApplicationJson)
			f.Response(propsError)
			return
		}
		propsError := f.NewResponseProps(w, ent.ResponseError{Error: e.ErrInternal.Error()},
			http.StatusInternalServerError, mc.ApplicationJson)
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
	var tenderData dto.TenderUpdateDataInput
	err = json.Unmarshal(body, &tenderData)
	if err != nil {
		d.logger.Info(err.Error(), zap.String(mc.RequestID, requestId))
		propsError := f.NewResponseProps(w, ent.ResponseReason{Reason: e.ErrRequestBody.Error()}, http.StatusBadRequest, mc.ApplicationJson)
		f.Response(propsError)
		return
	}
	isValid, err := f.Validate(tenderData)
	if err != nil || !isValid {
		d.logger.Info(err.Error(), zap.String(mc.RequestID, requestId))
		propsError := f.NewResponseProps(w, ent.ResponseReason{Reason: e.ErrRequestBody.Error()}, http.StatusBadRequest, mc.ApplicationJson)
		f.Response(propsError)
		return
	}
	tender, err := d.ucTender.UpdateTender(r.Context(), &tenderData, queryParams)
	if err != nil {
		d.logger.Info(err.Error(), zap.String(mc.RequestID, requestId))
		if errors.Is(err, e.ErrUserExist) || errors.Is(err, e.ErrResponsibilty) {
			propsError := f.NewResponseProps(w, ent.ResponseReason{Reason: err.Error()}, http.StatusUnauthorized, mc.ApplicationJson)
			f.Response(propsError)
			return
		}
		if errors.Is(err, e.ErrNoTenders) || errors.Is(err, e.ErrQPServiceType) {
			propsError := f.NewResponseProps(w, ent.ResponseReason{Reason: err.Error()}, http.StatusBadRequest, mc.ApplicationJson)
			f.Response(propsError)
			return
		}
		propsError := f.NewResponseProps(w, ent.ResponseError{Error: e.ErrInternal.Error()},
			http.StatusInternalServerError, mc.ApplicationJson)
		f.Response(propsError)
		return
	}

	tenderOutput := dto.NewTenderOutput(tender)
	responseData := f.NewResponseProps(w, tenderOutput, http.StatusOK, mc.ApplicationJson)
	f.Response(responseData)
}

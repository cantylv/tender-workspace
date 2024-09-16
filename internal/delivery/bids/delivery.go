package bids

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	ent "tender-workspace/internal/entity"
	"tender-workspace/internal/entity/dto"
	bqp "tender-workspace/internal/entity/dto/queries/bids"
	"tender-workspace/internal/usecase/bids"
	f "tender-workspace/internal/utils/functions"
	mc "tender-workspace/internal/utils/myconstants"
	e "tender-workspace/internal/utils/myerrors"

	"go.uber.org/zap"
)

type DeliveryLayer struct {
	ucBids bids.Usecase
	logger *zap.Logger
}

func NewDeliveryLayer(ucBids bids.Usecase, logger *zap.Logger) *DeliveryLayer {
	return &DeliveryLayer{
		ucBids: ucBids,
		logger: logger,
	}
}

func (d *DeliveryLayer) CreateBid(w http.ResponseWriter, r *http.Request) {
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
	var bidData dto.BidInput
	err = json.Unmarshal(body, &bidData)
	if err != nil {
		d.logger.Info(err.Error(), zap.String(mc.RequestID, requestId))
		propsError := f.NewResponseProps(w, ent.ResponseReason{Reason: e.ErrRequestBody.Error()}, http.StatusBadRequest, mc.ApplicationJson)
		f.Response(propsError)
		return
	}
	isValid, err := f.Validate(bidData)
	if err != nil || !isValid {
		d.logger.Info(err.Error(), zap.String(mc.RequestID, requestId))
		propsError := f.NewResponseProps(w, ent.ResponseReason{Reason: e.ErrRequestBody.Error()}, http.StatusBadRequest, mc.ApplicationJson)
		f.Response(propsError)
		return
	}

	bid, err := d.ucBids.CreateBid(r.Context(), &bidData)
	if err != nil {
		d.logger.Info(err.Error(), zap.String(mc.RequestID, requestId))
		if errors.Is(err, e.ErrBadStatusCreate) ||
			errors.Is(err, e.ErrTenderExist) ||
			errors.Is(err, e.ErrOrganizationExist) ||
			errors.Is(err, e.ErrUserAndOrg) ||
			errors.Is(err, e.ErrUserAlreadyHasBid) ||
			errors.Is(err, e.ErrOrgAlreadyHasBid) ||
			errors.Is(err, e.ErrBidYourself) {
			propsError := f.NewResponseProps(w, ent.ResponseReason{Reason: err.Error()}, http.StatusBadRequest, mc.ApplicationJson)
			f.Response(propsError)
			return
		}
		if errors.Is(err, e.ErrUserExist) {
			propsError := f.NewResponseProps(w, ent.ResponseReason{Reason: err.Error()}, http.StatusUnauthorized, mc.ApplicationJson)
			f.Response(propsError)
			return
		}
		propsError := f.NewResponseProps(w, ent.ResponseError{Error: e.ErrInternal.Error()}, http.StatusInternalServerError, mc.ApplicationJson)
		f.Response(propsError)
		return
	}

	tenderOutput := dto.NewBidOutput(bid)
	responseData := f.NewResponseProps(w, tenderOutput, http.StatusOK, mc.ApplicationJson)
	f.Response(responseData)
}

func (d *DeliveryLayer) GetUserBids(w http.ResponseWriter, r *http.Request) {
	requestId := r.Context().Value(mc.ContextKey(mc.RequestID)).(string)
	if r.Method != "GET" {
		d.logger.Info(e.ErrMethodNotAllowed.Error(), zap.String(mc.RequestID, requestId))
		propsError := f.NewResponseProps(w, ent.ResponseReason{Reason: e.ErrMethodNotAllowed.Error()}, http.StatusMethodNotAllowed, mc.ApplicationJson)
		f.Response(propsError)
		return
	}

	queryParams := new(bqp.ListUserBids)
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

	bids, err := d.ucBids.GetUserBids(r.Context(), queryParams)
	if err != nil {
		d.logger.Info(err.Error(), zap.String(mc.RequestID, requestId))
		if errors.Is(err, e.ErrUserExist) {
			propsError := f.NewResponseProps(w, ent.ResponseReason{Reason: err.Error()}, http.StatusUnauthorized, mc.ApplicationJson)
			f.Response(propsError)
			return
		}
		propsError := f.NewResponseProps(w, ent.ResponseError{Error: e.ErrInternal.Error()},
			http.StatusInternalServerError, mc.ApplicationJson)
		f.Response(propsError)
		return
	}
	if len(bids) == 0 {
		bids = make([]*ent.Bid, 0)
	}

	bidsOutput := dto.NewArrayBidOutput(bids)
	responseData := f.NewResponseProps(w, bidsOutput, http.StatusOK, mc.ApplicationJson)
	f.Response(responseData)
}

func (d *DeliveryLayer) GetTenderListOfBids(w http.ResponseWriter, r *http.Request) {
	requestId := r.Context().Value(mc.ContextKey(mc.RequestID)).(string)
	if r.Method != "GET" {
		d.logger.Info(e.ErrMethodNotAllowed.Error(), zap.String(mc.RequestID, requestId))
		propsError := f.NewResponseProps(w, ent.ResponseReason{Reason: e.ErrMethodNotAllowed.Error()}, http.StatusMethodNotAllowed, mc.ApplicationJson)
		f.Response(propsError)
		return
	}

	queryParams := new(bqp.TenderBidList)
	err := queryParams.GetParameters(r)
	if err != nil {
		d.logger.Info(err.Error(), zap.String(mc.RequestID, requestId))
		if errors.Is(err, e.ErrQPLimit) || errors.Is(err, e.ErrQPOffset) || errors.Is(err, e.ErrTenderID) {
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

	bids, err := d.ucBids.GetTenderBids(r.Context(), queryParams)
	if err != nil {
		d.logger.Info(err.Error(), zap.String(mc.RequestID, requestId))
		if errors.Is(err, e.ErrTenderExist) {
			propsError := f.NewResponseProps(w, ent.ResponseReason{Reason: err.Error()}, http.StatusBadRequest, mc.ApplicationJson)
			f.Response(propsError)
			return
		}
		if errors.Is(err, e.ErrUserExist) || errors.Is(err, e.ErrResponsibilty) {
			propsError := f.NewResponseProps(w, ent.ResponseReason{Reason: err.Error()}, http.StatusUnauthorized, mc.ApplicationJson)
			f.Response(propsError)
			return
		}
		propsError := f.NewResponseProps(w, ent.ResponseError{Error: e.ErrInternal.Error()},
			http.StatusInternalServerError, mc.ApplicationJson)
		f.Response(propsError)
		return
	}
	if len(bids) == 0 {
		bids = make([]*ent.Bid, 0)
	}

	bidsOutput := dto.NewArrayBidOutput(bids)
	responseData := f.NewResponseProps(w, bidsOutput, http.StatusOK, mc.ApplicationJson)
	f.Response(responseData)
}

func (d *DeliveryLayer) GetBidStatus(w http.ResponseWriter, r *http.Request) {
	requestId := r.Context().Value(mc.ContextKey(mc.RequestID)).(string)

	queryParams := new(bqp.BidStatus)
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

	status, err := d.ucBids.GetBidStatus(r.Context(), queryParams)
	if err != nil {
		d.logger.Info(err.Error(), zap.String(mc.RequestID, requestId))
		if errors.Is(err, e.ErrUserExist) {
			propsError := f.NewResponseProps(w, ent.ResponseReason{Reason: err.Error()}, http.StatusUnauthorized, mc.ApplicationJson)
			f.Response(propsError)
			return
		}
		if errors.Is(err, e.ErrNoBids) || errors.Is(err, e.ErrBadPermission) {
			propsError := f.NewResponseProps(w, ent.ResponseReason{Reason: err.Error()}, http.StatusBadRequest, mc.ApplicationJson)
			f.Response(propsError)
			return
		}
		propsError := f.NewResponseProps(w, ent.ResponseError{Error: e.ErrInternal.Error()},
			http.StatusInternalServerError, mc.ApplicationJson)
		f.Response(propsError)
		return
	}

	responseData := f.NewResponseProps(w, dto.BidStatus{Status: status}, http.StatusOK, mc.ApplicationJson)
	f.Response(responseData)
}

func (d *DeliveryLayer) UpdateBidStatus(w http.ResponseWriter, r *http.Request) {
	requestId := r.Context().Value(mc.ContextKey(mc.RequestID)).(string)

	queryParams := new(bqp.UpdateBidStatus)
	err := queryParams.GetParameters(r)
	if err != nil {
		d.logger.Info(err.Error(), zap.String(mc.RequestID, requestId))
		if errors.Is(err, e.ErrExistBidID) ||
			errors.Is(err, e.ErrBidID) ||
			errors.Is(err, e.ErrExistStatus) ||
			errors.Is(err, e.ErrQPBidStatusUpdate) {
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

	bid, err := d.ucBids.UpdateBidStatus(r.Context(), queryParams)
	if err != nil {
		d.logger.Info(err.Error(), zap.String(mc.RequestID, requestId))
		if errors.Is(err, e.ErrUserExist) {
			propsError := f.NewResponseProps(w, ent.ResponseReason{Reason: err.Error()}, http.StatusUnauthorized, mc.ApplicationJson)
			f.Response(propsError)
			return
		}
		if errors.Is(err, e.ErrNoBids) || errors.Is(err, e.ErrSetDeprecatedStatus) || errors.Is(err, e.ErrBadPermission) {
			propsError := f.NewResponseProps(w, ent.ResponseReason{Reason: err.Error()}, http.StatusBadRequest, mc.ApplicationJson)
			f.Response(propsError)
			return
		}
		propsError := f.NewResponseProps(w, ent.ResponseError{Error: e.ErrInternal.Error()},
			http.StatusInternalServerError, mc.ApplicationJson)
		f.Response(propsError)
		return
	}

	bidOutput := dto.NewBidOutput(bid)
	responseData := f.NewResponseProps(w, bidOutput, http.StatusOK, mc.ApplicationJson)
	f.Response(responseData)
}

func (d *DeliveryLayer) UpdateBid(w http.ResponseWriter, r *http.Request) {
	requestId := r.Context().Value(mc.RequestID).(string)
	if r.Method != "PATCH" {
		d.logger.Info(e.ErrMethodNotAllowed.Error(), zap.String(mc.RequestID, requestId))
		propsError := f.NewResponseProps(w, ent.ResponseReason{Reason: e.ErrMethodNotAllowed.Error()}, http.StatusMethodNotAllowed, mc.ApplicationJson)
		f.Response(propsError)
		return
	}

	queryParams := new(bqp.UpdateBidData)
	err := queryParams.GetParameters(r)
	if err != nil {
		d.logger.Info(err.Error(), zap.String(mc.RequestID, requestId))
		if errors.Is(err, e.ErrExistBidID) || errors.Is(err, e.ErrBidID) {
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
	var bidData dto.BidUpdateDataInput
	err = json.Unmarshal(body, &bidData)
	if err != nil {
		d.logger.Info(err.Error(), zap.String(mc.RequestID, requestId))
		propsError := f.NewResponseProps(w, ent.ResponseReason{Reason: e.ErrRequestBody.Error()}, http.StatusBadRequest, mc.ApplicationJson)
		f.Response(propsError)
		return
	}
	isValid, err := f.Validate(bidData)
	if err != nil || !isValid {
		d.logger.Info(err.Error(), zap.String(mc.RequestID, requestId))
		propsError := f.NewResponseProps(w, ent.ResponseReason{Reason: e.ErrRequestBody.Error()}, http.StatusBadRequest, mc.ApplicationJson)
		f.Response(propsError)
		return
	}

	bid, err := d.ucBids.UpdateBid(r.Context(), &bidData, queryParams)
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

	bidOutput := dto.NewBidOutput(bid)
	responseData := f.NewResponseProps(w, bidOutput, http.StatusOK, mc.ApplicationJson)
	f.Response(responseData)
}

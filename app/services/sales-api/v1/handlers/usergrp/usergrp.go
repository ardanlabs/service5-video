package usergrp

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/ardanlabs/service/business/core/user"
	"github.com/ardanlabs/service/business/data/page"
	"github.com/ardanlabs/service/business/data/transaction"
	"github.com/ardanlabs/service/business/web/v1/auth"
	"github.com/ardanlabs/service/business/web/v1/response"
	"github.com/ardanlabs/service/foundation/web"
)

// Handlers manages the set of user endpoints.
type Handlers struct {
	user *user.Core
	auth *auth.Auth
}

// New constructs a handlers for route access.
func New(user *user.Core, auth *auth.Auth) *Handlers {
	return &Handlers{
		user: user,
		auth: auth,
	}
}

// Create adds a new user to the system.
func (h *Handlers) Create(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	var app AppNewUser
	if err := web.Decode(r, &app); err != nil {
		return response.NewError(err, http.StatusBadRequest)
	}

	nc, err := toCoreNewUser(app)
	if err != nil {
		return response.NewError(err, http.StatusBadRequest)
	}

	usr, err := h.user.Create(ctx, nc)
	if err != nil {
		if errors.Is(err, user.ErrUniqueEmail) {
			return response.NewError(err, http.StatusConflict)
		}
		return fmt.Errorf("create: usr[%+v]: %w", usr, err)
	}

	return web.Respond(ctx, w, toAppUser(usr), http.StatusCreated)
}

// executeUnderTransaction constructs a new Handlers value with the core apis
// using a store transaction that was created via middleware.
func (h *Handlers) executeUnderTransaction(ctx context.Context) (*Handlers, error) {
	if tx, ok := transaction.Get(ctx); ok {
		user, err := h.user.ExecuteUnderTransaction(tx)
		if err != nil {
			return nil, err
		}

		h = &Handlers{
			user: user,
			auth: h.auth,
		}

		return h, nil
	}

	return h, nil
}

// CreateWithTran adds a new user to the system.
func (h *Handlers) CreateWithTran(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	h, err := h.executeUnderTransaction(ctx)
	if err != nil {
		return err
	}

	var app AppNewUser
	if err := web.Decode(r, &app); err != nil {
		return response.NewError(err, http.StatusBadRequest)
	}

	nc, err := toCoreNewUser(app)
	if err != nil {
		return response.NewError(err, http.StatusBadRequest)
	}

	usr, err := h.user.Create(ctx, nc)
	if err != nil {
		if errors.Is(err, user.ErrUniqueEmail) {
			return response.NewError(err, http.StatusConflict)
		}
		return fmt.Errorf("create: usr[%+v]: %w", usr, err)
	}

	return web.Respond(ctx, w, toAppUser(usr), http.StatusCreated)
}

// Query returns a list of users with paging.
func (h *Handlers) Query(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	page, err := page.Parse(r)
	if err != nil {
		return err
	}

	filter, err := parseFilter(r)
	if err != nil {
		return err
	}

	orderBy, err := parseOrder(r)
	if err != nil {
		return err
	}

	users, err := h.user.Query(ctx, filter, orderBy, page.Number, page.RowsPerPage)
	if err != nil {
		return fmt.Errorf("query: %w", err)
	}

	total, err := h.user.Count(ctx, filter)
	if err != nil {
		return fmt.Errorf("count: %w", err)
	}

	return web.Respond(ctx, w, response.NewPageDocument(toAppUsers(users), total, page.Number, page.RowsPerPage), http.StatusOK)
}

// QueryByID returns a user by its ID.
func (h *Handlers) QueryByID(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	id := auth.GetUserID(ctx)

	usr, err := h.user.QueryByID(ctx, id)
	if err != nil {
		switch {
		case errors.Is(err, user.ErrNotFound):
			return response.NewError(err, http.StatusNotFound)
		default:
			return fmt.Errorf("querybyid: id[%s]: %w", id, err)
		}
	}

	return web.Respond(ctx, w, toAppUser(usr), http.StatusOK)
}

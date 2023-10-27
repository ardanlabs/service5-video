package usergrp

import (
	"net/http"

	"github.com/ardanlabs/service/business/core/user"
	"github.com/ardanlabs/service/business/core/user/stores/userdb"
	db "github.com/ardanlabs/service/business/data/dbsql/pgx"
	"github.com/ardanlabs/service/business/web/v1/auth"
	"github.com/ardanlabs/service/business/web/v1/mid"
	"github.com/ardanlabs/service/foundation/logger"
	"github.com/ardanlabs/service/foundation/web"
	"github.com/jmoiron/sqlx"
)

// Config contains all the mandatory systems required by handlers.
type Config struct {
	Build string
	Log   *logger.Logger
	DB    *sqlx.DB
	Auth  *auth.Auth
}

// Routes adds specific routes for this group.
func Routes(app *web.App, cfg Config) {
	const version = "v1"

	authen := mid.Authenticate(cfg.Auth)
	ruleAdmin := mid.Authorize(cfg.Auth, auth.RuleAdminOnly)
	ruleAdminOrSubject := mid.Authorize(cfg.Auth, auth.RuleAdminOrSubject)
	tran := mid.ExecuteInTransation(cfg.Log, db.NewBeginner(cfg.DB))

	usrCore := user.NewCore(cfg.Log, userdb.NewStore(cfg.Log, cfg.DB))

	hdl := New(usrCore, cfg.Auth)
	app.Handle(http.MethodPost, version, "/users", hdl.Create)
	app.Handle(http.MethodPost, version, "/userstran", hdl.CreateWithTran, authen, ruleAdmin, tran)
	app.Handle(http.MethodPost, version, "/usersauth", hdl.Create, authen, ruleAdmin)
	app.Handle(http.MethodGet, version, "/users", hdl.Query, authen, ruleAdmin)
	app.Handle(http.MethodGet, version, "/users/:user_id", hdl.QueryByID, authen, ruleAdminOrSubject)
}

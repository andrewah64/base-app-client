package session

import (
	"context"
	"fmt"
	"log/slog"
)

import (
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

import (
	"github.com/andrewah64/base-app-client/internal/common/core/db"
)

type AuthSessionUser struct {
	AurId     int
	AurNm     string
	LngCd     string
	RolName   string
	Roles   []string
	LvlNb     int
	EppPt     string
	HrmNm     string
}

func AuthSessionUserInfo(ctx *context.Context, logger *slog.Logger, conn *pgxpool.Conn, tntId int, wauhsSsnTk string, eppPt string, hrmNm string) ([]AuthSessionUser, error){
	const (
		dbSchema = "web_core_auth_ssn_aur_inf"
		dbFunc   = "aur_inf"
	)

	results, err := db.DataSet[AuthSessionUser](ctx, logger, conn,
		func(ctx *context.Context, tx *pgx.Tx)(string, string, *pgx.Rows, error){
			qry := fmt.Sprintf("select %v.%v($1, $2, $3, $4, $5)", dbSchema, dbFunc)

			call, err := (*tx).Query(*ctx, qry, dbFunc, tntId, wauhsSsnTk, eppPt, hrmNm)
			if err != nil {

				slog.LogAttrs(*ctx, slog.LevelError, "AuthSessionUserInfo::Get dataset",
					slog.String("error"      , err.Error()),
					slog.String("qry"        , qry),
					slog.Int   ("tntId"      , tntId),
					slog.String("wauhsSsnTk" , wauhsSsnTk),
					slog.String("eppPt"      , eppPt),
					slog.String("hrmNm"      , hrmNm),
				)

				return qry, dbFunc, nil, fmt.Errorf("AuthSessionUserInfo::Call database function: %w", err)
			}

			return qry, dbFunc, &call, nil
		})

	return results, err
}

type UnauthSessionEndpoint struct {
	LvlNb int
}

func UnauthSessionEndpointInfo(ctx *context.Context, logger *slog.Logger, conn *pgxpool.Conn, tntId int, eppPt string, hrmNm string) ([]UnauthSessionEndpoint, error){
	const (
		dbSchema = "web_core_unauth_ssn_ep_inf"
		dbFunc   = "ep_inf"
	)

	results, err := db.DataSet[UnauthSessionEndpoint](ctx, logger, conn,
		func(ctx *context.Context, tx *pgx.Tx)(string, string, *pgx.Rows, error){
			qry := fmt.Sprintf("select %v.%v($1, $2, $3, $4)", dbSchema, dbFunc)

			call, err := (*tx).Query(*ctx, qry, dbFunc, tntId, eppPt, hrmNm)
			if err != nil {

				slog.LogAttrs(*ctx, slog.LevelError, "UnauthSessionEndpointInfo::Get dataset",
					slog.String("error" , err.Error()),
					slog.String("qry"   , qry),
					slog.Int   ("tntId" , tntId),
					slog.String("eppPt" , eppPt),
					slog.String("hrmNm" , hrmNm),
				)

				return qry, dbFunc, nil, fmt.Errorf("UnauthSessionEndpointInfo::Call database function: %w", err)
			}

			return qry, dbFunc, &call, nil
		})

	return results, err
}

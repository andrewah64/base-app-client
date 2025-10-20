package aur

import (
	"context"
	"fmt"
	"log/slog"
	"time"
)

import (
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5"
)

import (
	"github.com/go-webauthn/webauthn/protocol"
	"github.com/go-webauthn/webauthn/webauthn"
)

import (
	"github.com/andrewah64/base-app-client/internal/common/core/db"
	"github.com/andrewah64/base-app-client/internal/web/core/passkey"
)

type AukcAtnInf struct {
	PuvNm   string
	PkhNm []string
}

func GetAukcAtnInf (ctx *context.Context, logger *slog.Logger, conn *pgxpool.Conn, tntId int) ([]AukcAtnInf, error) {
	const (
		dbSchema = "web_core_unauth_ssn_aur_reg"
		dbFunc   = "aukc_atn_inf"
	)

	results, err := db.DataSet[AukcAtnInf](ctx, logger, conn, func(ctx *context.Context, tx *pgx.Tx)(string, string, *pgx.Rows, error){
		qry := fmt.Sprintf("select %v.%v($1, $2)", dbSchema, dbFunc)

		call, err := (*tx).Query(*ctx, qry, dbFunc, tntId)
		if err != nil {
			slog.LogAttrs(*ctx, slog.LevelError, "GetAukcAtnInf::get dataset",
				slog.String("error"   , err.Error()),
				slog.String("qry"     , qry),
				slog.Int   ("tntId"   , tntId),
			)

			return qry, dbFunc, nil, fmt.Errorf("GetAukcAtnInf::call database function: %w", err)
		}

		return qry, dbFunc, &call, nil
	})

	return results, err
}

type AumInf struct {
	AupcEnabled      bool
	AukcEnabled      bool
	Saml2S2i         bool
	Saml2S2s         bool
	GoogleEnabled    bool
	MicrosoftEnabled bool
}

func GetAumInf (ctx *context.Context, logger *slog.Logger, conn *pgxpool.Conn, tntId int) ([]AumInf, error) {
	const (
		dbSchema = "web_core_unauth_ssn_aur_reg"
		dbFunc   = "aum_inf"
	)

	results, err := db.DataSet[AumInf](ctx, logger, conn, func(ctx *context.Context, tx *pgx.Tx)(string, string, *pgx.Rows, error){
		qry := fmt.Sprintf("select %v.%v($1, $2)", dbSchema, dbFunc)

		call, err := (*tx).Query(*ctx, qry, dbFunc, tntId)
		if err != nil {
			slog.LogAttrs(*ctx, slog.LevelError, "GetAum::get dataset",
				slog.String("error"   , err.Error()),
				slog.String("qry"     , qry),
				slog.Int   ("tntId"   , tntId),
			)

			return qry, dbFunc, nil, fmt.Errorf("GetAum::call database function: %w", err)
		}

		return qry, dbFunc, &call, nil
	})

	return results, err
}

type AurPkyInf struct {
	AurId int
	SsnDn time.Duration
	EppPt string
}

func GetAurPkyInf (ctx *context.Context, logger *slog.Logger, conn *pgxpool.Conn, tntId int, aurNm string) ([]AurPkyInf, error) {
	const (
		dbSchema = "web_core_unauth_ssn_aur_reg"
		dbFunc   = "aur_pky_inf"
	)

	results, err := db.DataSet[AurPkyInf](ctx, logger, conn, func(ctx *context.Context, tx *pgx.Tx)(string, string, *pgx.Rows, error){
		qry := fmt.Sprintf("select %v.%v($1, $2, $3)", dbSchema, dbFunc)

		call, err := (*tx).Query(*ctx, qry, dbFunc, tntId, aurNm)
		if err != nil {
			slog.LogAttrs(*ctx, slog.LevelError, "GetPkyAur::get dataset",
				slog.String("error"   , err.Error()),
				slog.String("qry"     , qry),
				slog.Int   ("tntId"   , tntId),
				slog.String("aurNm"   , aurNm),
			)

			return qry, dbFunc, nil, fmt.Errorf("GetPkyAur::call database function: %w", err)
		}

		return qry, dbFunc, &call, nil
	})

	return results, err
}

type AurPwdInf struct {
	AurId          int
	AurHshPw       string
	SsnDn          time.Duration
	EppPt          string
	AupcMfaEnabled bool
}

func GetAurPwdInf (ctx *context.Context, logger *slog.Logger, conn *pgxpool.Conn, tntId int, aurNm string) ([]AurPwdInf, error) {
	const (
		dbSchema = "web_core_unauth_ssn_aur_reg"
		dbFunc   = "aur_pwd_inf"
	)

	results, err := db.DataSet[AurPwdInf](ctx, logger, conn, func(ctx *context.Context, tx *pgx.Tx)(string, string, *pgx.Rows, error){
		qry := fmt.Sprintf("select %v.%v($1, $2, $3)", dbSchema, dbFunc)

		call, err := (*tx).Query(*ctx, qry, dbFunc, tntId, aurNm)
		if err != nil {
			slog.LogAttrs(*ctx, slog.LevelError, "GetPwdAur::get dataset",
				slog.String("error"   , err.Error()),
				slog.String("qry"     , qry),
				slog.Int   ("tntId"   , tntId),
				slog.String("aurNm"   , aurNm),
			)

			return qry, dbFunc, nil, fmt.Errorf("GetPwdAur::call database function: %w", err)
		}

		return qry, dbFunc, &call, nil
	})

	return results, err
}

type AurNmInf struct {
	AurNmPass bool
}

func GetAurNmInf (ctx *context.Context, logger *slog.Logger, conn *pgxpool.Conn, tntId int, aurNm string) ([]AurNmInf, error) {
	rs, rErr := db.DataSet[AurNmInf](ctx, logger, conn,
		func(ctx *context.Context, tx *pgx.Tx)(string, string, *pgx.Rows, error){
			dbFunc := "nm_inf"
			qry    := fmt.Sprintf("select web_core_unauth_ssn_aur_reg.%v($1, $2, $3)", dbFunc)

			c, cErr := (*tx).Query(*ctx, qry, dbFunc, tntId, aurNm)
			if cErr != nil {
				slog.LogAttrs(*ctx, slog.LevelError, "get dataset",
					slog.String("error" , cErr.Error()),
					slog.String("qry"   , qry),
					slog.Int   ("tntId" , tntId),
					slog.String("aurNm" , aurNm),
				)

				return qry, dbFunc, nil, fmt.Errorf("call database function: %w", cErr)
			}

			return qry, dbFunc, &c, nil
		})

	return rs, rErr
}

type PkyInf struct {
	PkyCredentialId           []byte
	PkyPublicKey              []byte
	PkyAttestationType          string
	PkyAuthenticatorTransport []string
	PkyUserPresent              bool
	PkyUserVerified             bool
	PkyBackupEligible           bool
	PkyBackupState              bool
	PkyAaguid                 []byte
	PkySignCount                int
	PkyCloneWarning             bool
	PkyAttachment               string
	PkyClientDataJson         []byte
	PkyClientDataHash         []byte
	PkyAuthenticatorData      []byte
	PkyPublicKeyAlgorithm       int64
	PkyObject                 []byte
}

func GetPkyInf (ctx *context.Context, logger *slog.Logger, conn *pgxpool.Conn, tntId int, aurNm string) ([]PkyInf, error) {
	rs, rErr := db.DataSet[PkyInf](ctx, logger, conn,
		func(ctx *context.Context, tx *pgx.Tx)(string, string, *pgx.Rows, error){
			dbFunc := "pky_inf"
			qry    := fmt.Sprintf("select web_core_unauth_ssn_aur_reg.%v($1, $2, $3)", dbFunc)

			c, cErr := (*tx).Query(*ctx, qry, dbFunc, tntId, aurNm)
			if cErr != nil {
				slog.LogAttrs(*ctx, slog.LevelError, "get dataset",
					slog.String("error" , cErr.Error()),
					slog.String("qry"   , qry),
					slog.Int   ("tntId" , tntId),
					slog.String("aurNm" , aurNm),
				)

				return qry, dbFunc, nil, fmt.Errorf("call database function: %w", cErr)
			}

			return qry, dbFunc, &c, nil
		})

	return rs, rErr
}

type PlsInf struct {
	PlsJs []byte
}

func GetPlsInf (ctx *context.Context, logger *slog.Logger, conn *pgxpool.Conn, tntId int, aurNm string, pkyChallenge string) ([]PlsInf, error) {
	const (
		dbSchema = "web_core_unauth_ssn_aur_reg"
		dbFunc   = "pls_inf"
	)

	results, err := db.DataSet[PlsInf](ctx, logger, conn, func(ctx *context.Context, tx *pgx.Tx)(string, string, *pgx.Rows, error){
		qry := fmt.Sprintf("select %v.%v($1, $2, $3, $4)", dbSchema, dbFunc)

		call, err := (*tx).Query(*ctx, qry, dbFunc, tntId, aurNm, pkyChallenge)
		if err != nil {
			slog.LogAttrs(*ctx, slog.LevelError, "GetPlsInf::get dataset",
				slog.String("error"        , err.Error()),
				slog.String("qry"          , qry),
				slog.Int   ("tntId"        , tntId),
				slog.String("aurNm"        , aurNm),
				slog.String("pkyChallenge" , aurNm),
			)

			return qry, dbFunc, nil, fmt.Errorf("GetPlsInf::call database function: %w", err)
		}

		return qry, dbFunc, &call, nil
	})

	return results, err
}

type PwdInf struct {
	AupcAurPwdMinLen int
	AupcAurPwdMaxLen int
}

func GetPwdInf (ctx *context.Context, logger *slog.Logger, conn *pgxpool.Conn, tntId int) ([]PwdInf, error) {
	const (
		dbSchema = "web_core_unauth_ssn_aur_reg"
		dbFunc   = "pwd_inf"
	)

	results, err := db.DataSet[PwdInf](ctx, logger, conn, func(ctx *context.Context, tx *pgx.Tx)(string, string, *pgx.Rows, error){
		qry := fmt.Sprintf("select %v.%v($1, $2)", dbSchema, dbFunc)

		call, err := (*tx).Query(*ctx, qry, dbFunc, tntId)
		if err != nil {
			slog.LogAttrs(*ctx, slog.LevelError, "GetAur::get dataset",
				slog.String("error"   , err.Error()),
				slog.String("qry"     , qry),
				slog.Int   ("tntId"   , tntId),
			)

			return qry, dbFunc, nil, fmt.Errorf("GetAur::call database function: %w", err)
		}

		return qry, dbFunc, &call, nil
	})

	return results, err
}

func GetPkyAur (ctx *context.Context, conn *pgxpool.Conn, tntId int, aurNm string) (*passkey.User, error) {
	pkyRs, pkyRsErr := GetPkyInf(ctx, slog.Default(), conn, tntId, aurNm)
	if pkyRsErr != nil {
		return nil, fmt.Errorf("aaa")
	}

	var credentials []webauthn.Credential = make([]webauthn.Credential, len(pkyRs))
	for i, v := range pkyRs {
		var t []protocol.AuthenticatorTransport = make([]protocol.AuthenticatorTransport, len(v.PkyAuthenticatorTransport))
		for i, v := range v.PkyAuthenticatorTransport {
			t[i] = protocol.AuthenticatorTransport(v)
		}

		flags := webauthn.CredentialFlags {
			UserPresent    : v.PkyUserPresent,
			UserVerified   : v.PkyUserVerified,
			BackupEligible : v.PkyBackupEligible,
			BackupState    : v.PkyBackupState,
		}

		authn := webauthn.Authenticator {
			AAGUID       : v.PkyAaguid,
			SignCount    : uint32(v.PkySignCount),
			CloneWarning : v.PkyCloneWarning,
			Attachment   : protocol.AuthenticatorAttachment(v.PkyAttachment),
		}

		attsn := webauthn.CredentialAttestation {
			ClientDataJSON     : v.PkyClientDataJson,
			ClientDataHash     : v.PkyClientDataHash,
			AuthenticatorData  : v.PkyAuthenticatorData,
			PublicKeyAlgorithm : v.PkyPublicKeyAlgorithm,
			Object             : v.PkyObject,
		}

		credentials[i] = webauthn.Credential {
			ID              : v.PkyCredentialId,
			PublicKey       : v.PkyPublicKey,
			AttestationType : v.PkyAttestationType,
			Transport       : t,
			Flags           : flags,
			Authenticator   : authn,
			Attestation     : attsn,
		}
	}

	user := &passkey.User{
		Id          : protocol.URLEncodedBase64(aurNm),
		Name        : aurNm,
		DisplayName : aurNm,
		Credentials : credentials,
	}

	return user, nil
}

func PostNnc (ctx *context.Context, logger *slog.Logger, conn *pgxpool.Conn, tntId int, aurId int, nncNonce string, nncExpTs time.Time, exptErrs []string) error {
	var (
		sprocCall   = "call web_core_unauth_ssn_aur_reg.reg_nnc(@p_tnt_id, @p_aur_id, @p_nnc_nonce, @p_nnc_exp_ts)"
		sprocParams = pgx.NamedArgs{
			"p_tnt_id"     : tntId,
			"p_aur_id"     : aurId,
			"p_nnc_nonce"  : nncNonce,
			"p_nnc_exp_ts" : nncExpTs,
		}
	)

	sprocErr := db.Sproc(ctx, logger, conn, sprocCall, sprocParams, exptErrs)
	if sprocErr != nil {
		logger.LogAttrs(*ctx, slog.LevelDebug, "call sproc",
			slog.String("sprocCall"  , sprocCall),
			slog.String("error"      , sprocErr.Error()),
			slog.Int   ("tntId"      , tntId),
			slog.Int   ("aurId"      , aurId),
			slog.String("nncNonce"   , nncNonce),
			slog.Any   ("nncExpTs"   , nncExpTs),
			slog.Any   ("exptErrs"   , exptErrs),
		)

		return sprocErr
	}

	return nil
}

func PostPls (ctx *context.Context, logger *slog.Logger, conn *pgxpool.Conn, tntId int, aurNm string, plsChallenge string, plsJs []byte, exptErrs []string) error {
	var (
		sprocCall   = "call web_core_unauth_ssn_aur_reg.reg_pls(@p_tnt_id, @p_aur_nm, @p_pls_challenge, @p_pls_js)"
		sprocParams = pgx.NamedArgs{
			"p_tnt_id"        : tntId,
			"p_aur_nm"        : aurNm,
			"p_pls_challenge" : plsChallenge,
			"p_pls_js"        : plsJs,
		}
	)

	sprocErr := db.Sproc(ctx, logger, conn, sprocCall, sprocParams, exptErrs)
	if sprocErr != nil {
		logger.LogAttrs(*ctx, slog.LevelDebug, "call sproc",
			slog.String("sprocCall"     , sprocCall),
			slog.String("error"         , sprocErr.Error()),
			slog.Int   ("tntId"         , tntId),
			slog.String("aurNm"         , aurNm),
			slog.String("plsChallenge"  , plsChallenge),
			slog.Any   ("plsJs"         , plsJs),
			slog.Any   ("exptErrs"      , exptErrs),
		)

		return sprocErr
	}

	return nil
}

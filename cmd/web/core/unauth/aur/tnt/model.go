package aur

import (
	"context"
	"fmt"
	"log/slog"
)

import (
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5"
)

import (
	"github.com/andrewah64/base-app-client/internal/common/core/db"
)

type AukcInf struct {
	AukcAurNmMinLen  int
	AukcAurNmMaxLen  int
}

func GetAukcInf (ctx *context.Context, logger *slog.Logger, conn *pgxpool.Conn, tntId int) ([]AukcInf, error) {
	const (
		dbSchema = "web_core_unauth_aur_tnt_reg"
		dbFunc   = "aukc_inf"
	)

	results, err := db.DataSet[AukcInf](ctx, logger, conn, func(ctx *context.Context, tx *pgx.Tx)(string, string, *pgx.Rows, error){
		qry := fmt.Sprintf("select %v.%v($1, $2)", dbSchema, dbFunc)

		call, err := (*tx).Query(*ctx, qry, dbFunc, tntId)
		if err != nil {
			slog.LogAttrs(*ctx, slog.LevelError, "GetAukcInf::get dataset",
				slog.String("error" , err.Error()),
				slog.String("qry"   , qry),
				slog.Int   ("tntId" , tntId),
			)

			return qry, dbFunc, nil, fmt.Errorf("GetAukcInf::call database function: %w", err)
		}

		return qry, dbFunc, &call, nil
	})

	return results, err
}

type AukcRegInf struct {
	PkaNm   string
	PktNm   string
	PdcNm   string
	PuvNm   string
	PkgCd []int
	PkhNm []string
}

func GetAukcRegInf (ctx *context.Context, logger *slog.Logger, conn *pgxpool.Conn, tntId int) ([]AukcRegInf, error) {
	const (
		dbSchema = "web_core_unauth_aur_tnt_reg"
		dbFunc   = "aukc_reg_inf"
	)

	results, err := db.DataSet[AukcRegInf](ctx, logger, conn, func(ctx *context.Context, tx *pgx.Tx)(string, string, *pgx.Rows, error){
		qry := fmt.Sprintf("select %v.%v($1, $2)", dbSchema, dbFunc)

		call, err := (*tx).Query(*ctx, qry, dbFunc, tntId)
		if err != nil {
			slog.LogAttrs(*ctx, slog.LevelError, "GetAukcRegInf::get dataset",
				slog.String("error" , err.Error()),
				slog.String("qry"   , qry),
				slog.Int   ("tntId" , tntId),
			)

			return qry, dbFunc, nil, fmt.Errorf("GetAukcRegInf::call database function: %w", err)
		}

		return qry, dbFunc, &call, nil
	})

	return results, err
}

type AumInf struct {
	AupcEnabled      bool
	AukcEnabled      bool
	GoogleEnabled    bool
	MicrosoftEnabled bool
}

func GetAumInf (ctx *context.Context, logger *slog.Logger, conn *pgxpool.Conn, tntId int) ([]AumInf, error) {
	const (
		dbSchema = "web_core_unauth_aur_tnt_reg"
		dbFunc   = "aum_inf"
	)

	results, err := db.DataSet[AumInf](ctx, logger, conn, func(ctx *context.Context, tx *pgx.Tx)(string, string, *pgx.Rows, error){
		qry := fmt.Sprintf("select %v.%v($1, $2)", dbSchema, dbFunc)

		call, err := (*tx).Query(*ctx, qry, dbFunc, tntId)
		if err != nil {
			slog.LogAttrs(*ctx, slog.LevelError, "GetAumInf::get dataset",
				slog.String("error" , err.Error()),
				slog.String("qry"   , qry),
				slog.Int   ("tntId" , tntId),
			)

			return qry, dbFunc, nil, fmt.Errorf("GetAumInf::call database function: %w", err)
		}

		return qry, dbFunc, &call, nil
	})

	return results, err
}

type AupcInf struct {
	AupcAurNmMinLen  int
	AupcAurNmMaxLen  int
	AupcAurPwdMinLen int
	AupcAurPwdMaxLen int
	AupcAurPwdIncSym bool
	AupcAurPwdIncNum bool
}

func GetAupcInf (ctx *context.Context, logger *slog.Logger, conn *pgxpool.Conn, tntId int) ([]AupcInf, error) {
	const (
		dbSchema = "web_core_unauth_aur_tnt_reg"
		dbFunc   = "aupc_inf"
	)

	results, err := db.DataSet[AupcInf](ctx, logger, conn, func(ctx *context.Context, tx *pgx.Tx)(string, string, *pgx.Rows, error){
		qry := fmt.Sprintf("select %v.%v($1, $2)", dbSchema, dbFunc)

		call, err := (*tx).Query(*ctx, qry, dbFunc, tntId)
		if err != nil {
			slog.LogAttrs(*ctx, slog.LevelError, "GetAupcInf::get dataset",
				slog.String("error" , err.Error()),
				slog.String("qry"   , qry),
				slog.Int   ("tntId" , tntId),
			)

			return qry, dbFunc, nil, fmt.Errorf("GetAupcInf::call database function: %w", err)
		}

		return qry, dbFunc, &call, nil
	})

	return results, err
}

type MfaInf struct {
	AupcMfaEnabled bool
}

func GetMfaInf (ctx *context.Context, logger *slog.Logger, conn *pgxpool.Conn, tntId int) ([]MfaInf, error) {
	const (
		dbSchema = "web_core_unauth_aur_tnt_reg"
		dbFunc   = "mfa_inf"
	)

	results, err := db.DataSet[MfaInf](ctx, logger, conn, func(ctx *context.Context, tx *pgx.Tx)(string, string, *pgx.Rows, error){
		qry := fmt.Sprintf("select %v.%v($1, $2)", dbSchema, dbFunc)

		call, err := (*tx).Query(*ctx, qry, dbFunc, tntId)
		if err != nil {
			slog.LogAttrs(*ctx, slog.LevelError, "GetMfaInf::get dataset",
				slog.String("error" , err.Error()),
				slog.String("qry"   , qry),
				slog.Int   ("tntId" , tntId),
			)

			return qry, dbFunc, nil, fmt.Errorf("GetMfaInf::call database function: %w", err)
		}

		return qry, dbFunc, &call, nil
	})

	return results, err
}

type PrsInf struct {
	PrsJs []byte
}

func GetPrsInf (ctx *context.Context, logger *slog.Logger, conn *pgxpool.Conn, tntId int, aurNm string) ([]PrsInf, error) {
	const (
		dbSchema = "web_core_unauth_aur_tnt_reg"
		dbFunc   = "prs_inf"
	)

	results, err := db.DataSet[PrsInf](ctx, logger, conn, func(ctx *context.Context, tx *pgx.Tx)(string, string, *pgx.Rows, error){
		qry := fmt.Sprintf("select %v.%v($1, $2, $3)", dbSchema, dbFunc)

		call, err := (*tx).Query(*ctx, qry, dbFunc, tntId, aurNm)
		if err != nil {
			slog.LogAttrs(*ctx, slog.LevelError, "GetPrsInf::get dataset",
				slog.String("error" , err.Error()),
				slog.String("qry"   , qry),
				slog.Int   ("tntId" , tntId),
				slog.String("aurNm" , aurNm),
			)

			return qry, dbFunc, nil, fmt.Errorf("GetPrsInf::call database function: %w", err)
		}

		return qry, dbFunc, &call, nil
	})

	return results, err
}

func PostPwAur (ctx *context.Context, logger *slog.Logger, conn *pgxpool.Conn, tntId int, aurNm string, aurHshPw string, aurEa string, otpId *string, otpSecret *string, exptErrs []string) error {
	var (
		sprocCall   = "call web_core_unauth_aur_tnt_reg.reg_aur(@p_tnt_id, @p_aur_nm, @p_aur_hsh_pw, @p_aur_ea, @p_otp_id, @p_otp_secret)"
		sprocParams = pgx.NamedArgs{
			"p_tnt_id"     : tntId,
			"p_aur_nm"     : aurNm,
			"p_aur_hsh_pw" : aurHshPw,
			"p_aur_ea"     : aurEa,
			"p_otp_id"     : otpId,
			"p_otp_secret" : otpSecret,
		}
	)

	sprocErr := db.Sproc(ctx, logger, conn, sprocCall, sprocParams, exptErrs)
	if sprocErr != nil {
		logger.LogAttrs(*ctx, slog.LevelDebug, "call sproc",
			slog.String("sprocCall" , sprocCall),
			slog.String("error"     , sprocErr.Error()),
			slog.Int   ("tntId"     , tntId),
			slog.String("aurNm"     , aurNm),
			slog.String("aurEa"     , aurEa),
			slog.Any   ("exptErrs"  , exptErrs),
		)

		return sprocErr
	}

	return nil
}

func PostPkyAur (ctx *context.Context, logger *slog.Logger, conn *pgxpool.Conn, tntId int, aurNm string, pkyEnabled bool, pkyCredentialId []byte, pkyPublicKey []byte, pkyAttestationType string, pkyAuthenticatorTransport []string, pkyUserPresent bool, pkyUserVerified bool, pkyBackupEligible bool, pkyBackupState bool, pkyAaguid []byte, pkySignCount int, pkyCloneWarning bool, pkyAttachment string, pkyClientDataJson []byte, pkyClientDataHash []byte, pkyAuthenticatorData []byte, pkyPublicKeyAlgorithm int64, pkyObject []byte, exptErrs []string) error {
	var (
		sprocCall   = "call web_core_unauth_aur_tnt_reg.reg_aur(@p_tnt_id, @p_aur_nm, @p_pky_enabled, @p_pky_credential_id, @p_pky_public_key, @p_pky_attestation_type, @p_pky_authenticator_transport, @p_pky_user_present, @p_pky_user_verified, @p_pky_backup_eligible, @p_pky_backup_state, @p_pky_aaguid, @p_pky_sign_count, @p_pky_clone_warning, @p_pky_attachment, @p_pky_client_data_json, @p_pky_client_data_hash, @p_pky_authenticator_data, @p_pky_public_key_algorithm, @p_pky_object)"
		sprocParams = pgx.NamedArgs{
			"p_tnt_id"                      : tntId,
			"p_aur_nm"                      : aurNm,
			"p_pky_enabled"                 : pkyEnabled,
			"p_pky_credential_id"           : pkyCredentialId,
			"p_pky_public_key"              : pkyPublicKey,
			"p_pky_attestation_type"        : pkyAttestationType,
			"p_pky_authenticator_transport" : pkyAuthenticatorTransport,
			"p_pky_user_present"            : pkyUserPresent,
			"p_pky_user_verified"           : pkyUserVerified,
			"p_pky_backup_eligible"         : pkyBackupEligible,
			"p_pky_backup_state"            : pkyBackupState,
			"p_pky_aaguid"                  : pkyAaguid,
			"p_pky_sign_count"              : pkySignCount,
			"p_pky_clone_warning"           : pkyCloneWarning,
			"p_pky_attachment"              : pkyAttachment,
			"p_pky_client_data_json"        : pkyClientDataJson,
			"p_pky_client_data_hash"        : pkyClientDataHash,
			"p_pky_authenticator_data"      : pkyAuthenticatorData,
			"p_pky_public_key_algorithm"    : pkyPublicKeyAlgorithm,
			"p_pky_object"                  : pkyObject,
		}
	)

	sprocErr := db.Sproc(ctx, logger, conn, sprocCall, sprocParams, exptErrs)
	if sprocErr != nil {
		logger.LogAttrs(*ctx, slog.LevelDebug, "call sproc",
			slog.String("sprocCall"          , sprocCall),
			slog.String("error"              , sprocErr.Error()),
			slog.Int   ("tntId"              , tntId),
			slog.String("aurNm"              , aurNm),
			slog.Bool  ("pkyEnabled"         , pkyEnabled),
			slog.Any   ("exptErrs"           , exptErrs),
			slog.Any   ("pkyCredentialId"    , pkyCredentialId),
			slog.Any   ("pkyPublicKey"       , pkyPublicKey),
			slog.String("pkyAttestationType" , pkyAttestationType),
		)

		return sprocErr
	}

	return nil
}

func PostPrs (ctx *context.Context, logger *slog.Logger, conn *pgxpool.Conn, tntId int, aurNm string, prsJs []byte, exptErrs []string) error {
	var (
		sprocCall   = "call web_core_unauth_aur_tnt_reg.reg_prs(@p_tnt_id, @p_aur_nm, @p_prs_js)"
		sprocParams = pgx.NamedArgs{
			"p_tnt_id" : tntId,
			"p_aur_nm" : aurNm,
			"p_prs_js" : prsJs,
		}
	)

	sprocErr := db.Sproc(ctx, logger, conn, sprocCall, sprocParams, exptErrs)
	if sprocErr != nil {
		logger.LogAttrs(*ctx, slog.LevelDebug, "call sproc",
			slog.String("sprocCall" , sprocCall),
			slog.String("error"     , sprocErr.Error()),
			slog.Int   ("tntId"     , tntId),
			slog.String("aurNm"     , aurNm),
			slog.Any   ("prsJs"     , prsJs),
			slog.Any   ("exptErrs"  , exptErrs),
		)

		return sprocErr
	}

	return nil
}

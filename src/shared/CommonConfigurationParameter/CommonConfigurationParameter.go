package ccp

import (
	"encoding/json"
	"fmt"
)

const (
	ValueType_String             = "string"
	ValueType_StringArray        = "strings"
	ValueType_IntegerNumber      = "integer"
	ValueType_IntegerNumberArray = "integers"
	ValueType_Boolean            = "boolean"
	ValueType_MapStringString    = "map"
)

const (
	Administrator                = "administrator"
	AllowNativePasswords         = "allowNativePasswords"
	BodyTemplateForLogIn         = "bodyTemplateForLogIn"
	BodyTemplateForRegVCode      = "bodyTemplateForRegVCode"
	CacheRecordTtl               = "cacheRecordTtl"
	CacheSizeLimit               = "cacheSizeLimit"
	CertFile                     = "certFile"
	CheckConnLiveness            = "checkConnLiveness"
	DatabaseName                 = "dbName"
	DatabaseType                 = "databaseType"
	DriverName                   = "driverName"
	EmailChangeRequestTtl        = "emailChangeRequestTtl"
	EnableSelfSignedCertificate  = "enableSelfSignedCertificate"
	FileCacheItemTtl             = "fileCacheItemTtl"
	FileCacheSizeLimit           = "fileCacheSizeLimit"
	FileCacheVolumeLimit         = "fileCacheVolumeLimit"
	FilesCountToClean            = "filesCountToClean"
	Host                         = "host"
	ImageHeight                  = "imageHeight"
	ImageWidth                   = "imageWidth"
	ImagesFolder                 = "imagesFolder"
	IsAdminApprovalRequired      = "isAdminApprovalRequired"
	IsCacheEnabled               = "isCacheEnabled"
	IsDatabaseInitialisationUsed = "isDatabaseInitialisationUsed"
	IsDebugMode                  = "isDebugMode"
	IsImageCleanupAtStartUsed    = "isImageCleanupAtStartUsed"
	IsImageServerEnabled         = "isImageServerEnabled"
	IsImageStorageUsed           = "isImageStorageUsed"
	IsStorageCleaningEnabled     = "isStorageCleaningEnabled"
	KeyFile                      = "keyFile"
	LogInRequestTtl              = "logInRequestTtl"
	LogOutRequestTtl             = "logOutRequestTtl"
	MaxAllowedPacket             = "maxAllowedPacket"
	Moderator                    = "moderator"
	Name                         = "name"
	Net                          = "net"
	Params                       = "params"
	Password                     = "password"
	PasswordChangeRequestTtl     = "passwordChangeRequestTtl"
	Path                         = "path"
	Port                         = "port"
	PrivateKeyFilePath           = "privateKeyFilePath"
	PublicKeyFilePath            = "publicKeyFilePath"
	RecordCacheItemTtl           = "recordCacheItemTtl"
	RecordCacheSizeLimit         = "recordCacheSizeLimit"
	RegistrationRequestTtl       = "registrationRequestTtl"
	RequestIdLength              = "requestIdLength"
	Schema                       = "schema"
	SessionMaxDuration           = "sessionMaxDuration"
	SigningMethod                = "signingMethod"
	SiteName                     = "siteName"
	SubjectTemplateForRegVCode   = "subjectTemplateForRegVCode"
	User                         = "user"
	UserAgent                    = "userAgent"
	UserNameMaxLenInBytes        = "userNameMaxLenInBytes"
	UserPasswordMaxLenInBytes    = "userPasswordMaxLenInBytes"
	VerificationCodeLength       = "verificationCodeLength"
)

const (
	ErrF_UnknownConfigurationParameterType = "unknown configuration parameter type: %s"
)

type CommonConfigurationParameter struct {
	Name  string
	Type  string
	Value any
}

func ParseCommonConfigurationParameterValue(rt string, rv json.RawMessage) (v any, err error) {
	switch rt {
	case ValueType_String:
		return parseVariableData[string](rv)

	case ValueType_StringArray:
		return parseVariableData[[]string](rv)

	case ValueType_IntegerNumber:
		return parseVariableData[int](rv)

	case ValueType_IntegerNumberArray:
		return parseVariableData[[]int](rv)

	case ValueType_Boolean:
		return parseVariableData[bool](rv)

	case ValueType_MapStringString:
		return parseVariableData[map[string]string](rv)

	default:
		return nil, fmt.Errorf(ErrF_UnknownConfigurationParameterType, rt)
	}
}

func parseVariableData[T any](src json.RawMessage) (dst T, err error) {
	err = json.Unmarshal(src, &dst)
	if err != nil {
		return dst, err
	}
	return dst, nil
}

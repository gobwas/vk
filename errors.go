package vk

type ErrorCode int

const (
	ErrEmpty                        ErrorCode = 0
	ErrUnknown                      ErrorCode = 1
	ErrAppStopped                   ErrorCode = 2
	ErrBadMethod                    ErrorCode = 3
	ErrBadSignature                 ErrorCode = 4
	ErrNoAuth                       ErrorCode = 5
	ErrRateLimitExceeded            ErrorCode = 6
	ErrPermissionDenied             ErrorCode = 7
	ErrBadRequest                   ErrorCode = 8
	ErrTooManyActions               ErrorCode = 9
	ErrInternalError                ErrorCode = 10
	ErrTestModeError                ErrorCode = 11
	ErrCaptchaRequired              ErrorCode = 14
	ErrAccessDenied                 ErrorCode = 15
	ErrSecureLayerRequired          ErrorCode = 16
	ErrUserValidationRequired       ErrorCode = 17
	ErrPageRemovedOrBlocked         ErrorCode = 18
	ErrNotStandaloneProhibited      ErrorCode = 20
	ErrOnlyStandaloneAllowed        ErrorCode = 21
	ErrMethodDepricated             ErrorCode = 23
	ErrUserPermissionRequired       ErrorCode = 24
	ErrInvalidCommunityAccessCode   ErrorCode = 27
	ErrInvalidApplicationAccessCode ErrorCode = 28
	ErrInsufficientParameters       ErrorCode = 100
	ErrBadAPIID                     ErrorCode = 101
	ErrBadUserID                    ErrorCode = 113
	ErrBadTimestamp                 ErrorCode = 150
	ErrAccessDeniedForAlbum         ErrorCode = 200
	ErrAccessDeniedForAudio         ErrorCode = 201
	ErrAccessDeniedForGroup         ErrorCode = 203
	ErrAlbumOverflow                ErrorCode = 300
	ErrActionDenied                 ErrorCode = 500
	ErrCommercialPermissionDenied   ErrorCode = 600
	ErrCommercialError              ErrorCode = 603
)

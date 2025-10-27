package errs

import "errors"

var (
	ErrNotfound                    = errors.New("not found")
	ErrInvalidRequestBody          = errors.New("invalid request body")
	ErrInvalidFieldValue           = errors.New("invalid field value")
	ErrSomethingWentWrong          = errors.New("something went wrong")
	ErrInvalidID                   = errors.New("invalid id")
	ErrUserNotFound                = errors.New("user not found")
	ErrUsernameAlreadyExists       = errors.New("username already exists")
	ErrEmailAlreadyExists          = errors.New("email already exists")
	ErrPhoneAlreadyExists          = errors.New("phone already exists")
	ErrIncorrectUsernameOrPassword = errors.New("incorrect username or password")
	ErrProductNotfound             = errors.New("product not found")
	ErrInvalidProductID            = errors.New("invalid product id")
	ErrInvalidProductName          = errors.New("invalid product name, min 4 symbols")
	ErrShopNotFound                = errors.New("shop not found")
	ErrInvalidShopID               = errors.New("invalid shop id")
	ErrInvalidShopName             = errors.New("invalid shop name, min 3 symbols")
	ErrShopAlreadyExists           = errors.New("shop already exists for this owner")
	ErrInvalidToken                = errors.New("invalid token")
	ErrInsufficientStock           = errors.New("insufficient stock")
	ErrOrderNotFound               = errors.New("order not found")
)

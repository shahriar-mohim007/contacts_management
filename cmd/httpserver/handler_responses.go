package httpserver

import (
	"contacts/utils"
	"net/http"
)

var ValidDataNotFound = utilis.ResponseState{
	StatusCode: http.StatusBadRequest,
	Message:    "The provided information is invalid. Please recheck and try again.",
}

var UserAlreadyExist = utilis.ResponseState{
	StatusCode: http.StatusBadRequest,
	Message:    "User Already Exist With this Email",
}

var UserActivated = utilis.ResponseState{
	StatusCode: http.StatusOK,
	Message:    "User activated successfully",
}

var UserCreated = utilis.ResponseState{
	StatusCode: http.StatusCreated,
	Message:    "User created successfully",
}

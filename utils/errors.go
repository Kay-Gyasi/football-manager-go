package utils

import "errors"

var (
	CantUpdateNewDocument = errors.New("document does not exist")
	InvalidObjectID       = errors.New("invalid id")
	CantFindCoach         = errors.New("coach not found")
	CantFindPlayer        = errors.New("player not found")
)

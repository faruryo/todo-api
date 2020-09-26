package repositories

import "errors"

var ErrNoSuchEntity = errors.New("no such entity")
var ErrBadRequestIdMustBeZero = errors.New("bad request: ID must be 0")
var ErrBadRequestIdMustNotBeZero = errors.New("bad request: ID must not be 0")
var ErrBadRequestUpdateCreatedAt = errors.New("bad request: CreatedAt can't update")
var ErrBadRequestUpdateUpdatedAt = errors.New("bad request: UpdatedAt can't udpate")

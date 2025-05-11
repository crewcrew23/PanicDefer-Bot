package dberrs

import (
	"errors"

	"github.com/lib/pq"
)

var (
	ErrUniqueConstraint   = errors.New("unique constraint violation")
	ErrInvalidData        = errors.New("invalid service data")
	ErrIsNullField        = errors.New("some field is null")
	ErrGetRows            = errors.New("failed to get rows")
	ErrNoRows             = errors.New("no rows")
	ErrNotEnoughtArgument = errors.New("not enought arg")
	ErrDbOperation        = errors.New("data base operation err")
)

func IsNotEnoughtArgumentError(err error) bool {
	var pgErr *pq.Error
	return errors.As(err, &pgErr) && pgErr.Code == "08P01"
}

func IsUniqueConstraintError(err error) bool {
	var pgErr *pq.Error
	return errors.As(err, &pgErr) && pgErr.Code == "23505"
}

func IsCheckConstraintError(err error) bool {
	var pgErr *pq.Error
	return errors.As(err, &pgErr) && pgErr.Code == "23514"
}

func IsNullFieldError(err error) bool {
	var pgErr *pq.Error
	return errors.As(err, &pgErr) && pgErr.Code == "23502"
}

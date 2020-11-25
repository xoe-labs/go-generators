// Code generated by 'ddd-gen app command': DO NOT EDIT.

package command

import (
	"context"
	errwrap "github.com/hashicorp/errwrap"
	errors "github.com/xoe-labs/ddd-gen/internal/test-svc/app/errors"
	offers "github.com/xoe-labs/ddd-gen/internal/test-svc/app/ifaces/offers"
	requires "github.com/xoe-labs/ddd-gen/internal/test-svc/app/ifaces/requires"
	"reflect"
)

// Topic: Account

var (
	// ErrNotAuthorizedToMakeNewAccountWithOutId signals that the caller is not authorized to perform MakeNewAccountWithOutId
	ErrNotAuthorizedToMakeNewAccountWithOutId = errors.NewAuthorizationError("ErrNotAuthorizedToMakeNewAccountWithOutId")
	// ErrMakeNewAccountWithOutIdHasNoTarget signals that MakeNewAccountWithOutId's target was not distinguishable
	ErrMakeNewAccountWithOutIdHasNoTarget = errors.NewTargetIdentificationError("ErrMakeNewAccountWithOutIdHasNoTarget")
	// ErrMakeNewAccountWithOutIdLoadingFailed signals that MakeNewAccountWithOutId storage failed to load the entity
	ErrMakeNewAccountWithOutIdLoadingFailed = errors.NewStorageLoadingError("ErrMakeNewAccountWithOutIdLoadingFailed")
	// ErrMakeNewAccountWithOutIdSavingFailed signals that MakeNewAccountWithOutId failed to save the entity
	ErrMakeNewAccountWithOutIdSavingFailed = errors.NewStorageSavingError("ErrMakeNewAccountWithOutIdSavingFailed")
	// ErrMakeNewAccountWithOutIdFailedInDomain signals that MakeNewAccountWithOutId failed in the domain layer
	ErrMakeNewAccountWithOutIdFailedInDomain = errors.NewDomainError("ErrMakeNewAccountWithOutIdFailedInDomain")
)

// MakeNewAccountWithOutIdHandlerWrapper knows how to perform MakeNewAccountWithOutId
type MakeNewAccountWithOutIdHandlerWrapper struct {
	rw requires.StorageWriterReader
	p  requires.Policer
}

// NewMakeNewAccountWithOutIdHandlerWrapper returns MakeNewAccountWithOutIdHandlerWrapper
func NewMakeNewAccountWithOutIdHandlerWrapper(rw requires.StorageWriterReader, p requires.Policer) *MakeNewAccountWithOutIdHandlerWrapper {
	if reflect.ValueOf(rw).IsZero() {
		panic("no 'rw' provided!")
	}
	if reflect.ValueOf(p).IsZero() {
		panic("no 'p' provided!")
	}
	return &MakeNewAccountWithOutIdHandlerWrapper{rw: rw, p: p}
}

// Handle generically performs MakeNewAccountWithOutId
func (h MakeNewAccountWithOutIdHandlerWrapper) Handle(ctx context.Context, mnawoi requires.DomainCommandHandler, actor offers.Policeable, target offers.Distinguishable) error {
	// assert that target is distinguishable
	if !target.IsDistinguishable() {
		return ErrMakeNewAccountWithOutIdHasNoTarget
	}
	// load entity from store; handle + wrap error
	a, loadErr := h.rw.Load(ctx, target)
	if loadErr != nil {
		return errwrap.Wrap(ErrMakeNewAccountWithOutIdLoadingFailed, loadErr)
	}
	// assert authorization via policy interface
	if ok := h.p.Can(ctx, actor, "MakeNewAccountWithOutId", a); !ok {
		// return opaque error: handle potentially sensitive policy errors out-of-band!
		return ErrNotAuthorizedToMakeNewAccountWithOutId
	}
	// assert correct command handling by the domain
	if ok := mnawoi.Handle(ctx, a); !ok {
		var domErr error
		// mnawoi is an ErrorKeeper
		for i, e := range mnawoi.Errors() {
			if i == 0 {
				domErr = e
			} else {
				domErr = errwrap.Wrap(domErr, e)
			}
		}
		return ErrMakeNewAccountWithOutIdFailedInDomain
	}
	// save domain facts to storage
	saveErr := h.rw.SaveFacts(ctx, target, requires.FactKeeper(mnawoi))
	if saveErr != nil {
		return errwrap.Wrap(ErrMakeNewAccountWithOutIdSavingFailed, saveErr)
	}
	return nil
}

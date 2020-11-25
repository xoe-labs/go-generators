// Code generated by 'ddd-gen app command': DO NOT EDIT.

package command

import (
	"context"
	errwrap "github.com/hashicorp/errwrap"
	app "github.com/xoe-labs/ddd-gen/internal/test-svc/app"
	errors "github.com/xoe-labs/ddd-gen/internal/test-svc/app/errors"
	"reflect"
)

// Topic: Account

var (
	// ErrNotAuthorizedToDeleteAccount signals that the caller is not authorized to perform DeleteAccount
	ErrNotAuthorizedToDeleteAccount = errors.NewAuthorizationError("ErrNotAuthorizedToDeleteAccount")
	// ErrDeleteAccountHasNoTarget signals that DeleteAccount's target was not distinguishable
	ErrDeleteAccountHasNoTarget = errors.NewTargetIdentificationError("ErrDeleteAccountHasNoTarget")
	// ErrDeleteAccountLoadingFailed signals that DeleteAccount storage failed to load the entity
	ErrDeleteAccountLoadingFailed = errors.NewStorageLoadingError("ErrDeleteAccountLoadingFailed")
	// ErrDeleteAccountSavingFailed signals that DeleteAccount failed to save the entity
	ErrDeleteAccountSavingFailed = errors.NewStorageSavingError("ErrDeleteAccountSavingFailed")
	// ErrDeleteAccountFailedInDomain signals that DeleteAccount failed in the domain layer
	ErrDeleteAccountFailedInDomain = errors.NewDomainError("ErrDeleteAccountFailedInDomain")
)

// DeleteAccountHandlerWrapper knows how to perform DeleteAccount
type DeleteAccountHandlerWrapper struct {
	rw app.RequiresStorageWriterReader
	p  app.RequiresPolicer
}

// NewDeleteAccountHandlerWrapper returns DeleteAccountHandlerWrapper
func NewDeleteAccountHandlerWrapper(rw app.RequiresStorageWriterReader, p app.RequiresPolicer) *DeleteAccountHandlerWrapper {
	if reflect.ValueOf(rw).IsZero() {
		panic("no 'rw' provided!")
	}
	if reflect.ValueOf(p).IsZero() {
		panic("no 'p' provided!")
	}
	return &DeleteAccountHandlerWrapper{rw: rw, p: p}
}

// Handle generically performs DeleteAccount
func (h DeleteAccountHandlerWrapper) Handle(ctx context.Context, da app.RequiresDomainCommandHandler, actor app.OffersPoliceable, target app.OffersDistinguishable) error {
	// assert that target is distinguishable
	if !target.IsDistinguishable() {
		return ErrDeleteAccountHasNoTarget
	}
	// load entity from store; handle + wrap error
	a, loadErr := h.rw.Load(ctx, target)
	if loadErr != nil {
		return errwrap.Wrap(ErrDeleteAccountLoadingFailed, loadErr)
	}
	// assert authorization via policy interface
	if ok := h.p.Can(ctx, actor, "DeleteAccount", a); !ok {
		// return opaque error: handle potentially sensitive policy errors out-of-band!
		return ErrNotAuthorizedToDeleteAccount
	}
	// assert correct command handling by the domain
	if ok := da.Handle(ctx, a); !ok {
		var domErr error
		// da is an ErrorKeeper
		for i, e := range da.Errors() {
			if i == 0 {
				domErr = e
			} else {
				domErr = errwrap.Wrap(domErr, e)
			}
		}
		return ErrDeleteAccountFailedInDomain
	}
	// save domain facts to storage
	saveErr := h.rw.SaveFacts(ctx, target, app.OffersFactKeeper(da))
	if saveErr != nil {
		return errwrap.Wrap(ErrDeleteAccountSavingFailed, saveErr)
	}
	return nil
}

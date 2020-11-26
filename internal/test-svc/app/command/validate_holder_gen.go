// Code generated by 'ddd-gen app command': DO NOT EDIT.

package command

import (
	"context"
	errwrap "github.com/hashicorp/errwrap"
	app "github.com/xoe-labs/ddd-gen/internal/test-svc/app"
	errors "github.com/xoe-labs/ddd-gen/internal/test-svc/app/errors"
	domain "github.com/xoe-labs/ddd-gen/internal/test-svc/domain"
	"reflect"
)

// Topic: Holder

var (
	// ErrValidateHolderHasNoTarget signals that ValidateHolder's target was not distinguishable
	ErrValidateHolderHasNoTarget = errors.NewTargetIdentificationError("ErrValidateHolderHasNoTarget")
	// ErrValidateHolderLoadingFailed signals that ValidateHolder storage failed to load the entity
	ErrValidateHolderLoadingFailed = errors.NewStorageLoadingError("ErrValidateHolderLoadingFailed")
	// ErrValidateHolderSavingFailed signals that ValidateHolder failed to save the entity
	ErrValidateHolderSavingFailed = errors.NewStorageSavingError("ErrValidateHolderSavingFailed")
	// ErrValidateHolderFailedInDomain signals that ValidateHolder failed in the domain layer
	ErrValidateHolderFailedInDomain = errors.NewDomainError("ErrValidateHolderFailedInDomain")
)

// ValidateHolderHandlerWrapper knows how to perform ValidateHolder
type ValidateHolderHandlerWrapper struct {
	rw app.RequiresStorageWriterReader
}

// NewValidateHolderHandlerWrapper returns ValidateHolderHandlerWrapper
func NewValidateHolderHandlerWrapper(rw app.RequiresStorageWriterReader) *ValidateHolderHandlerWrapper {
	if reflect.ValueOf(rw).IsZero() {
		panic("no 'rw' provided!")
	}
	return &ValidateHolderHandlerWrapper{rw: rw}
}

// Handle generically performs ValidateHolder
func (h ValidateHolderHandlerWrapper) Handle(ctx context.Context, vh domain.ValidateHolder, actor app.OffersPoliceable, target app.OffersDistinguishable) error {
	// assert that target is distinguishable
	if !target.IsDistinguishable() {
		return ErrValidateHolderHasNoTarget
	}
	// load entity from store; handle + wrap error
	a, loadErr := h.rw.Load(ctx, target)
	if loadErr != nil {
		return errwrap.Wrap(ErrValidateHolderLoadingFailed, loadErr)
	}
	// assert correct command handling by the domain
	if ok := vh.Handle(ctx, a); !ok {
		var domErr error
		for i, e := range vh.Errors() {
			if i == 0 {
				domErr = e
			} else {
				domErr = errwrap.Wrap(domErr, e)
			}
		}
		return ErrValidateHolderFailedInDomain
	}
	// save domain facts to storage
	saveErr := h.rw.SaveFacts(ctx, target, app.OffersFactKeeper(&vh))
	if saveErr != nil {
		return errwrap.Wrap(ErrValidateHolderSavingFailed, saveErr)
	}
	return nil
}

// compile time assertions
var (
	_ app.RequiresCommandHandler = (*domain.ValidateHolder)(nil)
	_ app.RequiresErrorKeeper    = (*domain.ValidateHolder)(nil)
	_ app.OffersFactKeeper       = (*domain.ValidateHolder)(nil)
)

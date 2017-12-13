package etgate

import (
    "fmt"

    sdk "github.com/cosmos/cosmos-sdk"

    eth "github.com/ethereum/go-ethereum/core/types"
)

var (
    errConflictingChain          = fmt.Errorf("Conflicting chain")
    errNotInitialized            = fmt.Errorf("Not initialized")
    errMissingSignature          = fmt.Errorf("Missing signature")
    errNotValidator              = fmt.Errorf("Sender is not a validator")
    errHeaderOutOfVisibleRange   = fmt.Errorf("Header out of visible range")
    errNoncontinuousFinalization = fmt.Errorf("Noncontinuous finalization")
    errNotEnoughSigns            = fmt.Errorf("Not enough signs")
    errLogHeaderNotFound         = fmt.Errorf("Log header not found")
    errInvalidLogProof           = fmt.Errorf("Invalid logproof")
    errLogUnpackingError         = fmt.Errorf("Log unpacking error")
    errAncestorNotFound          = fmt.Errorf("Ancestor not found")
    errInvalidCoins              = fmt.Errorf("Invalid coins")

    ETGateCodeConflictingChain        uint32 = 1001
    ETGateCodeNonContinuousHeaderList uint32 = 1002
)
/*
func ErrConflictingChain(hash common.Hash) error {
    return errors.WithMessage(hash.Hex(), errConflictingChain, ETGateCodeConflictingChain)
}*/

func ErrInvalidLogProof(log eth.Log, err error) error {
    return err
}

func ErrAlreadySignedHeader(signer sdk.Actor) error {
    return nil
}

func ErrHeaderNotFound(height uint64) error {
    return nil
}

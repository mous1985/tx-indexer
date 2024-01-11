package tx

import (
	"encoding/base64"
	"errors"

	"github.com/gnolang/gno/tm2/pkg/bft/types"
	"github.com/gnolang/tx-indexer/serve/encode"
	"github.com/gnolang/tx-indexer/serve/metadata"
	"github.com/gnolang/tx-indexer/serve/spec"
	storageErrors "github.com/gnolang/tx-indexer/storage/errors"
)

type Handler struct {
	storage Storage
}

func NewHandler(storage Storage) *Handler {
	return &Handler{
		storage: storage,
	}
}

func (h *Handler) GetTxHandler(
	_ *metadata.Metadata,
	params []any,
) (any, *spec.BaseJSONError) {
	// Check the params
	if len(params) < 1 {
		return nil, spec.GenerateInvalidParamCountError()
	}

	// Extract the params
	requestedTx, ok := params[0].(string)
	if !ok {
		return nil, spec.GenerateInvalidParamError(1)
	}

	// Decode the hash from base64
	decodedHash, err := base64.StdEncoding.DecodeString(requestedTx)
	if err != nil {
		return nil, spec.GenerateInvalidParamError(1)
	}

	// Run the handler
	response, err := h.getTx(decodedHash)
	if err != nil {
		return nil, spec.GenerateResponseError(err)
	}

	if response == nil {
		return nil, nil
	}

	encodedResponse, err := encode.EncodeValue(response)
	if err != nil {
		return nil, spec.GenerateResponseError(err)
	}

	return encodedResponse, nil
}

// getTx fetches the tx from storage, if any
func (h *Handler) getTx(txHash []byte) (*types.TxResult, error) {
	tx, err := h.storage.GetTx(txHash)
	if errors.Is(err, storageErrors.ErrNotFound) {
		// Wrap the error
		return nil, nil
	}

	if err != nil {
		return nil, err
	}

	return tx, nil
}
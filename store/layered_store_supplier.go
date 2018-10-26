package store


type LayeredStoreSupplierResult struct {
	StoreResult
}

func NewSupplierResult() *LayeredStoreSupplierResult {
	return &LayeredStoreSupplierResult{}
}

type LayeredStoreSupplier interface {
	SetChainNext(LayeredStoreSupplier)
	Next() LayeredStoreSupplier
}

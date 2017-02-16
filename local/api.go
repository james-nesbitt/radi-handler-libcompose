package local

import (
	"context"

	handlers_bytesource "github.com/wunderkraut/radi-handlers/bytesource"
)

// Settings needed to make a local API
type LocalAPISettings struct {
	handlers_bytesource.BytesourceFileSettings
	Context context.Context
}

package assert

import "github.com/IsmailCLN/tapir/internal/sharedcontext"

var sharedCtx *sharedcontext.SharedContext

func SetSharedContext(ctx *sharedcontext.SharedContext) {
    sharedCtx = ctx
}

func Ctx() *sharedcontext.SharedContext { return sharedCtx }

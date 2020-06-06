package cli

import (
	"context"
	root "github.com/finitum/node-cli/internal/commands"
	"github.com/finitum/node-cli/opts"
	"github.com/finitum/node-cli/provider"
)

func Run(originalCtx context.Context, ctx context.Context, s *provider.Store, c *opts.Opts) (int, int, error) {
	return root.RunRootCommand(originalCtx, ctx, s, c)
}

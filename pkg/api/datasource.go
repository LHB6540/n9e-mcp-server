package api

import (
	"context"

	"github.com/n9e/n9e-mcp-server/pkg/client"
	"github.com/n9e/n9e-mcp-server/pkg/toolset"
	"github.com/n9e/n9e-mcp-server/pkg/types"

	"github.com/google/jsonschema-go/jsonschema"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

// ListDatasourcesInput represents datasources list query parameters
type ListDatasourcesInput struct {
	Limit int `json:"limit,omitempty"`
	Page  int `json:"p,omitempty"`
}

// RegisterDatasourceToolset registers datasource toolset
func RegisterDatasourceToolset(group *toolset.ToolsetGroup, getClient client.GetClientFunc) {
	ts := toolset.NewToolset("datasource", "Datasource management tools")

	ts.AddReadTools(
		listDatasourcesTool(getClient),
	)

	group.AddToolset(ts)
}

func listDatasourcesTool(getClient client.GetClientFunc) toolset.ServerTool {
	return toolset.NewServerTool(
		mcp.Tool{
			Name:        "list_datasources",
			Description: "List all available datasources",
			Annotations: &mcp.ToolAnnotations{
				Title:        "List Datasources",
				ReadOnlyHint: true,
			},
			InputSchema: &jsonschema.Schema{
				Type: "object",
				Properties: map[string]*jsonschema.Schema{
					"limit": {
						Type:        "integer",
						Description: "Page size (default 20)",
					},
					"p": {
						Type:        "integer",
						Description: "Page number (starts from 1)",
					},
				},
			},
		},
		toolset.MakeToolHandler(func(ctx context.Context, req *mcp.CallToolRequest, input ListDatasourcesInput) (*mcp.CallToolResult, error) {
			c := getClient(ctx)
			if c == nil {
				return toolset.NewToolResultError("failed to get n9e client from context"), nil
			}

			result, err := client.DoGet[[]types.Datasource](c, ctx, "/api/n9e/datasource/brief", nil)
			if err != nil {
				return toolset.NewToolResultError(err.Error()), nil
			}

			items, total := toolset.SlicePage(result, input.Page, input.Limit)
			return toolset.MarshalResult(types.PageResp[types.Datasource]{List: items, Total: total}), nil
		}),
	)
}

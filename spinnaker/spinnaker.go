// Package spinnaker is the spinnaker client
package spinnaker

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"

	"config"
	server "github.com/himanhsugusain/go-mcp"
	"go.uber.org/zap"

	"golang.org/x/oauth2"
	"google.golang.org/api/idtoken"
	"google.golang.org/api/option"

	"go.lsp.dev/jsonrpc2"
)

type Spinnaker struct {
	client *ClientWithResponses
	log    *zap.Logger
}

func NewSpinnakerClient(ctx context.Context, log *zap.Logger) (*Spinnaker, error) {
	cfg, err := config.NewConfig(ctx)
	log.Info("server", zap.String("endpoint", cfg.Gate.Endpoint))
	if err != nil {
		return nil, err
	}
	tokenSource, err := idtoken.NewTokenSource(ctx, cfg.Auth.OAuthClientID, option.WithCredentialsFile(cfg.Auth.ServiceAccountKeyPath))
	if err != nil {
		return nil, err
	}
	httpClient := oauth2.NewClient(ctx, tokenSource)
	client, err := NewClientWithResponses("", WithHTTPClient(httpClient), WithBaseURL(cfg.Gate.Endpoint))
	return &Spinnaker{
		client: client,
		log:    log,
	}, err
}

func (s *Spinnaker) GetCapabilities() server.Capabilities {
	return server.Capabilities{
		Prompts:   server.Prompts{},
		Resources: server.Resources{},
		Tools: server.Tools{
			ListChanged: true,
		},
	}
}

func (s *Spinnaker) ServerInfo() server.ServerInfo {
	return server.ServerInfo{
		Name:    "Spinnaker-mcp-server",
		Title:   "spinnaker-mcp",
		Version: "0.0.1",
	}
}

func (s *Spinnaker) ListTools() server.ListToolResponse {
	return server.ListToolResponse{
		Tools: []server.Tool{
			{
				Name:        "getApplications",
				Title:       "Spinnaker Applications",
				Description: "Get list of Applications from spinnaker",
				InputSchema: map[string]any{
					"type":     "object",
					"required": []string{},
				},
			},
			{
				Name:        "getPipelines",
				Title:       "Spinnaker Pipelines under the Application",
				Description: "Get list of pipelines in spinnaker",
				InputSchema: map[string]any{
					"type": "object",
					"properties": map[string]any{
						"application": map[string]string{
							"type":        "string",
							"description": "Get list of pipelines from application",
						},
					},
					"required": []string{"application"},
				},
			},
			{
				Name:        "getPipeline",
				Title:       "Retrieve pipeline executions",
				Description: "Get execution for a pipeline based on pipelineconfigId and limit",
				InputSchema: map[string]any{
					"type": "object",
					"properties": map[string]any{
						"pipelineConfigId": map[string]string{
							"type":        "string",
							"description": "PipelineConfigId of the pipeline to get executions",
						},
						"limit": map[string]string{
							"type":        "string",
							"description": "Number of latest executions to fetch",
						},
					},
					"required": []string{"pipelineConfigId", "limit"},
				},
			},
		},
		NextCursor: "",
	}
}

func getPtr[U any](x U) *U {
	p := x
	return &p
}

func (s *Spinnaker) ToolsCall(call *jsonrpc2.Call) map[string]any {
	var params server.ToolParams
	json.Unmarshal(call.Params(), &params)
	switch params.Name {
	case "getApplications":
		{
			resp, err := s.client.GetAllApplicationsWithResponse(context.Background(), &GetAllApplicationsParams{})
			if err != nil {
				return server.ToolsErrorText(err)
			}
			return server.ToolsResponseText(string(resp.Body))
		}
	case "getPipelines":
		{
			resp, err := s.client.GetPipelinesWithResponse(context.Background(), params.Arguments["application"], &GetPipelinesParams{})
			if err != nil {
				return server.ToolsErrorText(err)
			}
			return server.ToolsResponseText(string(resp.Body))
		}
	case "getPipeline":
		{
			s.log.Debug("input", zap.Any("params", params))
			limit, err := strconv.ParseInt(params.Arguments["limit"], 10, 32)
			if err != nil {
				return server.ToolsErrorText(err)
			}
			resp, err := s.client.GetLatestExecutionsByConfigIdsWithResponse(context.Background(), &GetLatestExecutionsByConfigIdsParams{
				PipelineConfigIds: getPtr(params.Arguments["pipelineConfigId"]),
				Limit:             getPtr(int32(limit)),
			})
			if err != nil {
				return server.ToolsErrorText(err)
			}
			return server.ToolsResponseText(string(resp.Body))
		}
	default:
		return server.ToolsErrorText(fmt.Errorf("tools call not found"))
	}
}

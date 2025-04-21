package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/cloudwego/eino-ext/components/model/openai"
	"github.com/cloudwego/eino/callbacks"
	einomodel "github.com/cloudwego/eino/components/model"
	"github.com/cloudwego/eino/components/tool"
	"github.com/cloudwego/eino/compose"
	"github.com/cloudwego/eino/flow/agent"
	"github.com/cloudwego/eino/flow/agent/react"
	"github.com/cloudwego/eino/schema"
	callbacks2 "github.com/cloudwego/eino/utils/callbacks"
	"io"
)

func main() {

	apiKey := ""
	apiURL := ""
	modelName := "gpt-4o"
	ctx := context.Background()
	chatModel, err := openai.NewChatModel(ctx, &openai.ChatModelConfig{
		APIKey:  apiKey,
		BaseURL: apiURL,
		Model:   modelName,
	})
	if err != nil {
		panic(err)
	}
	dateTool := NewDateTool()
	tools := []tool.BaseTool{
		dateTool,
	}
	handler := callbacks2.NewHandlerHelper().
		ChatModel(NewChatModelCallbackHandler()).
		Handler()
	agentTest, err := react.NewAgent(ctx, &react.AgentConfig{
		ToolCallingModel: chatModel,
		MaxStep:          1000,
		ToolsConfig: compose.ToolsNodeConfig{
			Tools: tools,
		},
	})
	if err != nil {
		panic(err)
	}
	//config, future := react.WithMessageFuture()
	_, err = agentTest.Generate(ctx, []*schema.Message{
		schema.UserMessage("写一首古诗"),
	}, agent.WithComposeOptions(
		compose.WithCallbacks(handler),
	))

	if err != nil {
		fmt.Printf("stream error: %v\n", err)
		return
	}

}

func NewChatModelCallbackHandler() *callbacks2.ModelCallbackHandler {
	return &callbacks2.ModelCallbackHandler{
		OnStart: func(ctx context.Context, runInfo *callbacks.RunInfo, input *einomodel.CallbackInput) context.Context {
			fmt.Println("Model is starting...")
			fmt.Printf("info: %s\n", runInfo.Name)
			return ctx
		},
		OnEndWithStreamOutput: func(ctx context.Context, runInfo *callbacks.RunInfo, output *schema.StreamReader[*einomodel.CallbackOutput]) context.Context {
			fmt.Println("Model has ended with stream output.")

			go func() {
				defer output.Close()
				for {
					frame, err := output.Recv()

					if errors.Is(err, io.EOF) {
						break
					}
					if err != nil {
						if errors.Is(err, context.Canceled) {
							fmt.Printf("接收流取消: %s\n", err)
							return
						}
						fmt.Printf("接收流出错: %s\n", err)
						return
					}
					if frame.Message == nil {
						fmt.Println("nil 跳过")
						continue
					}
					jsonData, _ := json.MarshalIndent(frame.Message, "", "  ")
					fmt.Printf("frame: %s\n", string(jsonData))

				}
			}()
			return ctx
		},
		OnError: func(ctx context.Context, runInfo *callbacks.RunInfo, err error) context.Context {
			fmt.Printf("chat An error occurred: %s\n", err)
			fmt.Printf("info: %v\n", runInfo)

			return ctx
		},
		OnEnd: func(ctx context.Context, runInfo *callbacks.RunInfo, output *einomodel.CallbackOutput) context.Context {
			jsonframe, err := json.MarshalIndent(output, "", " ")
			if err != nil {
				panic(err)
			}
			fmt.Println(string(jsonframe))
			return ctx
		},
	}
}

package main

import (
	"context"
	"github.com/cloudwego/eino/components/tool"
	"github.com/cloudwego/eino/schema"
	"time"
)

type DateTool struct {
}

func NewDateTool() *DateTool {
	return &DateTool{}
}

func (dt *DateTool) Info(ctx context.Context) (*schema.ToolInfo, error) {
	return &schema.ToolInfo{
		Name: "dateTool",
		Desc: "获取当前北京时间,返回时间格式为yyyy-MM-dd HH:mm:ss",
	}, nil
}

func (dt *DateTool) InvokableRun(ctx context.Context, argumentsInJSON string, opts ...tool.Option) (string, error) {
	loc, _ := time.LoadLocation("Asia/Shanghai")
	return time.Now().In(loc).Format("2006-01-02 15:04:05"), nil
}

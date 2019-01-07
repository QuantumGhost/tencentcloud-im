package ctxutil

import (
	"context"
	"github.com/leanovate/gopter"
	"github.com/leanovate/gopter/gen"
	"github.com/leanovate/gopter/prop"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestGetAndSetCommandName(t *testing.T) {
	type testcase struct {
		ctx           context.Context
		expectedValue string
		expectedExist bool
		msgTmpl       string
	}
	cases := []testcase{
		{
			ctx:           context.Background(),
			expectedValue: "",
			expectedExist: false,
			msgTmpl:       "test failed for empty context",
		},
		{
			ctx:           context.WithValue(context.Background(), ctxKeyCommand, 1),
			expectedValue: "",
			expectedExist: false,
			msgTmpl:       "tested failed for invalid type",
		},
	}
	for _, c := range cases {
		cmd, ok := GetCommand(c.ctx)
		assert.Equal(t, c.expectedValue, cmd, c.msgTmpl)
		assert.Equal(t, c.expectedExist, ok, c.msgTmpl)
	}
	property := gopter.NewProperties(nil)
	property.Property("set then set should return the original value", prop.ForAll(func(s string) bool {
		cmd, ok := GetCommand(SetCommand(context.Background(), s))
		return ok && cmd == s
	}, gen.AnyString()))
	property.TestingRun(t)
}

func TestGetAndSetServiceName(t *testing.T) {
	type testcase struct {
		ctx           context.Context
		expectedValue string
		expectedExist bool
		msgTmpl       string
	}
	cases := []testcase{
		{
			ctx:           context.Background(),
			expectedValue: "",
			expectedExist: false,
			msgTmpl:       "test failed for empty context",
		},
		{
			ctx:           context.WithValue(context.Background(), ctxKeyServiceName, 1),
			expectedValue: "",
			expectedExist: false,
			msgTmpl:       "tested failed for invalid type",
		},
	}
	for _, c := range cases {
		cmd, ok := GetServiceName(c.ctx)
		assert.Equal(t, c.expectedValue, cmd, c.msgTmpl)
		assert.Equal(t, c.expectedExist, ok, c.msgTmpl)
	}
	property := gopter.NewProperties(nil)
	property.Property("set then set should return the original value", prop.ForAll(func(s string) bool {
		svc, ok := GetServiceName(SetServiceName(context.Background(), s))
		return ok && svc == s
	}, gen.AnyString()))
	property.TestingRun(t)
}

func TestGetAndSetRequestStartTime(t *testing.T) {
	type testcase struct {
		ctx           context.Context
		expectedValue time.Time
		expectedExist bool
		msgTmpl       string
	}
	cases := []testcase{
		{
			ctx:           context.Background(),
			expectedValue: time.Time{},
			expectedExist: false,
			msgTmpl:       "test failed for empty context",
		},
		{
			ctx:           context.WithValue(context.Background(), ctxKeyRequestStartTime, 1),
			expectedValue: time.Time{},
			expectedExist: false,
			msgTmpl:       "tested failed for invalid type",
		},
	}
	for _, c := range cases {
		cmd, ok := GetRequestStartTime(c.ctx)
		assert.Equal(t, c.expectedValue, cmd, c.msgTmpl)
		assert.Equal(t, c.expectedExist, ok, c.msgTmpl)
	}
	property := gopter.NewProperties(nil)
	property.Property("set then set should return the original value", prop.ForAll(func(st time.Time) bool {
		startTime, ok := GetRequestStartTime(SetRequestStartTime(context.Background(), st))
		return ok && st.Equal(startTime)
	}, gen.AnyTime()))
	property.TestingRun(t)
}

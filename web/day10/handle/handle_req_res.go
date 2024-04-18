package handle

import (
	"context"
	"fmt"
	"reflect"
	"strings"
)

type ReqResFunc struct {
	fn reflect.Value // 函数调用入口

	ctx reflect.Type // 第一个请求参数：context.Context
	req reflect.Type // 第二个请求参数：XXXReq

	res reflect.Type // 第一个返回参数：XXXRes
	err reflect.Type // 第二个返回参数：error
}

// NewReqResFunc 返回 ReqResFunc
// 参数 reqResFunc 必须是 func(context.Context, *XXXReq) (*XXXRes, error) 格式，否则会触发 panic
func NewReqResFunc(reqRes any) *ReqResFunc {
	fn := reflect.ValueOf(reqRes)
	fnType := fn.Type()

	if fnType.NumIn() != 2 {
		panic("parameter must be context.Context and XXXReq")
	}
	if fnType.NumOut() != 2 {
		panic("return value must be XXXRes and error")
	}

	ctx, req := fnType.In(0), fnType.In(1)
	res, err := fnType.Out(0), fnType.Out(1)

	if !ctx.Implements(reflect.TypeOf((*context.Context)(nil)).Elem()) {
		panic("the first parameter must be context.Context")
	}
	if !err.Implements(reflect.TypeOf((*error)(nil)).Elem()) {
		panic("the second return value must be error")
	}

	// req.Kind() must be *Struct
	// todo add Struct
	if req.Kind() == reflect.Pointer && req.Elem().Kind() == reflect.Struct {
		if !strings.HasSuffix(req.String(), "Req") {
			panic(fmt.Sprintf(`invalid struct name for request: defined as "%s", but it should be named with "Res" suffix like "XxxReq" or "*XxxReq"`, req.String()))
		}
	} else {
		panic(fmt.Sprintf(`invalid handler: defined as "%s", but type of the  second input parameter should be like "BizReq" or "*BizReq"`, req.String()))
	}

	// res.Kind() must be Struct or *Struct
	if res.Kind() == reflect.Struct ||
		(res.Kind() == reflect.Pointer && res.Elem().Kind() == reflect.Struct) {
		if !strings.HasSuffix(res.String(), "Res") {
			panic(fmt.Sprintf(`invalid struct name for request: defined as "%s", but it should be named with "Res" suffix like "XxxRes" or "*XxxReq"`, res.String()))
		}
	} else {
		panic(fmt.Sprintf(`invalid handler: defined as "%s", but type of the first output parameter should be "BizRes" or "*BizRes"`, res.String()))
	}

	return &ReqResFunc{
		fn:  fn,
		ctx: ctx,
		req: req,
		res: res,
		err: err,
	}
}

func (f *ReqResFunc) Call(ctx context.Context, decode func(point any) error) (any, error) {
	req := reflect.New(f.req.Elem())
	point := req.Interface()

	if err := decode(point); err != nil {
		return nil, err
	}

	result := f.fn.Call([]reflect.Value{reflect.ValueOf(ctx), req})
	if err := result[1]; !err.IsNil() {
		return nil, err.Interface().(error)
	}
	return result[0].Interface(), nil
}

func (f *ReqResFunc) DecodeFunc() DecodeFunc {
	return f.Call
}

func (f *ReqResFunc) Req() reflect.Type {
	if f.req.Kind() == reflect.Pointer {
		return f.req.Elem()
	}
	return f.req
}
func (f *ReqResFunc) Res() reflect.Type { return f.res }

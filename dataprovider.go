// MIT License
//
// Copyright (c) 2025 Eigen
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in all
// copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
// SOFTWARE.

package dataprovider

import (
	"context"
	"reflect"
)

type inner struct {
	route             map[string]DataProvider
	writeOpIntercept  map[string]func(context.Context, any)
	deleteOpIntercept map[string]func(context.Context, any)
	fallback          DataProvider
}

func (i *inner) AddRoute(data interface{}, p DataProvider) {
	i.route[reflect.TypeOf(data).String()] = p
}

func (i *inner) AddWriteOpInterceptor(data interface{}, f func(context.Context, any)) {
	i.writeOpIntercept[reflect.TypeOf(data).String()] = f
}

func (i *inner) AddDeleteOpInterceptor(data interface{}, f func(context.Context, any)) {
	i.deleteOpIntercept[reflect.TypeOf(data).String()] = f
}

func (i *inner) Create(ctx context.Context, data any) error {
	defer i.writeOpInterceptor(ctx, data)
	if p, ok := i.route[reflect.TypeOf(data).String()]; ok {
		return p.Create(ctx, data)
	} else {
		return i.fallback.Create(ctx, data)
	}
}

func (i *inner) Set(ctx context.Context, data any) error {
	defer i.writeOpInterceptor(ctx, data)
	if p, ok := i.route[reflect.TypeOf(data).String()]; ok {
		return p.Set(ctx, data)
	} else {
		return i.fallback.Set(ctx, data)
	}
}

func (i *inner) Get(ctx context.Context, data any) error {
	if p, ok := i.route[reflect.TypeOf(data).String()]; ok {
		return p.Get(ctx, data)
	} else {
		return i.fallback.Get(ctx, data)
	}
}

func (i *inner) Delete(ctx context.Context, data any) error {
	defer i.deleteOpInterceptor(ctx, data)
	if p, ok := i.route[reflect.TypeOf(data).String()]; ok {
		return p.Delete(ctx, data)
	} else {
		return i.fallback.Delete(ctx, data)
	}
}

func (i *inner) writeOpInterceptor(ctx context.Context, data any) {
	if p, ok := i.writeOpIntercept[reflect.TypeOf(data).String()]; ok {
		p(ctx, data)
	}
}

func (i *inner) deleteOpInterceptor(ctx context.Context, data any) {
	if p, ok := i.deleteOpIntercept[reflect.TypeOf(data).String()]; ok {
		p(ctx, data)
	}
}

/*
 * Copyright (c) 2024 by Randy Bell.  All rights reserved.
 *
 * This Source Code Form is subject to the terms of the Apache Public License, version 2.0. If a copy of the APL was not distributed with this file, you can obtain one at https://www.apache.org/licenses/LICENSE-2.0.txt.
 */

package grpcServer

import (
	"context"
	"github.com/rbell/toolchest/server/example/proto"
)

type HelloService struct {
	proto.UnimplementedHelloServiceServer
}

func (h *HelloService) SayHello(ctx context.Context, request *proto.HelloRequest) (*proto.HelloReply, error) {
	return &proto.HelloReply{Message: "Hello, " + request.Name}, nil
}

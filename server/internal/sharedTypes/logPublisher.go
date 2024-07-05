/*
 * Copyright (c) 2024 by Randy Bell.  All rights reserved.
 *
 * This Source Code Form is subject to the terms of the Apache Public License, version 2.0. If a copy of the APL was not distributed with this file, you can obtain one at https://www.apache.org/licenses/LICENSE-2.0.txt.
 */

package sharedTypes

import "context"

type LogPublisher interface {
	InfoContext(ctx context.Context, msg string, args ...any)
	Info(msg string, args ...any)
	WarnContext(ctx context.Context, msg string, args ...any)
	Warn(msg string, args ...any)
	ErrorContext(ctx context.Context, msg string, args ...any)
	Error(msg string, args ...any)
}

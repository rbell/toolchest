/*
 * Copyright (c) 2024 by Randy Bell.  All rights reserved.
 *
 * This Source Code Form is subject to the terms of the Apache Public License, version 2.0. If a copy of the APL was not distributed with this file, you can obtain one at https://www.apache.org/licenses/LICENSE-2.0.txt.
 */

package contextUtil

import (
	"context"
	"log/slog"
	"time"
)

func WithSLog(ctx context.Context, log slog.Logger) context.Context {
	ctx := context.WithTimeout(ctx, 10*time.Second)
	return context.WithValue(ctx, "slog", log)
}
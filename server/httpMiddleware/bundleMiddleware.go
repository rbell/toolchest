/*
 * Copyright (c) 2024 by Randy Bell.  All rights reserved.
 *
 * This Source Code Form is subject to the terms of the Apache Public License, version 2.0. If a copy of the APL was not distributed with this file, you can obtain one at https://www.apache.org/licenses/LICENSE-2.0.txt.
 */

package httpMiddleware

import "net/http"

func BundleMiddleware(middleware ...HttpHandlerMiddleware) HttpHandlerMiddleware {
	return func(next http.HandlerFunc) http.HandlerFunc {
		if len(middleware) == 0 {
			return next
		}

		wrapped := next
		// loop in reverse to preserve middleware order
		for i := len(middleware) - 1; i >= 0; i-- {
			wrapped = middleware[i](wrapped)
		}
		return wrapped
	}
}

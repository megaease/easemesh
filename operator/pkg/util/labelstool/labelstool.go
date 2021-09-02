/*
 * Copyright (c) 2017, MegaEase
 * All rights reserved.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package labelstool

import (
	"fmt"
	"strings"
)

// Marshal transforms labels in type map to string.
func Marshal(labels map[string]string) string {
	labelsSlice := []string{}
	for k, v := range labels {
		labelsSlice = append(labelsSlice, k+"="+v)
	}
	// FIXME: Replace & with , in all in the future.
	return strings.Join(labelsSlice, "&")
}

// Unmarshal transforms labels in type string to map.
// The latest value will cover earlier one with the same key.
// So `k1=v1,k2=v2,k1=v3` will be {"k1": "v3", "k2": "v2"}
func Unmarshal(s string) (map[string]string, error) {
	result := make(map[string]string)
	if s == "" {
		return result, nil
	}

	kvs := strings.Split(s, ",")
	for _, kv := range kvs {
		label := strings.Split(kv, "=")
		if len(label) != 2 {
			return nil, fmt.Errorf("invalid labels")
		}
		result[label[0]] = label[1]
	}

	return result, nil
}

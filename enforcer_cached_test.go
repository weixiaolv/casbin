// Copyright 2018 The casbin Authors. All Rights Reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package casbin

import "testing"

func testEnforceCache(t *testing.T, e *CachedEnforcer, sub string, obj interface{}, act string, res bool) {
	t.Helper()
	if myRes, _ := e.Enforce(sub, obj, act); myRes != res {
		t.Errorf("%s, %v, %s: %t, supposed to be %t", sub, obj, act, myRes, res)
	}
}

func TestCache(t *testing.T) {
	e, _ := NewCachedEnforcer("examples/basic_model.conf", "examples/basic_policy.csv")
	// The cache is enabled by default for NewCachedEnforcer.

	testEnforceCache(t, e, "alice", "data1", "read", true)
	testEnforceCache(t, e, "alice", "data1", "write", false)
	testEnforceCache(t, e, "alice", "data2", "read", false)
	testEnforceCache(t, e, "alice", "data2", "write", false)

	// The cache is enabled, calling RemovePolicy, LoadPolicy or RemovePolicies will
	// also operate cached items.
	_, _ = e.RemovePolicy("alice", "data1", "read")

	testEnforceCache(t, e, "alice", "data1", "read", false)
	testEnforceCache(t, e, "alice", "data1", "write", false)
	testEnforceCache(t, e, "alice", "data2", "read", false)
	testEnforceCache(t, e, "alice", "data2", "write", false)

	e, _ = NewCachedEnforcer("examples/rbac_model.conf", "examples/rbac_policy.csv")

	testEnforceCache(t, e, "alice", "data1", "read", true)
	testEnforceCache(t, e, "bob", "data2", "write", true)
	testEnforceCache(t, e, "alice", "data2", "read", true)
	testEnforceCache(t, e, "alice", "data2", "write", true)

	_, _ = e.RemovePolicies([][]string{
		{"alice", "data1", "read"},
		{"bob", "data2", "write"},
	})

	testEnforceCache(t, e, "alice", "data1", "read", false)
	testEnforceCache(t, e, "bob", "data2", "write", false)
	testEnforceCache(t, e, "alice", "data2", "read", true)
	testEnforceCache(t, e, "alice", "data2", "write", true)
	
}

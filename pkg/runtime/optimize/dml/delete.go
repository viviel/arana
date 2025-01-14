/*
 * Licensed to the Apache Software Foundation (ASF) under one or more
 * contributor license agreements.  See the NOTICE file distributed with
 * this work for additional information regarding copyright ownership.
 * The ASF licenses this file to You under the Apache License, Version 2.0
 * (the "License"); you may not use this file except in compliance with
 * the License.  You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package dml

import (
	"context"
)

import (
	"github.com/pkg/errors"
)

import (
	"github.com/arana-db/arana/pkg/proto"
	"github.com/arana-db/arana/pkg/runtime/ast"
	"github.com/arana-db/arana/pkg/runtime/optimize"
	"github.com/arana-db/arana/pkg/runtime/plan"
	"github.com/arana-db/arana/pkg/runtime/plan/dml"
)

func init() {
	optimize.Register(ast.SQLTypeDelete, optimizeDelete)
}

func optimizeDelete(ctx context.Context, o *optimize.Optimizer) (proto.Plan, error) {
	var (
		stmt = o.Stmt.(*ast.DeleteStatement)
	)

	shards, err := o.ComputeShards(stmt.Table, stmt.Where, o.Args)
	if err != nil {
		return nil, errors.Wrap(err, "failed to optimize DELETE statement")
	}

	// TODO: delete from a child sharding-table directly

	if shards == nil {
		return plan.Transparent(stmt, o.Args), nil
	}

	ret := dml.NewSimpleDeletePlan(stmt)
	ret.BindArgs(o.Args)
	ret.SetShards(shards)

	return ret, nil
}

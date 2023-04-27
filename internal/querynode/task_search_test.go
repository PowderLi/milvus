// Licensed to the LF AI & Data foundation under one
// or more contributor license agreements. See the NOTICE file
// distributed with this work for additional information
// regarding copyright ownership. The ASF licenses this file
// to you under the Apache License, Version 2.0 (the
// "License"); you may not use this file except in compliance
// with the License. You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package querynode

import (
	"testing"

	"github.com/milvus-io/milvus/internal/proto/planpb"
	"github.com/stretchr/testify/suite"
)

type SearchTaskSuite struct {
	suite.Suite
}

func (s *SearchTaskSuite) TestMerge() {
	Params.Init()

	plan := &planpb.PlanNode{
		Node: &planpb.PlanNode_VectorAnns{},
	}

	s1 := &searchTask{
		NQ:   1,
		TopK: 1,
		plan: plan,
	}
	s2 := &searchTask{
		NQ:   1,
		TopK: 5,
		plan: plan,
	}
	s3 := &searchTask{
		NQ:   1,
		TopK: 25,
		plan: plan,
	}
	s4 := &searchTask{
		NQ:   5,
		TopK: 40,
		plan: plan,
	}

	s.Equal(s1.CanMergeWith(s2), s2.CanMergeWith(s1))

	s.True(s1.CanMergeWith(s2))
	s.True(s2.CanMergeWith(s3))
	s.True(s3.CanMergeWith(s4))

	// exceed the ratio (10)
	s.False(s1.CanMergeWith(s3))

	// merge s1, s2. now it's with nq=2, topk=10
	s1.Merge(s2)
	s.True(s1.CanMergeWith(s3))
}

func TestSearchTask(t *testing.T) {
	suite.Run(t, new(SearchTaskSuite))
}

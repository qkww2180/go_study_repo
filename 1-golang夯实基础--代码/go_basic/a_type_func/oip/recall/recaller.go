package recall

import "dqq/go/basic/a_type_func/oip/common"

type Recaller interface {
	Recall(n int) []*common.Product //生成一批推荐候选集
	Name() string
}

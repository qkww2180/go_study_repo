package db

import (
	"dqq/basic/init/logger"
	"fmt"
)

var DbClient int

func InitDb() {
	DbClient = logger.Log + 9
	fmt.Printf("DbClient=%d\n", DbClient)
}

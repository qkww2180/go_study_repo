package test

import (
	"fmt"
	"dqq/util"
	"testing"
	"time"

	"github.com/bytedance/sonic"
)

func TestParseTime(t *testing.T) {
	str := "2007-07-01T00:00:00+08:00"
	tm, _ := time.ParseInLocation(util.DATETIME_LAYOUT3, str, util.Loc)
	fmt.Println(tm.Format(util.DATETIME_LAYOUT1))
}

func TestParseDate(t *testing.T) {
	if tm, err := util.ParseDate("1980-2-15"); err != nil {
		fmt.Println(err)
		t.Fail()
	} else {
		fmt.Println(tm.Format(util.DATE_LAYOUT1))
	}

	if tm, err := util.ParseDate("1980-02-15"); err != nil {
		fmt.Println(err)
		t.Fail()
	} else {
		fmt.Println(tm.Format(util.DATE_LAYOUT1))
	}

	if tm, err := util.ParseDate("1980-02"); err != nil {
		fmt.Println(err)
		t.Fail()
	} else {
		fmt.Println(tm.Format(util.DATE_LAYOUT1))
	}

	if tm, err := util.ParseDate("1980.02.14"); err != nil {
		fmt.Println(err)
		t.Fail()
	} else {
		fmt.Println(tm.Format(util.DATE_LAYOUT1))
	}

	if tm, err := util.ParseDate("1980.02"); err != nil {
		fmt.Println(err)
		t.Fail()
	} else {
		fmt.Println(tm.Format(util.DATE_LAYOUT1))
	}

	if tm, err := util.ParseDate("1980年2月14日"); err != nil {
		fmt.Println(err)
		t.Fail()
	} else {
		fmt.Println(tm.Format(util.DATE_LAYOUT1))
	}

	if tm, err := util.ParseDate("1980年2月"); err != nil {
		fmt.Println(err)
		t.Fail()
	} else {
		fmt.Println(tm.Format(util.DATE_LAYOUT1))
	}

	if tm, err := util.ParseDate("1980年"); err != nil {
		fmt.Println(err)
		t.Fail()
	} else {
		fmt.Println(tm.Format(util.DATE_LAYOUT1))
	}

	if tm, err := util.ParseDate("1980/2/15"); err != nil {
		fmt.Println(err)
		t.Fail()
	} else {
		fmt.Println(tm.Format(util.DATE_LAYOUT1))
	}

	if tm, err := util.ParseDate("1980/02"); err != nil {
		fmt.Println(err)
		t.Fail()
	} else {
		fmt.Println(tm.Format(util.DATE_LAYOUT1))
	}

	if tm, err := util.ParseDate("1980"); err != nil {
		fmt.Println(err)
		t.Fail()
	} else {
		fmt.Println(tm.Format(util.DATE_LAYOUT1))
	}
}

func TestMyJsonDateTime(t *testing.T) {
	inst := struct { //匿名结构体，只使用一次
		Birthday util.MyJsonDateTime
	}{Birthday: util.MyJsonDateTime{Time: time.Now(), Format: util.DATETIME_LAYOUT1}}
	s, _ := sonic.MarshalString(inst)
	fmt.Println(s)
}

// go test -v ./util/test -run=^TestParseDate$ -count=1
// go test -v ./util/test -run=^TestMyJsonDateTime$ -count=1

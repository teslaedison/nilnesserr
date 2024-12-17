package a

import (
	"context"
	"fmt"
	"math/rand/v2"
	"strings"
)

func Do() error {
	if rand.Float64() > 0.5 {
		return fmt.Errorf("do err")
	}
	return nil
}

func Do2() error {
	if rand.Float64() > 0.5 {
		return fmt.Errorf("do err")
	}
	return nil
}

func Do3() (int, error) {
	if rand.Float64() > 0.5 {
		return 1, fmt.Errorf("do err")
	}
	return 0, nil
}

func Empty() int {
	var a int
	a += 1
	return a
}

func Call2() error {
	err := Do()
	if err != nil {
		return err
	}
	return err
}

func Call() error {
	err1 := Do()
	if err1 != nil {
		return err1
	}
	err2 := Do2()
	if err2 != nil {
		a := 1
		a = a + 2
		fmt.Println(a)
		if a > 10 {
			fmt.Println(a)
			if a > 11 {
				return err1 // want `return a nil value error`
			}
		}
	}
	return nil
}

func Call3() error {
	err := Do()
	if err == nil {
		return err
	}
	return err
}

func Call4() (error, error) {
	err := Do()
	if err != nil {
		return nil, err
	}
	err2 := Do2()
	if err2 != nil {
		return err, err2 // want `return a nil value error`
	}
	return nil, nil
}

func Call5() error {
	err := Do()
	if err != nil {
		return err
	}
	a, err := Do3()
	if err != nil {
		return err
	}
	_ = a
	if err := Do2(); err != nil {
		return err
	}
	return err
}

func Call6(ctx context.Context, in string) (int, error) {
	if err := Do(); err != nil {
		return 0, err
	}
	if !strings.Contains(in, "${{") || !strings.Contains(in, "}}") {
		return 1, nil
	}
	res, err := Do3()
	if err != nil {
		return 23, err
	}
	ret := res + 1
	if err := Do2(); err != nil {
		return 4, err
	}
	ret = ret + 1
	return 5, err
}

func Call7() error {
	var a any = int(1)
	switch a.(type) {
	case int:
		err := Do()
		if err != nil {
			return err
		}
		if err2 := Do2(); err2 != nil {
			return err2
		}
		return err
	case string:
	default:
		return nil
	}

	return nil
}

// bad case
func Call8() error {
	err := Do()
	if err != nil {
		return err
	}

	err2 := Do()
	if err2 == nil {
		return err
	} else {
		if _, err := Do3(); err != nil {
			return err
		}
		return err2
	}
}

// bad case
func Call9() (err error) {
	if err = Do(); err != nil {
		return
	} else if err = Do2(); err != nil {
		return
	}

	_, _ = Do3()
	return
}

// bad case
func Call10() (int, error) {
	res, err := Do3()
	if err == nil {
		num, err := Do3()
		num += res
		if err != nil {
			return 0, err
		}
		if num > 0 {
			err := Do()
			if err != nil {
				return num, err
			}
		}
		return num, err
	}
	return 0, err
}

func Call11() error {
	_, err := Do3()
	if err != nil {
		return err
	} else if err = Do2(); err != nil {
		return err
	}
	return err
}

func Call12() (err error) {
	err = Do()
	if err != nil {
		return err
	}
	err2 := Do2()
	if err2 != nil {
		return
	}
	return
}

func Call13() error {
	_, err := Do3()
	if err != nil {
		_, err := Do3()
		if err != nil {
			return err
		}
		return err
	}
	return nil
}

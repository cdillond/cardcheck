package cardcheck

import (
	"testing"
)

var Ex = []struct {
	ab_cnum uint64
	ans     uint64
	valid   bool
	err     error
}{
	{37828224631000, 5, true, nil},
	{42222222222, 2, true, nil},
	{10, 9, true, nil},
	{2, 6, true, nil},
	{0, 0, true, nil},
	{1234567890123456789, 4, true, nil},
}

var Exb = []struct {
	ab_cnum string
	fcnum   string
	ans     uint64
	valid   bool
	err     error
}{
	{"37828224631000", "378282246310005", 5, true, nil},
	{"-601111111111111", "-6011111111111110", 0, false, InvalidInputError{}},
	{"6*01111111111111", "6*011111111111110", 0, false, InvalidInputError{}},
	{"601111111111111&", "601111111111111&0", 0, false, InvalidInputError{}},
	{"42222222222", "422222222222", 2, true, nil},
	{"42222222222", "422222222223", 2, false, nil},
	{"10", "109", 9, true, nil},
	{"2", "26", 6, true, nil},
	{"0", "0", 0, true, nil},
	{"1234567890123456789", "12345678901234567894", 4, true, nil},
	{"1345678901234567894*1", "1345678901234567894*1", 0, false, nil},
	{"123456789012345678-945621654", "123456789012345678-9454656454", 0, false, InvalidInputError{}},
	{"18446744073709551615", "184467440737095516153", 3, true, nil},
	{"", "", 0, false, InvalidInputError{}},
	{"1844674407370955161518446744073709551615184467440737095516151844674", "18446744073709551615184467440737095516151844674407370955161518446748", 8, true, nil},
}

func BenchmarkGetCheckDigit(b *testing.B) {
	for i := 0; i < b.N; i++ {
		for _, val := range Ex {
			GetCheckDigit(val.ab_cnum)
		}
	}
}

func BenchmarkCheckLuhn(b *testing.B) {
	for i := 0; i < b.N; i++ {
		for _, val := range Ex {
			CheckLuhn(val.ab_cnum*10 + val.ans)
		}
	}
}

func TestGetCheckDigit(t *testing.T) {
	for _, val := range Ex {
		cd := GetCheckDigit(val.ab_cnum)

		if cd != val.ans {
			t.Errorf("Incorrect output. cnum: %d, wanted: %d, got: %d", val.ab_cnum, val.ans, cd)
		}
	}
}

func TestCheckLuhn(t *testing.T) {
	for _, val := range Ex {
		valid := CheckLuhn(val.ab_cnum*10 + val.ans)
		if valid != val.valid {
			t.Errorf("Incorrect output. cnum: %d, cd: %d", val.ab_cnum, val.ans)
		}
	}
}

func TestStrGetCD(t *testing.T) {
	for _, val := range Exb {
		cd, err := StrGetCD(val.ab_cnum)
		if err != nil {
			t.Log(err, val.ab_cnum)
		}
		if cd != val.ans {
			t.Errorf("%s wanted: %d, got: %d", val.ab_cnum, val.ans, cd)
		}
	}

}

func BenchmarkStrGetCD(b *testing.B) {
	for i := 0; i < b.N; i++ {
		for _, val := range Exb {
			StrGetCD(val.ab_cnum)
		}
	}
}

func BenchmarkStrCheckLuhn(b *testing.B) {
	for i := 0; i < b.N; i++ {
		for _, val := range Exb {
			StrCheckLuhn(val.fcnum)
		}
	}

}

func TestStrCheckLuhn(t *testing.T) {
	for _, val := range Exb {
		cd, err := StrCheckLuhn(val.fcnum)
		if err != nil {
			t.Log(err, val.fcnum)
		}
		if cd != val.valid {
			t.Errorf("%s, %v", val.fcnum, cd)
		}
	}

}

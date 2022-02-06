package cardcheck

import (
	"strconv"
)

type InvalidInputError struct {
}

func (e InvalidInputError) Error() string {
	return "invalid input"
}

// calculates luhn checkdigit for any valid uint64.
// leading zeros can be ignored without
// affecting accuracy. more performant
// than StrGetCD, but argument size is limited.
func GetCheckDigit(cnum uint64) uint64 {
	ev := true
	var csum uint64
	var dig uint64
	for ; cnum > 0; cnum /= 10 {
		dig = cnum % 10

		switch ev {
		case false:
			csum += dig
		case true:

			if dig > 4 {
				csum += 2*dig - 9
			} else {
				csum += 2 * dig
			}
		}
		ev = !ev
	}
	return (10 - csum%10) % 10
}

// checks if the final digit of a uint64
// matches the luhn checkdigit generated using
// the preceding digits in the uint64.
// more performant than StrCheckLuhn,
// but argument size is limited.
func CheckLuhn(fcnum uint64) bool {
	cd := GetCheckDigit(fcnum / 10)
	return cd == fcnum%10
}

// calculates the checkdigit for any
// arbitrarily long string of digits (signed prefixes
// are not permitted). returns 0 and an error if the
// argument cannot be parsed.
func StrGetCD(cnum string) (uint64, error) {
	val_uint, err := strconv.ParseUint(cnum, 10, 64)
	if err == nil {
		return GetCheckDigit(val_uint), nil
	}

	clen := len(cnum)
	var cnums []uint64
	start := clen % 18
	if start != 0 {
		chunk, err := strconv.ParseUint(cnum[:start], 10, 64)
		if err != nil {
			return 0, InvalidInputError{}
		}
		cnums = append(cnums, chunk)
	}
	for i := start; i < clen; i += 18 {
		chunk, err := strconv.ParseUint(cnum[i:i+18], 10, 64)
		if err != nil {
			return 0, InvalidInputError{}
		}
		cnums = append(cnums, chunk)
	}

	var csum uint64
	var dig uint64
	for i := len(cnums) - 1; i >= 0; i-- {
		ev := true

		num := cnums[i]
		for ; num > 0; num /= 10 {
			dig = num % 10

			switch ev {
			case false:
				csum += dig
			case true:

				if dig > 4 {
					csum += 2*dig - 9
				} else {
					csum += 2 * dig
				}
			}
			ev = !ev
		}
	}

	return (10 - csum%10) % 10, nil
}

// checks if the final digit of any
// arbitrarily long string of digits matches
// the checkdigit generated using the
// preceding digits in the string.
// returns false and an error if argument
// cannot be parsed. signed prefixes and
// other non-numeric characters are not permitted.
func StrCheckLuhn(fcnum string) (bool, uint64, error) {
	cd, err := strconv.ParseUint(fcnum[len(fcnum)-1:], 10, 64)
	if err != nil {
		return false, 0, InvalidInputError{}
	}
	res, err := StrGetCD(fcnum[:len(fcnum)-1])
	if err != nil {
		return false, 0, err
	}

	return cd == res, res, nil
}
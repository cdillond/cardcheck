package cardcheck

import (
	"strconv"
)

type InvalidInputError struct {
}

func (e InvalidInputError) Error() string {
	return "invalid input"
}

// GetCheckDigit calculates the Luhn checkdigit for any valid uint64.
// Leading zeros can be ignored without affecting accuracy.
// More performant than StrGetCD, but the argument size is limited.
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

// CheckLuhn checks if the final digit of a uint64 matches the Luhn checkdigit generated using the preceding digits in the uint64.
// More performant than StrCheckLuhn, but argument size is limited.
func CheckLuhn(fcnum uint64) bool {
	cd := GetCheckDigit(fcnum / 10)
	return cd == fcnum%10
}

// StrGetCD calculates the checkdigit for any arbitrarily long string of digits.
// Signed prefixes are not permitted.
// Returns 0 and an InvalidInputError if the argument cannot be parsed.
func StrGetCD(cnum string) (uint64, error) {
	val_uint, err := strconv.ParseUint(cnum, 10, 64)
	if err == nil {
		return GetCheckDigit(val_uint), nil
	}

	clen := len(cnum)
	if clen == 0 {
		return 0, InvalidInputError{}
	}
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

// StrCheckLuhn checks if the final digit of any arbitrarily long string of digits matches
// the checkdigit generated using the preceding digits in the string.
// Returns false and an InvalidInputError if argument cannot be parsed.
// Signed prefixes and other non-numeric characters are not permitted.
func StrCheckLuhn(fcnum string) (bool, error) {
	if len(fcnum) == 0 {
		return false, InvalidInputError{}
	}
	cd, err := strconv.ParseUint(fcnum[len(fcnum)-1:], 10, 64)
	if err != nil {
		return false, InvalidInputError{}
	}
	res, err := StrGetCD(fcnum[:len(fcnum)-1])
	if err != nil {
		return false, err
	}

	return cd == res, nil
}

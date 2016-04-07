// Package iban implements validation of IBAN as defined in ISO 13616
package iban

import (
	"errors"
	"math/big"
	"regexp"
	"strings"
)

// Validate performs a sanity check on the IBAN number provided
func Validate(iban string) error {
	iban = normalizeIBAN(iban)
	if !validIBANChars(iban) {
		return errors.New("invalid characters")
	}

	if f, ok := format[iban[:2]]; !ok {
		return errors.New("invalid country code")
	} else {
		// Validating the BBAN number
		if len(iban) != f.size {
			return errors.New("invalid BBAN length")
		}
	}

	rearranged := []byte(iban[4:] + iban[:4])
	t := make([]string, 0, len(rearranged))
	for _, c := range rearranged {
		t = append(t, alphaToNum[c])
	}

	x := big.NewInt(0)
	x, ok := x.SetString(strings.Join(t, ""), 10)
	if !ok {
		return errors.New("internal error")
	}

	modulo := big.NewInt(0)
	modulo.Mod(x, big.NewInt(97))
	if modulo.Int64() != 1 {
		return errors.New("invalid IBAN")
	}
	return nil
}

func normalizeIBAN(iban string) string {
	iban = strings.ToUpper(iban)
	return replaceIgnoredChars(iban, "")
}

var (
	replaceIgnoredChars = regexp.MustCompile(`[\s\-]+`).ReplaceAllString
	validIBANChars      = regexp.MustCompile(`^[0-9A-Z]{15,34}$`).MatchString
)

var alphaToNum = map[byte]string{
	'0': "0",
	'1': "1",
	'2': "2",
	'3': "3",
	'4': "4",
	'5': "5",
	'6': "6",
	'7': "7",
	'8': "8",
	'9': "9",
	'A': "10",
	'B': "11",
	'C': "12",
	'D': "13",
	'E': "14",
	'F': "15",
	'G': "16",
	'H': "17",
	'I': "18",
	'J': "19",
	'K': "20",
	'L': "21",
	'M': "22",
	'N': "23",
	'O': "24",
	'P': "25",
	'Q': "26",
	'R': "27",
	'S': "28",
	'T': "29",
	'U': "30",
	'V': "31",
	'W': "32",
	'X': "33",
	'Y': "34",
	'Z': "35",
}

type country struct {
	name     string
	size     int
	bban     string
	code     string
	format   string
	comment  string
	standard bool
}

var format = map[string]country{
	"AL": country{name: "Albania", size: 28, bban: "8n, 16c", code: "AL", format: "ALkk bbbs sssx cccc cccc cccc cccc", comment: "b = National bank code s = Branch code x = National check digit c = Account number", standard: true},
	"AD": country{name: "Andorra", size: 24, bban: "8n,12c", code: "AD", format: "ADkk bbbb ssss cccc cccc cccc", comment: "b = National bank code s = Branch code c = Account number", standard: true},
	"AT": country{name: "Austria", size: 20, bban: "16n", code: "AT", format: "ATkk bbbb bccc cccc cccc", comment: "b = National bank code c = Account number", standard: true},
	"AZ": country{name: "Azerbaijan", size: 28, bban: "4c,20n", code: "AZ", format: "AZkk bbbb cccc cccc cccc cccc cccc", comment: "b = National bank code c = Account number", standard: true},
	"BH": country{name: "Bahrain", size: 22, bban: "4a,14c", code: "BH", format: "BHkk bbbb cccc cccc cccc cc", comment: "b = National bank code c = Account number", standard: true},
	"BE": country{name: "Belgium", size: 16, bban: "12n", code: "BE", format: "BEkk bbbc cccc ccxx", comment: "b = National bank code c = Account number x = National check digits", standard: true},
	"BA": country{name: "Bosnia and Herzegovina", size: 20, bban: "16n", code: "BA", format: "BAkk bbbs sscc cccc ccxx", comment: "k = IBAN check digits (always 39) b = National bank code s = Branch code c = Account number x = National check digits", standard: true},
	"BR": country{name: "Brazil", size: 29, bban: "23n, 1a, 1c", code: "BR", format: "BRkk bbbb bbbb ssss sccc cccc ccct n", comment: "k = IBAN check digits (Calculated by MOD 97-10) b = National bank code s = Branch code c = Account Number t = Account type (Cheque account, Savings account etc.) n = Owner account number (1, 2 etc.)[31]", standard: true},
	"BG": country{name: "Bulgaria", size: 22, bban: "4a,6n,8c", code: "BG", format: "BGkk bbbb ssss ddcc cccc cc", comment: "b = BIC bank code s = Branch (BAE) number d = Account type c = Account number", standard: true},
	"CR": country{name: "Costa Rica", size: 21, bban: "17n", code: "CR", format: "CRkk bbbc cccc cccc cccc c", comment: "b = bank code c = Account number", standard: true},
	"HR": country{name: "Croatia", size: 21, bban: "17n", code: "HR", format: "HRkk bbbb bbbc cccc cccc c", comment: "b = Bank code c = Account number", standard: true},
	"CY": country{name: "Cyprus", size: 28, bban: "8n,16c", code: "CY", format: "CYkk bbbs ssss cccc cccc cccc cccc", comment: "b = National bank code s = Branch code c = Account number", standard: true},
	"CZ": country{name: "Czech Republic", size: 24, bban: "20n", code: "CZ", format: "CZkk bbbb ssss sscc cccc cccc", comment: "b = National bank code s = Account number prefix c = Account number", standard: true},
	"DK": country{name: "Denmark", size: 18, bban: "14n", code: "DK", format: "DKkk bbbb cccc cccc cc", comment: "b = National bank code c = Account number", standard: true},
	"DO": country{name: "Dominican Republic", size: 28, bban: "4a,20n", code: "DO", format: "DOkk bbbb cccc cccc cccc cccc cccc", comment: "b = Bank identifier c = Account number", standard: true},
	"TL": country{name: "East Timor", size: 23, bban: "19n", code: "TL", format: "TLkk bbbc cccc cccc cccc cxx", comment: "k = IBAN check digits (always = '38') b = Bank identifier c = Account number x = National check digit", standard: true},
	"EE": country{name: "Estonia", size: 20, bban: "16n", code: "EE", format: "EEkk bbss cccc cccc cccx", comment: "b = National bank code s = Branch code c = Account number x = National check digit", standard: true},
	"FO": country{name: "Faroe Islands", size: 18, bban: "14n", code: "FO", format: "FOkk bbbb cccc cccc cx", comment: "b = National bank code c = Account number x = National check digit", standard: true},
	"FI": country{name: "Finland", size: 18, bban: "14n", code: "FI", format: "FIkk bbbb bbcc cccc cx", comment: "b = Bank and branch code c = Account number x = National check digit", standard: true},
	"FR": country{name: "France", size: 27, bban: "10n,11c,2n", code: "FR", format: "FRkk bbbb bggg ggcc cccc cccc cxx", comment: "b = National bank code g = Branch code (fr:code guichet) c = Account number x = National check digits (fr:clé RIB)", standard: true},
	"GE": country{name: "Georgia", size: 22, bban: "2c,16n", code: "GE", format: "GEkk bbcc cccc cccc cccc cc", comment: "b = National bank code c = Account number", standard: true},
	"DE": country{name: "Germany", size: 22, bban: "18n", code: "DE", format: "DEkk bbbb bbbb cccc cccc cc", comment: "b = Bank and branch identifier (de:Bankleitzahl or BLZ) c = Account number", standard: true},
	"GI": country{name: "Gibraltar", size: 23, bban: "4a,15c", code: "GI", format: "GIkk bbbb cccc cccc cccc ccc", comment: "b = BIC bank code c = Account number", standard: true},
	"GR": country{name: "Greece", size: 27, bban: "7n,16c", code: "GR", format: "GRkk bbbs sssc cccc cccc cccc ccc", comment: "b = National bank code s = Branch code c = Account number", standard: true},
	"GL": country{name: "Greenland", size: 18, bban: "14n", code: "GL", format: "GLkk bbbb cccc cccc cc", comment: "b = National bank code c = Account number", standard: true},
	"GT": country{name: "Guatemala", size: 28, bban: "4c,20c", code: "GT", format: "GTkk bbbb mmtt cccc cccc cccc cccc", comment: "b = National bank code c = Account number m = Currency t = Account type ", standard: true},
	"HU": country{name: "Hungary", size: 28, bban: "24n", code: "HU", format: "HUkk bbbs sssk cccc cccc cccc cccx", comment: "b = National bank code s = Branch code c = Account number x = National check digit", standard: true},
	"IS": country{name: "Iceland", size: 26, bban: "22n", code: "IS", format: "ISkk bbbb sscc cccc iiii iiii ii", comment: "b = National bank code s = Branch code c = Account number i = holder's kennitala (national identification number).", standard: true},
	"IE": country{name: "Ireland", size: 22, bban: "4c,14n", code: "IE", format: "IEkk aaaa bbbb bbcc cccc cc", comment: "a = BIC bank code b = Bank/branch code (sort code) c = Account number", standard: true},
	"IK": country{name: "Israel", size: 23, bban: "19n", code: "IK", format: "ILkk bbbn nncc cccc cccc ccc", comment: "b = National bank code n = Branch number c = Account number 13 digits (padded with zeros)", standard: true},
	"IT": country{name: "Italy", size: 27, bban: "1a,10n,12c", code: "IT", format: "ITkk xaaa aabb bbbc cccc cccc ccc", comment: "x = Check char (CIN) a = National bank code (it:Associazione bancaria italiana or Codice ABI ) b = Branch code (it:Coordinate bancarie or CAB – Codice d'Avviamento Bancario) c = Account number", standard: true},
	"JO": country{name: "Jordan", size: 30, bban: "4a, 22n", code: "JO", format: "JOkk bbbb nnnn cccc cccc cccc cccc cc", comment: "b = National bank code n = Branch code c = Account number ", standard: true},
	"KZ": country{name: "Kazakhstan", size: 20, bban: "3n,13c", code: "KZ", format: "KZkk bbbc cccc cccc cccc", comment: "b = National bank code c = Account number ", standard: true},
	"XK": country{name: "Kosovo", size: 20, bban: "4n,10n,2n", code: "XK", format: "XKkk bbbb cccc cccc cccc", comment: "b = National bank code c = Account number", standard: true},
	"KW": country{name: "Kuwait", size: 30, bban: "4a, 22c", code: "KW", format: "KWkk bbbb cccc cccc cccc cccc cccc cc", comment: "b = National bank code c = Account number.", standard: true},
	"LV": country{name: "Latvia", size: 21, bban: "4a,13c", code: "LV", format: "LVkk bbbb cccc cccc cccc c", comment: "b = BIC Bank code c = Account number", standard: true},
	"LB": country{name: "Lebanon", size: 28, bban: "4n,20c", code: "LB", format: "LBkk bbbb cccc cccc cccc cccc cccc", comment: "b = National bank code c = Account number", standard: true},
	"LI": country{name: "Liechtenstein", size: 21, bban: "5n,12c", code: "LI", format: "LIkk bbbb bccc cccc cccc c", comment: "b = National bank code c = Account number", standard: true},
	"LT": country{name: "Lithuania", size: 20, bban: "16n", code: "LT", format: "LTkk bbbb bccc cccc cccc", comment: "b = National bank code c = Account number", standard: true},
	"LU": country{name: "Luxembourg", size: 20, bban: "3n,13c", code: "LU", format: "LUkk bbbc cccc cccc cccc", comment: "b = National bank code c = Account number", standard: true},
	"MK": country{name: "Macedonia", size: 19, bban: "3n,10c,2n", code: "MK", format: "MKkk bbbc cccc cccc cxx", comment: "k = IBAN check digits (always = '07') b = National bank code c = Account number x = National check digits", standard: true},
	"MT": country{name: "Malta", size: 31, bban: "4a,5n,18c", code: "MT", format: "MTkk bbbb ssss sccc cccc cccc cccc ccc", comment: "b = BIC bank code s = Branch code c = Account number", standard: true},
	"MR": country{name: "Mauritania", size: 27, bban: "23n", code: "MR", format: "MRkk bbbb bsss sscc cccc cccc cxx", comment: "k = IBAN check digits (always 13) b = National bank code s = Branch code (fr:code guichet) c = Account number x = National check digits (fr:clé RIB)", standard: true},
	"MU": country{name: "Mauritius", size: 30, bban: "4a,19n,3a", code: "MU", format: "MUkk bbbb bbss cccc cccc cccc 000d dd", comment: "b = National bank code s = Branch identifier c = Account number 0 = Zeroes d = Currency Symbol ", standard: true},
	"MC": country{name: "Monaco", size: 27, bban: "10n,11c,2n", code: "MC", format: "MCkk bbbb bsss sscc cccc cccc cxx", comment: "b = National bank code s = Branch code (fr:code guichet) c = Account number x = National check digits (fr:clé RIB). ", standard: true},
	"MD": country{name: "Moldova", size: 24, bban: "2c,18c", code: "MD", format: "MDkk bbcc cccc cccc cccc cccc", comment: "b = National bank code c = Account number", standard: true},
	"ME": country{name: "Montenegro", size: 22, bban: "18n", code: "ME", format: "MEkk bbbc cccc cccc cccc xx", comment: "k = IBAN check digits (always = '25') b = Bank code c = Account number x = National check digits", standard: true},
	"NL": country{name: "Netherlands", size: 18, bban: "4a,10n", code: "NL", format: "NLkk bbbb cccc cccc cc", comment: "b = BIC Bank code c = Account number", standard: true},
	"NO": country{name: "Norway", size: 15, bban: "11n", code: "NO", format: "NOkk bbbb cccc ccx", comment: "b = National bank code c = Account number x = Modulo-11 national check digit", standard: true},
	"PK": country{name: "Pakistan", size: 24, bban: "4c,16n", code: "PK", format: "PKkk bbbb cccc cccc cccc cccc", comment: "b = National bank code c = Account number", standard: true},
	"PS": country{name: "Palestinian territories", size: 29, bban: "4c,21n", code: "PS", format: "PSkk bbbb xxxx xxxx xccc cccc cccc c", comment: "b = National bank code c = Account number x = Not specified", standard: true},
	"PL": country{name: "Poland", size: 28, bban: "24n", code: "PL", format: "PLkk bbbs sssx cccc cccc cccc cccc", comment: "b = National bank code s = Branch code x = National check digit c = Account number, ", standard: true},
	"PT": country{name: "Portugal", size: 25, bban: "21n", code: "PT", format: "PTkk bbbb ssss cccc cccc cccx x", comment: "k = IBAN check digits (always = '50') b = National bank code s = Branch code C = Account number x = National check digit", standard: true},
	"QA": country{name: "Qatar", size: 29, bban: "4a, 21c", code: "QA", format: "QAkk bbbb cccc cccc cccc cccc cccc c", comment: "b = National bank code c = Account number[34]", standard: true},
	"RO": country{name: "Romania", size: 24, bban: "4a,16c", code: "RO", format: "ROkk bbbb cccc cccc cccc cccc", comment: "b = BIC Bank code c = Branch code and account number (bank-specific format) ", standard: true},
	"SM": country{name: "San Marino", size: 27, bban: "1a,10n,12c", code: "SM", format: "SMkk xaaa aabb bbbc cccc cccc ccc", comment: "x = Check char (it:CIN) a = National bank code (it:Associazione bancaria italiana or Codice ABI) b = Branch code (it:Coordinate bancarie or CAB – Codice d'Avviamento Bancario) c = Account number", standard: true},
	"SA": country{name: "Saudi Arabia", size: 24, bban: "2n,18c", code: "SA", format: "SAkk bbcc cccc cccc cccc cccc", comment: "b = National bank code c = Account number preceded by zeros, if required", standard: true},
	"RS": country{name: "Serbia", size: 22, bban: "18n", code: "RS", format: "RSkk bbbc cccc cccc cccc xx", comment: "b = National bank code c = Account number x = Account check digits", standard: true},
	"SK": country{name: "Slovakia", size: 24, bban: "20n", code: "SK", format: "SKkk bbbb ssss sscc cccc cccc", comment: "b = National bank code s = Account number prefix c = Account number", standard: true},
	"SI": country{name: "Slovenia", size: 19, bban: "15n", code: "SI", format: "SIkk bbss sccc cccc cxx", comment: "k = IBAN check digits (always = '56') b = National bank code s = Branch code c = Account number x = National check digits", standard: true},
	"ES": country{name: "Spain", size: 24, bban: "20n", code: "ES", format: "ESkk bbbb gggg xxcc cccc cccc", comment: "b = National bank code g = Branch code x = Check digits c = Account number", standard: true},
	"SE": country{name: "Sweden", size: 24, bban: "20n", code: "SE", format: "SEkk bbbc cccc cccc cccc cccc", comment: "b = National bank code c = Account number ", standard: true},
	"CH": country{name: "Switzerland", size: 21, bban: "5n,12c", code: "CH", format: "CHkk bbbb bccc cccc cccc c", comment: "b = National bank code c = Account number", standard: true},
	"TN": country{name: "Tunisia", size: 24, bban: "20n", code: "TN", format: "TNkk bbss sccc cccc cccc cccc", comment: "k = IBAN check digits (always 59) b = National bank code s = Branch code c = Account number", standard: true},
	"TR": country{name: "Turkey", size: 26, bban: "5n,17c", code: "TR", format: "TRkk bbbb bxcc cccc cccc cccc cc", comment: "b = National bank code x = Reserved for future use (currently '0') c = Account number", standard: true},
	"GB": country{name: "United Kingdom", size: 22, bban: "4a,14n", code: "GB", format: "GBkk bbbb ssss sscc cccc cc", comment: "b = BIC bank code s = Bank and branch code (sort code) c = Account number", standard: true},
	"AE": country{name: "United Arab Emirates", size: 23, bban: "3n,16n", code: "AE", format: "AEkk bbbc cccc cccc cccc ccc", comment: "b = National bank code c = Account number ", standard: true},
	"VG": country{name: "Virgin Islands, British", size: 24, bban: "4c,16n", code: "VG", format: "VGkk bbbb cccc cccc cccc cccc", comment: "b = National bank code c = Account number", standard: true},
}

package cmds

import (
	"bytes"
	"fmt"
	"sort"

	"github.com/pobri19/sqlboiler/dbdrivers"
)

func combineImports(a, b imports) imports {
	var c imports

	c.standard = removeDuplicates(combineStringSlices(a.standard, b.standard))
	c.thirdparty = removeDuplicates(combineStringSlices(a.thirdparty, b.thirdparty))

	sort.Sort(c.standard)
	sort.Sort(c.thirdparty)

	return c
}

func combineTypeImports(a imports, b map[string]imports, columns []dbdrivers.Column) imports {
	tmpImp := imports{
		standard:   make(importList, len(a.standard)),
		thirdparty: make(importList, len(a.thirdparty)),
	}

	copy(tmpImp.standard, a.standard)
	copy(tmpImp.thirdparty, a.thirdparty)

	for _, col := range columns {
		for key, imp := range b {
			if col.Type == key {
				tmpImp.standard = append(tmpImp.standard, imp.standard...)
				tmpImp.thirdparty = append(tmpImp.thirdparty, imp.thirdparty...)
			}
		}
	}

	tmpImp.standard = removeDuplicates(tmpImp.standard)
	tmpImp.thirdparty = removeDuplicates(tmpImp.thirdparty)

	sort.Sort(tmpImp.standard)
	sort.Sort(tmpImp.thirdparty)

	return tmpImp
}

func buildImportString(imps imports) []byte {
	stdlen, thirdlen := len(imps.standard), len(imps.thirdparty)
	if stdlen+thirdlen < 1 {
		return []byte{}
	}

	if stdlen+thirdlen == 1 {
		var imp string
		if stdlen == 1 {
			imp = imps.standard[0]
		} else {
			imp = imps.thirdparty[0]
		}
		return []byte(fmt.Sprintf("import %s", imp))
	}

	buf := &bytes.Buffer{}
	buf.WriteString("import (")
	for _, std := range imps.standard {
		fmt.Fprintf(buf, "\n\t%s", std)
	}
	if stdlen != 0 && thirdlen != 0 {
		buf.WriteString("\n")
	}
	for _, third := range imps.thirdparty {
		fmt.Fprintf(buf, "\n\t%s", third)
	}
	buf.WriteString("\n)\n")

	return buf.Bytes()
}

func combineStringSlices(a, b []string) []string {
	c := make([]string, len(a)+len(b))
	if len(a) > 0 {
		copy(c, a)
	}
	if len(b) > 0 {
		copy(c[len(a):], b)
	}

	return c
}

func removeDuplicates(dedup []string) []string {
	if len(dedup) <= 1 {
		return dedup
	}

	for i := 0; i < len(dedup)-1; i++ {
		for j := i + 1; j < len(dedup); j++ {
			if dedup[i] != dedup[j] {
				continue
			}

			if j != len(dedup)-1 {
				dedup[j] = dedup[len(dedup)-1]
				j--
			}
			dedup = dedup[:len(dedup)-1]
		}
	}

	return dedup
}

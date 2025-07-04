package drivers

import (
	"regexp"

	"github.com/aarondl/strmangle"
)

var rgxEnum = regexp.MustCompile(`^enum(\.\w+)?\([^)]+\)$`)

// Column holds information about a database column.
// Types are Go types, converted by TranslateColumnType.
type Column struct {
	Name          string `json:"name" toml:"name"`
	Type          string `json:"type" toml:"type"`
	DBType        string `json:"db_type" toml:"db_type"`
	Default       string `json:"default" toml:"default"`
	Comment       string `json:"comment" toml:"comment"`
	Nullable      bool   `json:"nullable" toml:"nullable"`
	Unique        bool   `json:"unique" toml:"unique"`
	Validated     bool   `json:"validated" toml:"validated"`
	AutoGenerated bool   `json:"auto_generated" toml:"auto_generated"`

	// Postgres only extension bits
	// ArrType is the underlying data type of the Postgres
	// ARRAY type. See here:
	// https://www.postgresql.org/docs/9.1/static/infoschema-element-types.html
	ArrType *string `json:"arr_type" toml:"arr_type"`
	UDTName string  `json:"udt_name" toml:"udt_name"`
	// DomainName is the domain type name associated to the column. See here:
	// https://www.postgresql.org/docs/10/extend-type-system.html#EXTEND-TYPE-SYSTEM-DOMAINS
	DomainName *string `json:"domain_name" toml:"domain_name"`

	// MySQL only bits
	// Used to get full type, ex:
	// tinyint(1) instead of tinyint
	// Used for "tinyint-as-bool" flag
	FullDBType string `json:"full_db_type" toml:"full_db_type"`
}

// ColumnNames of the columns.
func ColumnNames(cols []Column) []string {
	names := make([]string, len(cols))
	for i, c := range cols {
		names[i] = c.Name
	}

	return names
}

// ColumnDBTypes of the columns.
func ColumnDBTypes(cols []Column) map[string]string {
	types := map[string]string{}

	for _, c := range cols {
		types[strmangle.TitleCase(c.Name)] = c.DBType
	}

	return types
}

// FilterColumnsByAuto generates the list of columns that have autogenerated values
func FilterColumnsByAuto(auto bool, columns []Column) []Column {
	var cols []Column

	for _, c := range columns {
		if (auto && c.AutoGenerated) || (!auto && !c.AutoGenerated) {
			cols = append(cols, c)
		}
	}

	return cols
}

// FilterColumnsByDefault generates the list of columns that have default values
func FilterColumnsByDefault(defaults bool, columns []Column) []Column {
	var cols []Column

	for _, c := range columns {
		if (defaults && len(c.Default) != 0) || (!defaults && len(c.Default) == 0) {
			cols = append(cols, c)
		}
	}

	return cols
}

// FilterColumnsByEnum generates the list of columns that are enum values.
func FilterColumnsByEnum(columns []Column) []Column {
	var cols []Column

	for _, c := range columns {
		if rgxEnum.MatchString(c.DBType) {
			cols = append(cols, c)
		}
	}

	return cols
}

// IsEnumDBType reports whether the column type is Enum
func IsEnumDBType(dbType string) bool {
	return rgxEnum.MatchString(dbType)
}

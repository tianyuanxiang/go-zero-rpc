package dict

const (
	Other  = 0
	Insert = 1
	Update = 2
	Delete = 3
	Query  = 4
)

var OperNumToType = map[int64]string{
	Other:  "Other",
	Insert: "Insert",
	Update: "Update",
	Delete: "Delete",
	Query:  "Query",
}

var OperTypeToNum = map[string]int{
	"Other":  Other,
	"Insert": Insert,
	"Update": Update,
	"Delete": Delete,
	"Query":  Query,
}

package examples

// The migrate command pushes the schema from this file to the database
//go:generate agoradb migrate --db test_users --auth_token $(agoradb whoami -t) schema.go

// The generate command pulls down the schema from the database and generates the client code
//go:generate agoradb generate --db test_users --auth_token $(agoradb whoami -t) --out-dir ./usersdb --package users

// User with related objects: 1:1, 1:n, n:n relations
type User struct {	
	Id			string	`db:"primary_key,type=uuid"`
	FirstName	string
	LastName	string
	Gender		string
	Age			int
	
	// Relational data
	Friends		[]User	`db:"manyToMany=users,by=Id"`
	Address		Address	`db:"one=user_addresses,by=Id,having=mainAddress;true"`
	Cars		[]Car	`db:"many=userCars,by=Id"`
}

// Address with composite primary key and a custom table name
type Address struct {
	TableName 			`db:"user_addresses"`
	
	UserId		string	`db:"primary_key,type=uuid"`
	MainAddress bool	`db:"primary_key,type=bit"`
	Street		string
	City		string
	ZipCode		int
	
	// Relational data
	User User `db:"one=users,by=UserId"`
}

// Car with composite primary key
type Car struct {
	UserId		string	`db:"primary_key,type=uuid"`
	Model		string	`db:"primary_key"`
	Color		Colors	`db:"enum"`
	
	// Relational data
	User User `db:"one=users,by=UserId"`
}

type Colors int

const (
	Red Colors = iota
	Green
	Blue
	Black
)

func (c Colors) String() string {
	return [...]string{"Red", "Green", "Blue", "Black"}[c]
}
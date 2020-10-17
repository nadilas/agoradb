type User struct {	
	Id			string	`db:primary_key,type=uuid`
	FirstName	string
	LastName	string
	Gender		string
	Age			int
	
	// Relational data
	Friends		[]User	`db:manyToMany=users,by=Id,via=friends.friend_id`
	Address		Address	`db:one=user_addresses,by=Id,having=mainAddress,equal=true`
	Cars		[]Car	`db:many=userCars,by=Id`
}

type Address struct {
	TableName 			`db:user_addresses`
	
	UserId		string	`db:primary_key,type=uuid`
	MainAddress bool	`db:primary_key,type=bit`
	Street		string
	City		string
	ZipCode		int
	
	// Relational data
	User User `db:one=users,by=UserId`
}

type Car struct {
	UserId		string	`db:primary_key,type=uuid`
	Model		string	`db:primary_key`
	Color		Colors
	
	// Relational data
	User User `db:one=users,by=UserId`
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
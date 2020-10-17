import (
	agoradb "github.com/featme-inc/agoradb/client"
	q "github.com/featme-inc/agoradb/client/query"
	users "github.com/featme-inc/agoradb/examples/usersdb"
)

func main() {
	// Connect with user & pass using the agoradb library
	client, err := agoradb.New(agoradb.Config{
		Auth: agoradb.Auth{
			User: username,
			Password: password,
		},
	})
	if err != nil {
		log.Fatal(err)
	}
	
	// After having generated the client code (go:generate or agoradb client generate ...)
	// Use the generated client code to access the database
	db := users.NewUsersDatabase(client)
	
	db.GetUser(
		q.And(users.User{Firstname: q.Like("Joh%")}, users.User{Age: q.Eq(22)}),
	).
	// use middlewares to preload relationships
	With(
		users.WithFriends(), users.WithAddress(), users.WithCars(),
	).
	// use the finishers to get the execute the query
	// the finishers such as One(), First(int), All(), StreamOne(), StreamFirst(int), StreamAll() are the
	// ones which actually execute the grpc query to the backend.
	// The grpc response is interpreted and turned into a local representation before the finisher returns
	StreamAll(func (users []User, err error) {
		// use response ....
		for _, u := range users {
			log.Println(u.Fristname + " has " + len(u.Cars) + " car(s).")
		}
	})
}
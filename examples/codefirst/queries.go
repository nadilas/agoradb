package examples

import (
	"context"
	
	agoradb "github.com/featme-inc/agoradb/client"
	q "github.com/featme-inc/agoradb/client/query"
	"github.com/featme-inc/agoradb/examples/codefirst/usersdb"
)

func main() {
	client, err := agoradb.New(agoradb.Config{
		Auth: agoradb.Auth{
			User: username,
			Password: password,
		},
	})
	if err != nil {
		log.Fatal(err)
	}
	
	db := users.NewDatabase(client)
	
	getUsers(db)
}

// Simple read with preloading of relational data
//
// Traditional pseudo SQL representation would be:
// select * from users as USERS where u.firstname like 'Joh%' and u.age = 22
//. + select * from user_addresses as UA where UA.user_id in (USERS...)
//. + select * from cars as UC where UC.user_id in (USERS...)
//. + select * from users as FRIENDS where FRIENDS.id in (select UF.friend_id from users_friends as UF where UF.user_id in (USERS...))
func getUsersContext(ctx context.Context, db *UsersDatabase) {
	db.GetUserContext(
		ctx,
		q.And(users.UserQ{Firstname: q.Like("Joh%")}, users.UserQ{Age: q.Eq(22)}),
	).
	// use middlewares to preload selected relationships
	With(
		users.WithUserFriends(),
		users.WithUserAddress(),
		users.WithUserCars(),
	).
	// use the finishers to get the execute the query
	// the finishers such as One(), First(int), All(), StreamOne(), StreamFirst(int), StreamAll() are the
	// ones which actually execute the grpc query to the backend.
	// The grpc response is interpreted and turned into a local DTO representation before the finisher returns
	StreamAll(func (users []users.User, err error) error {
		// use response ....
		for _, u := range users {
			log.Println(u.Fristname + " has " + len(u.Cars) + " car(s) and " + len(u.Friends) + " friend(s).")
		}
		return nil
	})
}
package users

import (
	q "github.com/featme-inc/agoradb/client/query"
)

// User query helper

type UserQ struct {
	Id			q.StringQuery
	FirstName	q.StringQuery
	LastName	q.StringQuery
	Gender		q.StringQuery
	Age			q.IntQuery
}

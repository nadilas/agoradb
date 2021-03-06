package users

// TODO THIS FILE SHOULD BE GENERATED by agoradb migrage

import (
	"context"
	
	q "github.com/featme-inc/agoradb/client/query"
)

??? which User struct to use? Own representation from the protobuf generated code or use the one from the original file (is that possible at all)?

type GetUserFinisher interface {
	// Returns exactly one User
	// If the query returns more than one rows, this method returns with an error
	One() (user User, err error)
	
	// Returns the first n Users
	First(n int) (users []User, err error)
	
	// Returns all Users
	All() (users []User, err error)
	
	// Returns exactly one User asynchronously
	// If the query returns more than one rows, this method returns with an error
	StreamOne(func (user User, err error) error)
	
	// Returns the first n Users asynchronously
	StreamFirst(int, func (users []User, err error) error)

	// Returns all Users asynchronously
	StreamAll(func (users []User, err error) error)
}

type GetUserMiddleware interface {
	With(middleware ...q.With) GetUserFinisher
}

type SaveUserMiddleware interface {
	
}

type UsersDatabase interface {
	GetUser(q ...q.Query) GetUserMiddleware
	GetUserContext(ctx context.Context, q ...q.Query) GetUserMiddleware
	SaveUser(users ...User) SaveUserMiddleware
}

type usersDatabase struct {
	client *agoradb.Client
}

func NewDatabase(client *agoradb.Client) UsersDatabase {
	return &usersDatabase{client: client}
}

func (db *userDatabase) GetUserContext(ctx context.Context, q ...q.Query) GetUserMiddleware {
	return &getUserQuery{ctx: ctx, c: db.client, q: q...}
}

type getUserQuery {
	ctx context.Context
	c *agoradb.Client
	q []q.Query
	m []q.Middleware
}

func (c *getUserQuery) buildQuery() []*GetUserQueryParams {
	// todo turn the middleware input into wire compatible 
	var qp []*GetUserQueryParams
	for _, q := range c.q {
		mp = append(qp, q.Wire())
	}
	return mp
}

func (c *getUserQuery) buildMiddlewares() []*GetUserMiddlewareParams {
	// todo turn the middleware input into wire compatible 
	var mp []*GetUserMiddlewareParams
	for _, mw := range c.m {
		mp = append(mp, mw.Wire())
	}
	return mp
}

func (c *getUserQuery) With(middleware ...q.With) GetUserFinisher {
	c.m = append(c.m, middleware...)
	return c
}

// Returns exactly one User
// If the query returns more than one rows, this method returns with an error
func (c *getUserQuery) One() (user User, err error) {
	resp, err := c.c.GetUsers(c.ctx, &GetUsersRequest{query: c.buildQuery(), middlewares: c.buildMiddlewares()})
	if err != nil {
		return User{}, err
	}
	
	resc := len(resp.Users)
	if resc < 1 {
		return User{}, q.NoResultsFound
	} else if resc > 1 {
		return User{}, q.TooManyResults
	}
	
	return resp.Users[0], nil
}

// Returns the first n Users
//
// n=0 means all entries
func (c *getUserQuery) First(n int) (users []User, err error) {
	if n == 0 {
		return c.All()
	}
	
	resp, err := c.c.GetUsers(c.ctx, &GetUsersRequest{query: c.buildQuery(), middlewares: c.buildMiddlewares()})
	if err != nil {
		return User{}, err
	}
	
	resc := len(resp.Users)
	if resc < 1 {
		return User{}, q.NoResultsFound
	}
	
	max := n-1
	if max > resc {
		max = resc - 1
	}
	
	return resp.Users[0..max], nil
}

// Returns all Users
func (c *getUserQuery) All() (users []User, err error) {
	resp, err := c.c.GetUsers(c.ctx, &GetUsersRequest{query: c.buildQuery(), middlewares: c.buildMiddlewares()})
	if err != nil {
		return User{}, err
	}
	
	resc := len(resp.Users)
	if resc < 1 {
		return User{}, q.NoResultsFound
	}
	
	return resp.Users, nil
}

// Returns exactly one User asynchronously
// If the query returns more than one rows, this method returns with an error
func (c *getUserQuery) StreamOne(fn func (user User, err error) error) {
	stream, err := c.c.StreamUsers(c.ctx, &StreamUsersRequest{query: c.buildQuery(), middlewares: c.buildMiddlewares()})
	if err != nil {
	  _ = fn(nil, err)
	  return
	}
	c := 0
	for {
	    user, err := stream.Recv()
	    if err == io.EOF {
			_ = fn(io.EOF)
	        break
	    }
		if err != nil {
	        log.Println("ERR: %v.StreamUsers(_) = _, %v", c.c, err)
			break
	    }
		c++
		if c == 2 {
			_ = fn(User{}, q.TooManyResults)
			break
		}
	    err = fn(user, nil)
		if err != nil {
			log.Println("ERR: Cancelled %v.StreamUsers(_) = _, %v", c.c, err)
			break
		}
	}
}

// Returns the first n Users asynchronously
func (c *getUserQuery) StreamFirst(n int, fn func (users []User, err error) error) {
	if n == 0 {
		c.StreamAll(fn)
		return
	}
	
	ctx, cancel := context.WithCancel(c.ctx)
	stream, err := c.c.StreamUsers(ctx, &StreamUsersRequest{query: c.buildQuery(), middlewares: c.buildMiddlewares()})
	if err != nil {
	  _ = fn(nil, err)
	  return
	}
	objCounter := 0
	for {
	    user, err := stream.Recv()
	    if err == io.EOF {
			_ = fn(io.EOF)
	        break
	    }
		if err != nil {
	        log.Println("ERR: %v.StreamUsers(_) = _, %v", c.c, err)
			break
	    }
		objCounter++
	    err = fn(user, nil)
		if err != nil {
			log.Println("ERR: Cancelled %v.StreamUsers(_) = _, %v", c.c, err)
			break
		}
		if objCounter+1 > n {
			cancel()
			break
		}
	}
}

func (c *getUserQuery) StreamAll(fn func (users User, err error) error) {
	stream, err := c.c.StreamUsers(c.ctx, &StreamUsersRequest{query: c.buildQuery(), middlewares: c.buildMiddlewares()})
	if err != nil {
	  _ = fn(nil, err)
	  return
	}
	for {
	    user, err := stream.Recv()
	    if err == io.EOF {
			_ = fn(io.EOF)
	        break
	    }
	    if err != nil {
	        log.Println("ERR: %v.StreamUsers(_) = _, %v", c.c, err)
			break
	    }
	    err = fn(user, nil)
		if err != nil {
			log.Println("ERR: Cancelled %v.StreamUsers(_) = _, %v", c.c, err)
			break
		}
	}
}
package surreal

import (
	"time"

	"github.com/gookit/color"
	"github.com/surrealdb/surrealdb.go"
)

type Surreal struct {
	connAttempts int
	connTimeout  time.Duration
	*surrealdb.DB
}

func New(url, ns, database, user, password string) (*Surreal, error) {
	surreal := &Surreal{
		connAttempts: 10,
		connTimeout:  time.Second,
	}

	var lastErr error

	for surreal.connAttempts > 0 {
		db, err := surrealdb.New(url)
		if err != nil {
			surreal.connAttempts--

			lastErr = err

			color.Redln("Surreal is trying to connect, attempts left: ", surreal.connAttempts)

			time.Sleep(surreal.connTimeout)

			continue
		}

		_, err = db.Signin(map[string]interface{}{
			"user": user,
			"pass": password,
		})
		if err != nil {
			surreal.connAttempts--

			lastErr = err

			color.Redln("Surreal is trying to connect, attempts left: ", surreal.connAttempts)

			time.Sleep(surreal.connTimeout)

			continue
		}

		_, err = db.Use(ns, database)
		if err != nil {
			surreal.connAttempts--

			lastErr = err

			color.Redln("Surreal is trying to connect, attempts left: ", surreal.connAttempts)

			time.Sleep(surreal.connTimeout)

			continue
		}

		surreal.DB = db

		break
	}

	if lastErr != nil {
		return nil, lastErr
	}

	return surreal, nil
}

func (r *Surreal) Close() {
	r.DB.Close()
}

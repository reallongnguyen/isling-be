package surreal

import (
	"isling-be/pkg/logger"
	"time"

	"github.com/surrealdb/surrealdb.go"
)

type Surreal struct {
	connAttempts int
	connTimeout  time.Duration
	log          logger.Interface
	*surrealdb.DB
}

func New(url, ns, database, user, password string, log logger.Interface) (*Surreal, error) {
	surreal := &Surreal{
		connAttempts: 10,
		connTimeout:  time.Second,
		log:          log,
	}

	err := surreal.Connect(url, ns, database, user, password)

	return surreal, err
}

func (r *Surreal) Connect(url, ns, database, user, password string) error {
	var lastErr error

	for r.connAttempts > 0 {
		db, err := surrealdb.New(url)
		if err != nil {
			r.connAttempts--

			lastErr = err

			r.log.Error("Surreal is trying to connect, attempts left: %d", r.connAttempts)

			time.Sleep(r.connTimeout)

			continue
		}

		_, err = db.Signin(map[string]interface{}{
			"user": user,
			"pass": password,
		})
		if err != nil {
			r.connAttempts--

			db.Close()

			lastErr = err

			r.log.Error("Surreal is trying to connect, attempts left: %d", r.connAttempts)

			time.Sleep(r.connTimeout)

			continue
		}

		_, err = db.Use(ns, database)
		if err != nil {
			r.connAttempts--

			lastErr = err

			r.log.Error("Surreal is trying to connect, attempts left: %d", r.connAttempts)

			time.Sleep(r.connTimeout)

			continue
		}

		r.DB = db
		lastErr = nil

		break
	}

	return lastErr
}

func (r *Surreal) Close() {
	if r.DB != nil {
		r.DB.Close()
	}
}

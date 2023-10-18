package watermill

import (
	"context"
	"isling-be/pkg/logger"

	"github.com/ThreeDotsLabs/watermill"
	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/ThreeDotsLabs/watermill/pubsub/gochannel"
)

type Watermill struct {
	Router      *message.Router
	internalPub message.Publisher
	internalSub message.Subscriber
}

func NewWatermill(l logger.Interface) *Watermill {
	logAdp := &LogAdapter{log: l}

	// router, err := message.NewRouter(message.RouterConfig{}, logAdp)
	// if err != nil {
	// 	l.Fatal("NewWatermill: %w", err)
	// }

	// router.AddPlugin(plugin.SignalsHandler)

	// router.AddMiddleware(
	// 	middleware.CorrelationID,
	// 	middleware.Retry{
	// 		MaxRetries:      5,
	// 		InitialInterval: time.Millisecond * 100,
	// 		Multiplier:      2,
	// 		Logger:          logAdp,
	// 	}.Middleware,
	// 	middleware.Recoverer,
	// )

	goChan := gochannel.NewGoChannel(gochannel.Config{
		OutputChannelBuffer:            1024,
		BlockPublishUntilSubscriberAck: false,
	}, logAdp)

	return &Watermill{
		// Router:      router,
		internalPub: goChan,
		internalSub: goChan,
	}
}

// func (r *Watermill) RunRouter() {
// 	go func() {
// 		if err := r.Router.Run(context.Background()); err != nil {
// 			log.Fatal("RunRouter: run watermill router: %w", err)
// 		}
// 	}()

// 	<-r.Router.Running()
// }

func (r *Watermill) Publish(topic string, payload []byte, metadata map[string]string) error {
	msg := message.NewMessage(watermill.NewShortUUID(), payload)

	for k, v := range metadata {
		msg.Metadata.Set(k, v)
	}

	return r.internalPub.Publish(topic, msg)
}

func (r *Watermill) Subscribe(topic string, handler func(uuid string, payload []byte, metadata map[string]string) error) error {
	msgChan, err := r.internalSub.Subscribe(context.Background(), topic)
	if err != nil {
		return err
	}

	go func() {
		for msg := range msgChan {
			err := handler(msg.UUID, msg.Payload, msg.Metadata)

			if err == nil {
				msg.Ack()
			} else {
				msg.Nack()
			}
		}
	}()

	return nil
}

func (r *Watermill) Close() {
	r.internalPub.Close()
	r.internalSub.Close()
}

package queuer

import (
	"encoding/json"

	etcd "github.com/coreos/etcd/clientv3"

	"github.com/tddhit/tools/log"
	"github.com/tddhit/tools/msgqueue"
	"github.com/tddhit/tools/msgqueue/option"
	"github.com/tddhit/zerg/types"
)

type NsqQueuer struct {
	producer *msgqueue.NsqProducer
	consumer *msgqueue.NsqConsumer
	msgChan  chan *msgqueue.Message
	topics   []string
}

func NewNsqQueuer(client *etcd.Client, topics []string,
	popt option.NsqProducer, copt option.NsqConsumer) (*NsqQueuer, error) {
	producer, err := msgqueue.NewNsqProducer(client, popt)
	if err != nil {
		log.Error(err)
		return nil, err
	}
	consumer, err := msgqueue.NewNsqConsumer(client, copt)
	if err != nil {
		log.Error(err)
		return nil, err
	}
	q := &NsqQueuer{
		producer: producer,
		consumer: consumer,
		msgChan:  make(chan *msgqueue.Message, 1000),
		topics:   topics,
	}
	q.pop()
	return q, nil
}

func (q *NsqQueuer) Push(req *types.Request) {
	if err := q.producer.Publish(req.Parser, []byte(req.RawURL)); err != nil {
		log.Error(err)
	}
}

func (q *NsqQueuer) Pop() (req *types.Request) {
	var jsonMsg struct {
		URL    string
		Parser string
		Proxy  string
	}
	for {
		msg := <-q.msgChan
		if err := json.Unmarshal(msg.Body, &jsonMsg); err != nil {
			log.Error(err)
			continue
		}
		req, err := types.NewRequest(jsonMsg.URL, jsonMsg.Parser, jsonMsg.Proxy)
		if err != nil {
			log.Error(err)
			continue
		}
		return req
	}
	return nil
}

func (q *NsqQueuer) pop() {
	for _, topic := range q.topics {
		go func(topic string) {
			for msg := range q.consumer.Messages(topic) {
				q.msgChan <- msg
			}
		}(topic)
	}
}

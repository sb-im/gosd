package luavm

import (
	"context"
	"encoding/json"
	"fmt"
	"testing"
	"time"

	"sb.im/gosd/app/model"
)

func TestLuaNotification(t *testing.T) {
	task := &model.Task{}
	task.ID = 1

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	w := newWorker(t)
	ch := make(chan error)
	go func() {
		time.Sleep(100 * time.Millisecond)
		ch <- w.doRun(task, []byte(`
function main(task)
  print("### RUN Notification RUN ###")

  task:Notification("notification")

  sleep("1s")

  task:Notification("notification", 3)

  print("### END Notification END ###")
end
`))
	}()

	topic := fmt.Sprintf(topic_notification, task.ID)
	keyspace := "__keyspace@*__:%s"
	pubsub := w.rdb.PSubscribe(ctx, fmt.Sprintf(keyspace, topic))
	ev := pubsub.Channel()

	// Notification
	<-ev
	if val, err := w.rdb.Get(ctx, topic).Bytes(); err != nil {
		t.Error(err)
	} else {
		d := Notification{}
		if err = json.Unmarshal(val, &d); err != nil {
			t.Error(err)
		}
	}

	// Notification 2
	<-ev
	if val, err := w.rdb.Get(ctx, topic).Bytes(); err != nil {
		t.Error(err)
	} else {
		d := Notification{}
		if err = json.Unmarshal(val, &d); err != nil {
			t.Error(err)
		}
	}

	pubsub.Close()

	select {
	case <-ctx.Done():
		t.Error("Time Out")
	case err := <-ch:
		if err != nil {
			t.Error(err)
		}
	}
}

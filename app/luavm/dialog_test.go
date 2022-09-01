package luavm

import (
	"context"
	"encoding/json"
	"fmt"
	"testing"
	"time"

	"sb.im/gosd/app/model"
)

func TestLuaDialog(t *testing.T) {
	task := &model.Task{}
	task.ID = 1

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	w := newWorker(t)
	ch := make(chan error)
	go func() {
		time.Sleep(100 * time.Millisecond)
		ch <- w.doRun(context.Background(), task, []byte(`
function main(task)
  print("### RUN Dialog RUN ###")
  task:CleanDialog()
  sleep("1s")

  dialog = {
    name = "Checker ~",
    message = "Wow Wow Wow ~",
    level = "success",
    items = {
      {name = "nn", message = 'mm', level = 'info'},
      {name = "n2", message = 'ok', level = 'success'},
      {name = "n3", message = 'Not ok', level = 'danger'},
      {name = "n4", message = '...', level = 'warning'},
    },
    buttons = {
      {name = "Cancel", message = 'cancel', level = 'primary'},
      {name = "Confirm", message = 'confirm', level = 'danger'},
    }
  }

	task:ToggleDialog(dialog)
  sleep("1s")

  ask_status = {
    name = "ARE YOU OK ?",
    buttons = {
      {name = "Fine, thank you.", message = 'fine', level = 'primary'},
      {name = "I feel bad.", message = 'bad', level = 'danger'},
    }
  }

  task:ToggleDialog(ask_status)
  sleep("1s")

  task:CleanDialog()
  print("### END Dialog END ###")
end
`))
	}()

	topic := fmt.Sprintf(topic_dialog, task.ID)
	keyspace := "__keyspace@*__:%s"
	pubsub := w.rdb.PSubscribe(ctx, fmt.Sprintf(keyspace, topic))
	ev := pubsub.Channel()

	// CleanDialog
	<-ev
	if val, err := w.rdb.Get(ctx, topic).Result(); err != nil {
		t.Error(err)
	} else if val != "{}" {
		t.Error(val)
	}

	// ToggleDialog
	<-ev
	if val, err := w.rdb.Get(ctx, topic).Bytes(); err != nil {
		t.Error(err)
	} else {
		d := Dialog{}
		if err = json.Unmarshal(val, &d); err != nil {
			t.Error(err)
		}
	}

	// ToggleDialog 2
	<-ev
	if val, err := w.rdb.Get(ctx, topic).Bytes(); err != nil {
		t.Error(err)
	} else {
		d := Dialog{}
		if err = json.Unmarshal(val, &d); err != nil {
			t.Error(err)
		}
	}

	// CleanDialog 2
	<-ev
	if val, err := w.rdb.Get(ctx, topic).Result(); err != nil {
		t.Error(err)
	} else if val != "{}" {
		t.Error(val)
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

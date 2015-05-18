package dynatest

import (
	"github.com/underarmour/dynago"
	"testing"
)

func TestExecutor(t *testing.T) {
	client := dynago.NewClientExecutor(&Executor{})
	client.Query("foo") // TODO
}

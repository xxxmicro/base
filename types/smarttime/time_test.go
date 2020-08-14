package smarttime_test

import(
	"testing"
	"time"
	"github.com/xxxmicro/base/types/smarttime"
)

type User struct {
	ID string				`json:"id"`
	Ctime smarttime.Time			`json:"ctime"`
}

func TestSmartTime(t *testing.T) {
	time.Now()

	user := &User{
		ID: "1",
	}

	var err error
	user.Ctime, err = smarttime.Parse("1597401486073")
	if err != nil {
		t.Fatal(err)
		return
	}
	t.Logf("ctime: %s", time.Time(user.Ctime).Format("2006-01-02 15:04"))
}

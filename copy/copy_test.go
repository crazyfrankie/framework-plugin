package copy

import (
	"testing"
)

//func TestCopy(t *testing.T) {
//	a := A{
//		Name:  "tom",
//		Age:   10,
//		Phone: "13117127070",
//	}
//	var b B
//
//	err := NewCopier(a, &b, "Phone").Builder()
//	assert.NoError(t, err)
//
//	fmt.Println(a)
//	fmt.Println(b)
//}

func BenchmarkCopier(t *testing.B) {
	a := A{
		ID:       1,
		Name:     "tom",
		Password: "1234567",
		Avatar:   "github.com/crazyfrankie/static/default.png",
		Phone:    "13117127070",
		Ctime:    11,
		Utime:    11,
	}
	var b B

	for i := 0; i < t.N; i++ {
		err := NewCopier(a, &b, "Password").Builder()
		if err != nil {
			t.Fatal(err)
		}
	}
}

func BenchmarkCopy(t *testing.B) {
	a := A{
		ID:       1,
		Name:     "tom",
		Password: "1234567",
		Avatar:   "github.com/crazyfrankie/static/default.png",
		Phone:    "13117127070",
		Ctime:    11,
		Utime:    11,
	}
	var b B

	for i := 0; i < t.N; i++ {
		b.ID = a.ID
		b.Name = a.Name
		b.Avatar = a.Avatar
		b.Phone = a.Phone
		b.Ctime = a.Ctime
		b.Utime = a.Utime
	}
}

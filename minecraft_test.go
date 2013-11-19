// minecraft_test.go
package minecraft

import (
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestGetUser(t *testing.T) {

	Convey("CLone1018 should equal clone1018", t, func() {
		user := GetUser("CLone1018")

		So(user.Name, ShouldEqual, "clone1018")
	})

}

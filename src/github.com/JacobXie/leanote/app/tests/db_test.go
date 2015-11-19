package tests

import (
	"github.com/JacobXie/leanote/app/db"
	"testing"
	//	. "github.com/JacobXie/leanote/app/lea"
	//	"github.com/JacobXie/leanote/app/service"
	//	"gopkg.in/mgo.v2"
	//	"fmt"
)

func TestDBConnect(t *testing.T) {
	db.Init("mongodb://localhost:27017/leanote", "leanote")
}

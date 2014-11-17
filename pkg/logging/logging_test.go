package logging

import (
	"testing"

	"github.com/Sirupsen/logrus"
	. "github.com/anthonybishopric/gotcha"
)

func TestLoggingCanMergeFields(t *testing.T) {
	fields1 := logrus.Fields{
		"foo": "a",
		"bar": "b",
	}
	fields2 := logrus.Fields{
		"foo": "z",
		"baz": "q",
	}
	res := Merge(fields1, fields2)
	Assert(t).AreEqual("z", res["foo"], "Should have taken new value's foo")
	Assert(t).AreEqual("b", res["bar"], "Should have taken old value's bar")
	Assert(t).AreEqual("q", res["baz"], "Should have taken old value's baz")
}

func TestSubLoggerMergesFields(t *testing.T) {
	logger := NewLogger(logrus.Fields{
		"foo": "a",
		"bar": "b",
	})
	sub := logger.SubLogger(logrus.Fields{
		"foo": "z",
		"baz": "q",
	})
	res := sub.baseFields
	Assert(t).AreEqual("z", res["foo"], "Should have taken new value's foo")
	Assert(t).AreEqual("b", res["bar"], "Should have taken old value's bar")
	Assert(t).AreEqual("q", res["baz"], "Should have taken old value's baz")

	Assert(t).AreEqual("a", logger.baseFields["foo"], "Should not have overwritten original")
}

func TestLoggingMergeDoesNotModifyOriginalMap(t *testing.T) {
	fields1 := logrus.Fields{
		"foo": "a",
	}
	fields2 := logrus.Fields{
		"foo": "b",
	}
	Merge(fields1, fields2)

	Assert(t).AreEqual("a", fields1["foo"], "Should not have modified the original fields")
}

func TestWithFieldsCombinesBaseFieldsAndGiven(t *testing.T) {
	logger := NewLogger(logrus.Fields{
		"foo": "a",
	})
	entry := logger.WithFields(logrus.Fields{
		"baz": "c",
	})

	Assert(t).AreEqual("a", entry.Data["foo"], "should have kept foo")
	Assert(t).AreEqual("c", entry.Data["baz"], "should have merged baz")
}
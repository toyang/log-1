package log

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"strings"
	"testing"
	"time"
)

func TestFileLine(t *testing.T) {
	buf := bytes.NewBuffer(make([]byte, 100))
	l := new(logger)
	l.name = "name"
	l.lv = INFO
	l.handlers = make(map[Handler]bool)

	hdr, err := NewStreamHandler(buf, "{{level}} {{date}} {{time}} {{name}} {{file_line}} {{}}")
	if err != nil {
		t.Fatalf("NewStreamHandler Error:%v", err)
	}
	l.AddHandler(hdr)
	SetGlobalAppID("samaritan.test")
	defer SetGlobalAppID("")
	l.Info("TEST_TEST")

	strs := strings.Split(buf.String(), " ")
	if strs[4] != "log_test.go:27" {
		t.Errorf("FileLine Error: %s", buf.String())
	}
}

func newLogger(t *testing.T, w io.Writer, f string) Logger {
	l := NewWithWriter("test", nil)
	h, err := NewStreamHandler(w, f)
	h.Colored(false)
	l.AddHandler(h)
	if err != nil {
		t.Error("error creating stream handler: ", err)
		t.FailNow()
	}
	return l
}

func TestGlobalLevel(t *testing.T) {
	expected := "W: WarnLog\n"
	var b bytes.Buffer
	l := newLogger(t, &b, "{{l}}: {{}}")
	SetGlobalLevel(WARN)
	defer SetGlobalLevel(NOTSET)
	l.SetLevel(INFO)

	l.Debug("DebugLog")
	l.Info("InfoLog")
	l.Warn("WarnLog")

	if b.String() != expected {
		t.Errorf("Expected:\n%v\nGot:\n%v", expected, b.String())
	}
}

func TestLevel(t *testing.T) {
	expected := "W: WarnLog\n"
	var b bytes.Buffer
	l := newLogger(t, &b, "{{l}}: {{}}")
	l.SetLevel(WARN)
	l.Debug("DebugLog")
	l.Info("InfoLog")
	l.Warn("WarnLog")

	if b.String() != expected {
		t.Errorf("Expected:\n%v\nGot:\n%v", expected, b.String())
	}
}

func TestGlobalAppID(t *testing.T) {
	var buf bytes.Buffer
	l := newLogger(t, &buf, "[{{app_id}}] ## {{}}")

	expected_nil := "[-] ## InfoLog\n"
	l.Info("InfoLog")
	if buf.String() != expected_nil {
		t.Errorf("Expected:\n%s\nGot:\n%s", expected_nil, buf.String())
	}

	buf.Reset()
	expected := "[test.appid] ## InfoLog\n"
	SetGlobalAppID("test.appid")
	defer SetGlobalAppID("")

	l.Info("InfoLog")
	if buf.String() != expected {
		t.Errorf("Expected:\n%s\nGot:\n%s", expected, buf.String())
	}
}

func TestSetRPCID(t *testing.T) {
	var buf bytes.Buffer
	l := newLogger(t, &buf, "[{{rpc_id}}] ## {{}}")
	rpcLog := l.(RPCLogger)

	expected_nil := "[-] ## InfoLog\n"
	rpcLog.Info("InfoLog")
	if buf.String() != expected_nil {
		t.Errorf("Expected:\n%s\nGot:\n%s", expected_nil, buf.String())
	}

	buf.Reset()
	expected := "[test.rpcid] ## InfoLog\n"
	rpcLog.SetRPCID("test.rpcid")
	rpcLog.Info("InfoLog")
	if buf.String() != expected {
		t.Errorf("Expected:\n%s\nGot:\n%s", expected, buf.String())
	}
}

func TestSetRequestID(t *testing.T) {
	var buf bytes.Buffer
	l := newLogger(t, &buf, "[{{request_id}}] ## {{}}")
	rpcLog := l.(RPCLogger)

	expected_nil := "[-] ## InfoLog\n"
	rpcLog.Info("InfoLog")
	if buf.String() != expected_nil {
		t.Errorf("Expected:\n%s\nGot:\n%s", expected_nil, buf.String())
	}

	buf.Reset()
	expected := "[test.request_id] ## InfoLog\n"
	rpcLog.SetRequestID("test.request_id")
	rpcLog.Info("InfoLog")
	if buf.String() != expected {
		t.Errorf("Expected:\n%s\nGot:\n%s", expected, buf.String())
	}
}

func TestTemplate(t *testing.T) {
	expected := `long: INFO
short: I
duplicate: I
content: hi
`
	var b bytes.Buffer
	l := newLogger(t, &b, "long: {{ level }}\nshort: {{     l }}\nduplicate: {{l}}\ncontent: {{}}")
	l.Info("hi")

	if b.String() != expected {
		t.Errorf("Expected:\n%v\nGot:\n%v\n", expected, b.String())
	}
}

func TestSyslogtpl(t *testing.T) {
	buf := bytes.NewBuffer(make([]byte, 100))

	l := new(logger)
	l.name = "name"
	l.lv = INFO
	l.handlers = make(map[Handler]bool)

	hdr, err := NewStreamHandler(buf, syslogTpl)
	if err != nil {
		t.Fatalf("NewStreamHandler Error:%v", err)
	}
	l.AddHandler(hdr)
	SetGlobalAppID("samaritan.test")
	defer SetGlobalAppID("")
	l.Info("TEST_TEST")

	strs := strings.Split(buf.String(), " ")

	if strs[4] != "[samaritan.test" || strs[8] != "TEST_TEST\n" ||
		strs[5] != "-" || strs[6] != "-]" {
		t.Errorf("SyslogTpl Error: %s", buf.String())
	}
}

func ExampleLogger() {
	l := NewWithWriter("test", nil)
	h, err := NewStreamHandler(os.Stdout, "{{level}} {{}}")
	if err != nil {
		fmt.Println("error creating stream handler: ", err)
		return
	}
	h.Colored(false)
	l.AddHandler(h)
	l.Debug("default level is INFO")
	l.Info("thus debug is not printed")

	l.SetLevel(DEBUG)
	l.Debug("this enables debug")
	// Output:
	// INFO thus debug is not printed
	// DEBUG this enables debug
}

func ExampleLevel() {
	l := NewWithWriter("test", nil)
	l.SetLevel(DEBUG)
	h, err := NewStreamHandler(os.Stdout, "{{level}} {{}}")
	if err != nil {
		fmt.Println("error creating stream handler: ", err)
		return
	}
	h.Colored(false)
	l.AddHandler(h)
	l.Debug("Debug, turned off by default")
	l.Info("Info, default log level")
	l.Warn("Warning, errors are handled, need attention")
	// l.Fatal("Fatal, will os.Exit(1)")
	// FATA Fatal, will os.Exit(1)

	// Output:
	// DEBUG Debug, turned off by default
	// INFO Info, default log level
	// WARN Warning, errors are handled, need attention
}

// Benchmarks
func dateM(n time.Time) string {
	year, month, day := n.Date()
	return fmt.Sprintf("%04d-%02d-%02d", year, month, day)
}

func timeM(n time.Time) string {
	hour, min, sec := n.Clock()
	return fmt.Sprintf("%02d:%02d:%02d", hour, min, sec)
}

func BenchmarkDateM(b *testing.B) {
	n := time.Now()
	for i := 0; i < b.N; i++ {
		dateM(n)
	}
}

// func BenchmarkDate(b *testing.B) {
// 	dftLogger.cur.now = time.Now()
// 	for i := 0; i < b.N; i++ {
// 		dftLogger.date(false)
// 	}
// }

func BenchmarkTimeM(b *testing.B) {
	n := time.Now()
	for i := 0; i < b.N; i++ {
		timeM(n)
	}
}

// func BenchmarkTime(b *testing.B) {
// 	dftLogger.cur.now = time.Now()
// 	for i := 0; i < b.N; i++ {
// 		dftLogger.time(false)
// 	}
// }

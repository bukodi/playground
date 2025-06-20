package freejson

import (
	"encoding/json"
	"fmt"
	"os"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"testing"
)

type result struct {
	Level string
	Pkg   string
	Msg   string
	Count string
}

type record struct {
	Preview bool
	Result  result
}

func TestMDASplunk(t *testing.T) {
	data, err := os.ReadFile("/home/lbukodi/Downloads/MDA_client_logs_24_hours.json")
	if err != nil {
		t.Error(err)
	}

	lines := strings.Split(string(data), "\n")
	t.Logf("%d lines", len(lines))
	var errorResults map[string]int = make(map[string]int)
	for i, line := range lines {
		var rec record
		err = json.Unmarshal([]byte(line), &rec)
		if err != nil {
			t.Errorf("can't parse %d line: %s\n%+v", i, line, err)
		}
		count, _ := strconv.Atoi(rec.Result.Count)
		msg := normalizeMsg(rec.Result.Msg)
		msg = "[" + rec.Result.Level + "] " + msg
		errorResults[msg] += count
	}

	t.Logf("errors by types: %d", len(errorResults))

	var messages []string
	for msg := range errorResults {
		messages = append(messages, msg)
	}
	sort.Strings(messages)
	for i, msg := range messages {
		fmt.Printf("#%3d. %6d : %s\n", i, errorResults[msg], msg)
	}

}

var thumbprintRegex = regexp.MustCompile(`(humbprint:\s*)([a-fA-F0-9]{40})`)
var thumbprint64Regex = regexp.MustCompile(`[a-f0-9]{64}`)
var logIdthumbprintegex = regexp.MustCompile(`\[([A-Z0-9]{6})\]`)
var correlationIdRegex = regexp.MustCompile(`([A-Z0-9]{16}_[0-9]{1,3})`)
var correlationId2Regex = regexp.MustCompile(` [A-Z0-9]{16} `)
var taskIndexRegex = regexp.MustCompile(`\(([0-9]{1,2})\.\) execution `)
var idSpRegex = regexp.MustCompile(`\(Id: ([0-9]{1,3}),`)
var idRegex = regexp.MustCompile(`\(Id:([0-9]{1,3}),`)
var ceridRegex = regexp.MustCompile(`CER_ID=([0-9]{3,7})`)
var ip74443Regex = regexp.MustCompile(`[0-9]{1,3}\.[0-9]{1,3}\.[0-9]{1,3}\.[0-9]{1,3}:7443`)
var ipWithPortRegex = regexp.MustCompile(`[0-9]{1,3}\.[0-9]{1,3}\.[0-9]{1,3}\.[0-9]{1,3}:[0-9]{3,5}`)
var bytesRegex = regexp.MustCompile(`[0-9]{1,6} bytes`)
var userCertHostUpperRegex = regexp.MustCompile(`USER.CERTIFICATES.CEE.[A-Z]*.NSROOT.NET`)
var userCertHostRegex = regexp.MustCompile(`user.certificates.cee.[a-z]*.nsroot.net`)
var mdaServerHostRegex = regexp.MustCompile(`mda-server:[a-z0-9]*.nam.nsroot.net`)
var mdaTaskEnumRegex = regexp.MustCompile(`[0-9]. : mrgapi.MdaTask`)
var mdaTaskReceivedRegex = regexp.MustCompile(`Received [0-9] tasks`)

func normalizeMsg(msg string) string {
	if strings.Contains(msg, "20210112085537-c389da54e794/notifyicon.go") {
		return msgNotifyIcon
	}
	if strings.Contains(msg, "EndUser not found by EU_UPN / EU_CUSTOM_ID_1 ") {
		return msgEndUserNotFound
	}
	if strings.Contains(msg, "Can't connect to wss://user.certificates") {
		return msgCanNotConnect
	}
	if strings.Contains(msg, "connection was forcibly closed by the remote host") {
		return msgConnForciblyClosed
	}
	if strings.Contains(msg, "connection was aborted by the software") {
		return msgConnAbortedByHost
	}
	if strings.Contains(msg, "Key pair with certificate imported: [MRG_") {
		return msgKpImported
	}
	if strings.Contains(msg, "Certificate imported: [MRG_") {
		return msgCertImported
	}
	if strings.Contains(msg, "Key pair deleted: [MRG_") {
		return msgKpDeleted
	}

	// Replace thumbprint values with zeros using regex
	msg = thumbprintRegex.ReplaceAllString(msg, "${1}?00000000000000000000000000000000000000?")
	msg = thumbprint64Regex.ReplaceAllString(msg, "000?000?000?000?000?000?000?000?000?000?000?000?000?000?000?000?")
	msg = logIdthumbprintegex.ReplaceAllString(msg, "[?00000]")
	msg = correlationIdRegex.ReplaceAllString(msg, "?000000000000000_?")
	msg = correlationId2Regex.ReplaceAllString(msg, "?00000000000000?")
	msg = taskIndexRegex.ReplaceAllString(msg, "(?.) execution ")
	msg = idSpRegex.ReplaceAllString(msg, "(Id: ?,")
	msg = idRegex.ReplaceAllString(msg, "(Id:?,")
	msg = ceridRegex.ReplaceAllString(msg, "CER_ID=?00000")
	msg = ip74443Regex.ReplaceAllString(msg, "10.?.?.?:7443")
	msg = ipWithPortRegex.ReplaceAllString(msg, "10.?.?.0:?000")
	msg = bytesRegex.ReplaceAllString(msg, "1?00 bytes")
	msg = mdaTaskEnumRegex.ReplaceAllString(msg, "?. : mrgapi.MdaTask")
	msg = mdaTaskReceivedRegex.ReplaceAllString(msg, "Received ? tasks")

	msg = userCertHostRegex.ReplaceAllString(msg, "user.certificates.cee.???.nsroot.net")
	msg = userCertHostUpperRegex.ReplaceAllString(msg, "USER.CERTIFICATES.CEE.???.NSROOT.NET")
	msg = mdaServerHostRegex.ReplaceAllString(msg, "mda-server:???.nam.nsroot.net")

	return msg
}

const msgNotifyIcon = `Shell_NotifyIcon\\n\\nStack:\\ngoroutine 205 [running]:\\nruntime/debug.Stack()\\n\\t/opt/gos/go1.22.6/src/runtime/debug/stack.go:24 +0x5e\\ngithub.com/lxn/walk.newErr(...)\\n\\t/home/lbukodi/go/pkg/mod/github.com/lxn/walk@v0.0.0-20210112085537-c389da54e794/error.go:81\\ngithub.com/lxn/walk.newError({0xd36be7, 0x10})\\n\\t/home/lbukodi/go/pkg/mod/github.com/lxn/walk@v0.0.0-20210112085537-c389da54e794/error.go:85 +0x25\\ngithub.com/lxn/walk.(*NotifyIcon).SetIcon(0xc0001da750, {0x12f8838, 0xc000072480})\\n\\t/home/lbukodi/go/pkg/mod/github.com/lxn/walk@v0.0.0-20210112085537-c389da54e794/notifyicon.go:323 +0x14b\\nmain.(*windowsUIFunctions).ChangeTaskbarIcon(0x166a710, {0xd3179d?, 0x0?}, {0xd49d43, 0x1d})\\n\\t/src/margareta_citi/mda/cmd/mdawin/main.go:202 +0x137\\nmargareta.noreg.hu/mda/pkg/mrgapi.SynchronizeWithServer({0x12fa850, 0x166a710})\\n\\t/src/margareta_citi/mda/pkg/mrgapi/syncproc.go:41 +0xf9\\ncreated by main.(*MDAMainWindow).WndProc in goroutine 1\\n\\t/src/margareta_citi/mda/cmd/mdawin/mdamainwnd.go:88 +0xa5\\n: 1`
const msgEndUserNotFound = `Sending error to core failed: server side error: EndUser not found by EU_UPN / EU_CUSTOM_ID_1 = ????@NAM.NSROOT.NET (Search in server log for code [000000] at ????)`
const msgCanNotConnect = `Sending error to core failed: Can't connect to wss://user.certificates.?.nsroot.net:7443/margareta-mda-srv/mda : write tcp 10.?.100.100:?000?->169.172.77.?:7443: i/o timeout: 1`
const msgConnForciblyClosed = `Sending error to core failed: write tcp 10.?.100.100:?000?->169.172.77.?:7443: wsasend: An existing connection was forcibly closed by the remote host.`
const msgConnAbortedByHost = `Sending error to core failed: write tcp 10.?.100.100:?000?->169.172.77.?:7443: wsasend: An established connection was aborted by the software in your host machine`
const msgKpImported = `Key pair with certificate imported: [MRG_?0000?] ?soeid? ?type? - ?cerid?`
const msgCertImported = `Certificate imported: [MRG_?0000?] ?soeid? ?type? - ?cerid?`
const msgKpDeleted = `Key pair deleted: [MRG_?0000?] ?soeid? ?type? - ?cerid?`

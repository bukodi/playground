package freejson

import (
	"fmt"
	"os"
	"sort"
	"strings"
	"testing"
)

/*

Example:
eg56209,"2025-06-19-10-55",H02PQLI1312P125
eg56209,"2025-06-19-10-55",H02PQLI1312P125
eg56209,"2025-06-19-10-55",H02PQLI1312P125
NC22161,"2025-06-19-10-47",H18NXLI1313P102
NC22161,"2025-06-19-10-47",H18NXLI1313P102
NC22161,"2025-06-19-10-47",H18NXLI1313P102
NC22161,"2025-06-19-10-47",H18NXLI1313P102
NC22161,"2025-06-19-10-47",H18NXLI1313P102
ju93120,"2025-06-19-10-37",H02KXLI1306P032
ju93120,"2025-06-19-10-37",H02KXLI1306P032
*/

type rec struct {
	soeid    string
	date     string
	hostname string
	isFSL0   bool
}

var bySoeID = map[string][]*rec{}
var byHostName = map[string][]*rec{}
var countBySoeIdAndHostName = map[[2]string]int{}
var fsl0Soeid = map[string]any{}
var nonFsl0Soeid = map[string]any{}

func TestParseCSV(t *testing.T) {
	data, err := os.ReadFile("/home/lbukodi/Downloads/nte_exists_hostname.csv")
	if err != nil {
		t.Error(err)
	}
	lines := strings.Split(string(data), "\n")
	fmt.Printf("\n------- process file --------\n")
	fmt.Printf("%d lines in the file\n", len(lines))
	for _, line := range lines {
		if line == "" {
			continue
		}
		r := parseLine(line)
		bySoeID[r.soeid] = append(bySoeID[r.soeid], r)
		byHostName[r.hostname] = append(byHostName[r.hostname], r)
		countBySoeIdAndHostName[[2]string{r.soeid, r.hostname}]++
		if r.isFSL0 {
			fsl0Soeid[r.soeid] = true
		} else {
			nonFsl0Soeid[r.soeid] = true
		}
	}
	fmt.Printf("Number of distinct SoeIds : %d\n", len(bySoeID))
	fmt.Printf("   - SoeIds with FSL0 suffix: %d\n", len(fsl0Soeid))
	fmt.Printf("   - SoeIds without FSL0 suffix: %d\n", len(nonFsl0Soeid))
	fmt.Printf("Number of distinct hostnames: %d\n", len(byHostName))
	fmt.Printf("Number of distinct soeid-hostname pairs: %d\n", len(countBySoeIdAndHostName))

	countHostsBySoeID := map[string]int{}
	countSoeIdByHosts := map[string]int{}
	for soeidAndHostname, _ := range countBySoeIdAndHostName {
		countHostsBySoeID[soeidAndHostname[0]]++
		countSoeIdByHosts[soeidAndHostname[1]]++
	}

	{
		x := map[int]int{}
		for _, numHosts := range countHostsBySoeID {
			x[numHosts]++
		}
		keys := []int{}
		for k := range x {
			keys = append(keys, k)
		}
		sort.Ints(keys)
		fmt.Printf("\n------ Number of hosts per soeid ----------\n")
		checkSumm := 0
		for _, k := range keys {
			explain := ""
			if k == 1 {
				explain = fmt.Sprintf("( There are %d hosts, what used by only one SoeId. )", x[k])
			} else if k == 2 {
				explain = fmt.Sprintf("( There are %d hosts, what used by exactly two SoeIds. )", x[k])
			} else if k == 3 {
				explain = fmt.Sprintf("( There are %d hosts, what used by exactly three SoeIds. )", x[k])
			}
			fmt.Printf("    %2d: %5d  %s\n", k, x[k], explain)
			checkSumm += k * x[k]
		}
		fmt.Printf(" (check: summ of distinct soeids-host pairs: %d)\n", checkSumm)
	}
	{
		x := map[int]int{}
		for _, numSoeids := range countSoeIdByHosts {
			x[numSoeids]++
		}
		keys := []int{}
		for k := range x {
			keys = append(keys, k)
		}
		sort.Ints(keys)
		fmt.Printf("\n--------- Number of soeids per host -----------\n")
		checkSumm := 0
		for _, k := range keys {
			explain := ""
			if k == 1 {
				explain = fmt.Sprintf("( There are %d SoeIds, what uses only one host. )", x[k])
			} else if k == 2 {
				explain = fmt.Sprintf("( There are %d SoeIds, what uses exactly two hosts. )", x[k])
			} else if k == 3 {
				explain = fmt.Sprintf("( There are %d SoeIds, what uses exactly three hosts. )", x[k])
			}
			fmt.Printf("    %2d: %5d  %s\n", k, x[k], explain)
			checkSumm += k * x[k]
		}
		fmt.Printf(" (check: summ of distinct soeids-host pairs: %d)\n", checkSumm)
	}

	onoToOnePairs := [][2]string{}
	for pair, _ := range countBySoeIdAndHostName {
		if countSoeIdByHosts[pair[1]] == 1 && countHostsBySoeID[pair[0]] == 1 {
			onoToOnePairs = append(onoToOnePairs, pair)
		}
	}
	fmt.Printf("----------- one-to-one pairs-------------)\n")
	fmt.Printf("Number of SoeId-Hostname pairs, where a SoeID uses only one host and that host is not used by other SoeId: %d\n", len(onoToOnePairs))

	hostNames := []string{}
	for pair, _ := range countBySoeIdAndHostName {
		if pair[1][0] != "H"[0] {
			hostNames = append(hostNames, pair[1])
		}
	}
	sort.Strings(hostNames)
	fmt.Printf("\n------- Hostnames not started with H:%d (%f %% of all hostnames) --------\n", len(hostNames), float32(len(hostNames)*100)/float32(len(byHostName)))
	for _, hostName := range hostNames {
		fmt.Printf("    %s\n", hostName)
	}
}

func parseLine(line string) *rec {
	r := rec{}
	// Split the line by commas
	parts := strings.Split(line, ",")

	// Remove quotes from the date field
	if len(parts) != 3 {
		panic("Invalid line: " + line)
	}
	// First part is soeid
	r.soeid = parts[0]

	// Second part is date (remove quotes)
	date := parts[1]
	date = strings.Trim(date, "\"")
	r.date = date

	// Third part is hostname
	r.hostname = parts[2]

	// Normalize line
	{
		r.hostname = strings.ToUpper(r.hostname)
		r.soeid = strings.ToUpper(r.soeid)
		r.soeid = strings.ReplaceAll(r.soeid, `"`, "")
		if strings.HasSuffix(r.soeid, ".FSL0") {
			r.soeid = strings.TrimSuffix(r.soeid, ".FSL0")
			r.isFSL0 = true
		}
	}

	return &r
}

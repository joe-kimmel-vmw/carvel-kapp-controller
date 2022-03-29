package main

import (
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

func main() {
	numPkgs := []int{20, 40, 100, 200, 400, 600, 1000, 1600}
	deployTimes := []string{}
	deleteTimes := []string{}
	for _, numCnoops := range numPkgs {
		time.Sleep(1 * time.Second) // I have a vague feeling like part of the problem is we just get ratelimited doing this test too fast.
		totalCnoops := numCnoops
		fmt.Printf("\n===========\n\t Starting for %d noops\n===========\n", totalCnoops)
		fname := writePkgr(numCnoops)
		// defer os.Remove(fname)

		t1 := time.Now()
		cmd := exec.Command("kapp", "deploy", "-f", fname, "-a", fname[:len(fname)-5], "-y", "--wait-resource-timeout=0s")
		stdout, err := cmd.Output()
		fmt.Println(string(stdout))
		if err != nil {
			fmt.Println(err.Error())
			kctlo, kctlerr := exec.Command("kubectl", "get", "pkgrs", "-A", "-o", "yaml").Output()
			fmt.Println("kubectl get pkgrs: \n", string(kctlo), kctlerr.Error())
			panic(err)
		}
		t2 := time.Now()
		deployTime := t2.Sub(t1).Seconds()

		cmd = exec.Command("kapp", "delete", "-a", fname[:len(fname)-5], "-y")
		stdout, err = cmd.Output()
		fmt.Println(string(stdout))
		if err != nil {
			fmt.Println(err.Error())
			panic(err)
		}
		t3 := time.Now()
		deleteTime := t3.Sub(t2).Seconds()
		fmt.Printf("\n===========\n\t Finished %d Cnoops in %f seconds (%f deploy ; %f delete)\n===========\n", totalCnoops, t3.Sub(t1).Seconds(), deployTime, deleteTime)
		deployTimes = append(deployTimes, strconv.FormatFloat(deployTime, 'f', 1, 64))
		deleteTimes = append(deleteTimes, strconv.FormatFloat(deleteTime, 'f', 1, 64))
	}

	fname := fmt.Sprintf("results-%v.csv", time.Now().Unix())
	f, err := os.Create(fname)
	if err != nil {
		panic(err)
	}
	defer f.Close()
	f.WriteString(strings.Join(toStringArr(numPkgs), ", ")) // TODO: it's numPkgs * numVersions ...
	f.WriteString("\n")
	f.WriteString(strings.Join(deployTimes, ", "))
	f.WriteString("\n")
	f.WriteString(strings.Join(deleteTimes, ", "))
	f.WriteString("\n")

}

// sure do miss constructs like [str(x) for x in inp]
func toStringArr(inp []int) []string {
	icantbelieveihavetodothis := []string{}
	for _, russcoxissupersmarttho := range inp {
		icantbelieveihavetodothis = append(icantbelieveihavetodothis, strconv.Itoa(russcoxissupersmarttho))
	}
	return icantbelieveihavetodothis
}

func writePkgr(numCnoops int) string {
	totalCnoops := numCnoops

	preamble := fmt.Sprintf(`
apiVersion: packaging.carvel.dev/v1alpha1
kind: PackageRepository
metadata:
  name: minimal-repo-%d.tanzu.carvel.dev
  namespace: default
  annotations:
    kapp.k14s.io/disable-original: ""
spec:
  fetch:
    inline:
      paths:
`, totalCnoops)

	pkgStr := `
        packages/pkg.test.carvel.dev/noop%[1]d.test.carvel.dev.0.0.1.yml: |
          ---
          kind: CarvelNoop
          apiVersion: data.packaging.carvel.dev/v1alpha1
          metadata:
            name: "foo"
            namespace: "default"
`
	fname := fmt.Sprintf("pkgr-%d.yaml", totalCnoops)
	f, err := os.Create(fname)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	f.WriteString(preamble)
	for i := 0; i < numCnoops; i++ {
		_, err := f.WriteString(fmt.Sprintf(pkgStr, i))
		if err != nil {
			panic(err)
		}
	}
	return fname
}

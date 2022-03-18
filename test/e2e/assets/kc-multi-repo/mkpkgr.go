package main

import (
	"flag"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

func main() {
	numVersions := flag.Int("numVers", 2, "The number of minor versions for each Package (total packages will be numPkgs x numVers)")
	flag.Parse()

	numPkgs := []int{10, 20, 50, 100, 200, 300, 500, 800}
	deployTimes := []string{}
	deleteTimes := []string{}
	for _, numPackages := range numPkgs {
		totalPackages := numPackages * *numVersions
		fmt.Printf("\n===========\n\t Starting for %d Packages\n===========\n", totalPackages)
		fname := writePkgr(numPackages, *numVersions)
		// defer os.Remove(fname)

		t1 := time.Now()
		cmd := exec.Command("kapp", "deploy", "-f", fname, "-a", fname[:len(fname)-5], "-y", "--wait-resource-timeout=0s", "--wait-timeout=0s")
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
		fmt.Printf("\n===========\n\t Finished %d Packages in %f seconds (%f deploy ; %f delete)\n===========\n", totalPackages, t3.Sub(t1).Seconds(), deployTime, deleteTime)
		deployTimes = append(deployTimes, strconv.FormatFloat(deployTime, 'f', 1, 64))
		deleteTimes = append(deleteTimes, strconv.FormatFloat(deleteTime, 'f', 1, 64))
	}

	fname := fmt.Sprintf("results-%v.csv", time.Now().Unix())
	f, err := os.Create(fname)
	if err != nil {
		panic(err)
	}
	defer f.Close()
	f.WriteString(strings.Join(toStringArr(numPkgs), ", "))
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

func writePkgr(numPackages int, numVersions int) string {
	totalPackages := numPackages * numVersions

	preamble := fmt.Sprintf(`
apiVersion: packaging.carvel.dev/v1alpha1
kind: PackageRepository
metadata:
  name: minimal-repo-%d.tanzu.carvel.dev
  # Adds it to global namespace (as defined by kapp-controller)
  # which makes packages available in all namespaces
  namespace: kapp-controller-packaging-global
  annotations:
    kapp.k14s.io/disable-original: ""
spec:
  fetch:
    inline:
      paths:
`, totalPackages)

	pkgStr := `
        packages/pkg.test.carvel.dev/pkg%[1]d.test.carvel.dev.0.%[2]d.0.yml: |
          ---
          apiVersion: data.packaging.carvel.dev/v1alpha1
          kind: Package
          metadata:
            name: pkg%[1]d.test.carvel.dev.0.%[2]d.0
          spec:
            refName: pkg%[1]d.test.carvel.dev
            version: 0.%[2]d.0
            template:
              spec: {}
`
	fname := fmt.Sprintf("pkgr-%d.yaml", totalPackages)
	f, err := os.Create(fname)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	f.WriteString(preamble)
	for i := 0; i < numPackages; i++ {
		for j := 0; j < numVersions; j++ {
			_, err := f.WriteString(fmt.Sprintf(pkgStr, i, j))
			if err != nil {
				panic(err)
			}
		}
	}
	return fname
}

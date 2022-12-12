package main

import (
	"fmt"
	"github.com/diamondburned/gotk4/pkg/gio/v2"
	"github.com/diamondburned/gotk4/pkg/gtk/v4"
	"github.com/prometheus/procfs"
	"github.com/prometheus/procfs/blockdevice"
	"os"
	"strconv"
	"strings"
	"time"
)

func main() {
	app := gtk.NewApplication("com.github.diamondburned.gotk4-examples.gtk4.simple", gio.ApplicationFlagsNone)
	app.ConnectActivate(func() { activate(app) })

	if code := app.Run(os.Args); code > 0 {
		os.Exit(code)
	}
}

func activate(app *gtk.Application) {
	window := gtk.NewApplicationWindow(app)
	window.SetTitle("gotk4 Example")

	box := gtk.NewBox(gtk.OrientationVertical, 0)
	window.SetChild(box)

	//sysuid := syscall.Getuid()
	fs, _ := procfs.NewDefaultFS()
	cpuInfo, _ := fs.CPUInfo()
	fsd, _ := blockdevice.NewFS("/proc", "/sys")

	stats, _ := fsd.ProcDiskstats()

	label1 := gtk.NewLabel("cpu model: " + cpuInfo[0].ModelName)
	label2 := gtk.NewLabel("cpu cores: " + strconv.Itoa(int(cpuInfo[0].CPUCores)))
	label3 := gtk.NewLabel("read IO's: " + strconv.Itoa(int(stats[0].IOStats.ReadIOs)))
	label4 := gtk.NewLabel("write IO's: " + strconv.Itoa(int(stats[0].IOStats.WriteIOs)))

	box.Append(label1)
	box.Append(label2)
	box.Append(label3)
	box.Append(label4)

	window.SetDefaultSize(400, 300)
	window.Show()
	go updateCpuLoad(label1)
	go updateRW(label3, label4)
}

func updateRW(labelR *gtk.Label, labelW *gtk.Label) {
	for true {
		fsd, _ := blockdevice.NewFS("/proc", "/sys")
		stats, _ := fsd.ProcDiskstats()
		var read uint64
		var write uint64
		read = 0
		write = 0
		for _, stat := range stats {
			read += stat.ReadIOs
			write += stat.WriteIOs
		}
		labelR.SetText("read IO's: " + strconv.Itoa(int(read)))
		labelW.SetText("write IO's: " + strconv.Itoa(int(write)))
		time.Sleep(time.Second)
	}
}

func updateCpuLoad(label *gtk.Label) {
	prevTotal := 0
	prevWork := 0

	for true {
		currentTotal := 0
		currentWork := 0
		data, _ := os.ReadFile("/proc/stat")
		str := string(data)
		cpuStr := strings.Split(str, "\n")[0]
		cpuStrSplt := strings.Split(cpuStr, " ")
		for i := 1; i < len(cpuStrSplt); i++ {
			jiff, _ := strconv.Atoi(cpuStrSplt[i])
			currentTotal += jiff
			if i < 4 {
				currentWork += jiff
			}
		}
		workOverPeriod := float64(currentWork - prevWork)
		totalOverPeriod := float64(currentTotal - prevTotal)
		println(workOverPeriod, totalOverPeriod, workOverPeriod/totalOverPeriod)
		cpuLoad := workOverPeriod / totalOverPeriod * 100
		prevTotal = currentTotal
		prevWork = currentWork
		label.SetText("cpu load: " + fmt.Sprintf("%f", cpuLoad))
		time.Sleep(time.Second)
	}
}

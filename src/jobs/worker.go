package jobs

import (
	"os"
	"os/exec"
	"strings"
	"time"

	"path/filepath"

	"github.com/lfkeitel/inca3/src/models"
	"github.com/lfkeitel/inca3/src/utils"
)

type worker struct {
	e         *utils.Environment
	job       *models.Job
	startTime time.Time
	cancel    chan bool
	done      chan bool
	errors    chan error
}

func newWorker(e *utils.Environment) *worker {
	return &worker{
		e:      e,
		cancel: make(chan bool, 1),
		done:   make(chan bool, 1),
		errors: make(chan error, 1),
	}
}

func (w *worker) start() error {
	w.job.Status = models.Starting

	var devices []*models.Device
	var err error

	if len(w.job.Devices) == 0 {
		devices, err = models.GetAllDevices(w.e)
	} else {
		devices, err = models.GetDevicesByIDs(w.e, w.job.Devices)
	}

	if err != nil {
		w.job.Status = models.Stopped
		return err
	}

	w.job.Status = models.Running
	go func(d []*models.Device) {
		w.run(d)
	}(devices)

	return nil
}

func (w *worker) run(devices []*models.Device) {
	wg := NewLimitGroup(w.e.Config.Job.MaxConnections)
	date := time.Now().Format("2006-01-02T15:03:04")
	for _, d := range devices {
		// Check for cancel signal; if given, break
		select {
		case <-w.cancel:
			w.job.Status = models.Stopping
			break
		default:
		}

		configFileDir := filepath.Join(w.e.Config.DirPaths.BaseDir, d.Address)
		if err := os.MkdirAll(configFileDir, 0755); err != nil {
			w.e.Log.WithField("Path", configFileDir).Error("Failed to make directories")
			continue
		}

		configFile := filepath.Join(configFileDir, date+".conf")
		args := w.getArguments(d.Type.Args, d.Address, configFile)

		wg.Add(1)
		go func(de *models.Device, cFile string, argList []string) {
			defer wg.Done()
			// Run job script
			err := w.execScript(filepath.Join(w.e.Config.DirPaths.ScriptDir, de.Type.Script), argList)
			if err != nil {
				w.e.Log.WithField("Err", err).Error("Failed to get config")
				return
			}

			// Build a configuration entry
			c := models.NewConfig(w.e)
			c.Slug = de.Slug + "_" + date
			c.DeviceID = de.ID
			c.Filename = filepath.Join(de.Address, date+".conf")
			c.Created = time.Now()
			c.Compressed = false

			if err := c.Save(); err != nil {
				w.e.Log.WithField("Err", err).Error("Failed to save config")
			}
		}(d, configFile, args)

		wg.Wait()
	}

	wg.WaitAll()
	w.job.Status = models.Finished
}

func (w *worker) getArguments(argStr string, host string, filename string) []string {
	args := strings.Split(argStr, ";")
	argList := make([]string, len(args))
	for i, a := range args {
		switch a {
		case "$address":
			argList[i] = host
			break
		case "$username":
			argList[i] = w.e.Config.Job.RemoteUsername
			break
		case "$password":
			argList[i] = w.e.Config.Job.RemotePassword
			break
		case "$logfile":
			argList[i] = filename
			break
		case "$enablepw":
			argList[i] = w.e.Config.Job.EnablePassword
			break
		}
	}
	return argList
}

func (w *worker) execScript(sfn string, args []string) error {
	_, err := exec.Command(sfn, args...).Output()
	if err != nil {
		w.e.Log.WithField("Err", err).Error()
		return err
	}
	//stdOutLogger.Info(string(out))
	return nil
}

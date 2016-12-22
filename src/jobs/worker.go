package jobs

import (
	"os"
	"os/exec"
	"time"

	"path/filepath"

	"github.com/lfkeitel/inca3/src/models"
	"github.com/lfkeitel/inca3/src/utils"
	"github.com/lfkeitel/verbose"
)

type worker struct {
	e         *utils.Environment
	job       *job
	startTime time.Time
	cancel    chan bool
	done      chan bool
	errors    chan error
	running   bool
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

	w.e.Log.Debug("Worker: Getting devices")
	if len(w.job.Devices) == 0 {
		devices, err = models.GetAllDevices(w.e)
	} else {
		devices, err = models.GetDevicesByIDs(w.e, w.job.Devices)
	}

	if err != nil {
		w.job.Status = models.Stopped
		w.done <- true
		return err
	}

	w.job.Total = len(devices)

	w.job.Status = models.Running
	w.startTime = time.Now()
	w.job.Start = w.startTime
	go func(d []*models.Device) {
		w.e.Log.Debug("Worker: Running job")
		w.run(d)
		w.e.Log.Debug("Worker: Job finished")
		w.done <- true
	}(devices)

	return nil
}

func (w *worker) run(devices []*models.Device) {
	wg := NewLimitGroup(w.e.Config.Job.MaxConnections)
	date := time.Now().Format("2006-01-02T15:04:05")

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
			w.job.finished++
			continue
		}

		configFile := filepath.Join(configFileDir, date+".conf")
		args := w.getArguments(d.Address, configFile)

		wg.Add(1)
		go func(de *models.Device, cFile string, argList []string) {
			w.e.Log.WithFields(verbose.Fields{
				"Address": de.Address,
				"File":    cFile,
				"Script":  de.Type.Script,
			}).Debug("Worker: Running script")
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

			w.e.Log.WithField("Address", de.Address).Debug("Script finished, config saved")
			w.job.finished++
		}(d, configFile, args)

		wg.Wait()
	}

	wg.WaitAll()
	w.job.Status = models.Finished
	w.job.End = time.Now()
	if err := w.job.Save(); err != nil {
		select {
		case w.errors <- err:
		default:
		}
		w.e.Log.WithField("Err", err).Error("Failed to save job")
	}
}

func (w *worker) getArguments(host string, filename string) []string {
	return []string{
		host,
		w.e.Config.Job.RemoteUsername,
		w.e.Config.Job.RemotePassword,
		filename,
		w.e.Config.Job.EnablePassword,
	}
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

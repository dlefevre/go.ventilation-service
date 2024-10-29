package controller

import (
	"sync"
	"time"

	"github.com/dlefevre/go.ventilation-service/config"
	"github.com/dlefevre/go.ventilation-service/gpio"
	"github.com/rs/zerolog/log"
)

// Enum pseudo-type.
type Enum uint

// Queue size for the command channel.
const queueSize = 3

// Enumeration of commands.
const (
	CmdDummy   Enum = iota // CmdDummy does nothing, but prevents errors when closing the channel.
	CmdSpeed1              // CmdSpeed1 identifies the speed 1 request command
	CmdSpeed2              // CmdSpeed2 identifies the speed 2 request command
	CmdSpeed3              // CmdSpeed3 identifies the speed 3 request command
	CmdAway                // CmdAway identifies the away request command
	CmdAuto                // CmdAuto identifies the auto request command
	CmdTimer15             // CmdTimer15 identifies the timer request command (10')
	CmdTimer30             // CmdTimer30 identifies the timer request command (30')
	CmdTimer60             // CmdTimer60 identifies the timer request command (60')
)

const (
	reactTime time.Duration = 100 * time.Millisecond
)

var (
	instance *VentilationControllerService
	once     sync.Once
)

// VentilationControllerService implements the service for controlling the ventilation and reporting its state.
type VentilationControllerService struct {
	command chan Enum
	adapter gpio.GPIOAdapter
	wg      sync.WaitGroup
	lock    sync.RWMutex
	backoff time.Duration
}

// GetVentilationControllerService returns the one and only VentilationControllerServiceImpl instance.
func GetVentilationControllerService() *VentilationControllerService {
	once.Do(func() {
		instance = newVentilationControllerService()
	})
	return instance
}

// Creates a new VentilationControllerServiceImpl object.
func newVentilationControllerService() *VentilationControllerService {
	return &VentilationControllerService{
		command: nil,
		adapter: gpio.GetGPIOAdapter(),
		wg:      sync.WaitGroup{},
		backoff: time.Duration(config.GetGPIOBackoff()) * time.Millisecond,
	}
}

// Main loop for handling commands.
func (d *VentilationControllerService) commandLoop() {
	defer d.wg.Done()

	for {
		cmd, ok := <-d.command
		if !ok {
			log.Info().Msg("command channel closed")
			break
		}
		switch cmd {
		case CmdSpeed1:
			d.toggle(d.adapter.WriteSpeed1Pin)
		case CmdSpeed2:
			d.toggle(d.adapter.WriteSpeed2Pin)
		case CmdSpeed3:
			d.toggle(d.adapter.WriteSpeed3Pin)
		case CmdAway:
			d.toggle(d.adapter.WriteAwayPin)
		case CmdAuto:
			d.toggle(d.adapter.WriteAutoPin)
		case CmdTimer15:
			d.toggleX(d.adapter.WriteTimerPin, 1)
		case CmdTimer30:
			d.toggleX(d.adapter.WriteTimerPin, 2)
		case CmdTimer60:
			d.toggleX(d.adapter.WriteTimerPin, 3)
		default:
			log.Warn().Msgf("unknown command: %v", d.command)
		}
	}

	log.Info().Msg("commandLoop exiting")
}

func (d *VentilationControllerService) toggle(pinHandle func(bool)) {
	pinHandle(true)
	time.Sleep(reactTime)
	pinHandle(false)
	time.Sleep(d.backoff)
}

func (d *VentilationControllerService) toggleX(pinHandle func(bool), times int) {
	for i := 0; i < times; i++ {
		pinHandle(true)
		time.Sleep(reactTime)
		pinHandle(false)
		time.Sleep(reactTime)
	}
	time.Sleep(d.backoff - reactTime)
}

// Start all goroutines.
func (d *VentilationControllerService) Start() {
	d.lock.Lock()
	defer d.lock.Unlock()

	d.command = make(chan Enum, queueSize)
	go d.commandLoop()
	d.wg.Add(1)
}

// Stop all goroutines, gracefully.
func (d *VentilationControllerService) Stop() {
	d.lock.Lock()
	close(d.command)
	d.lock.Unlock()
	log.Info().Msg("Stopping VentilationControllerService")

	d.wg.Wait()
	log.Info().Msg("VentilationControllerService stopped")
}

func (d *VentilationControllerService) SendCommand(command Enum) {
	select {
	case d.command <- command:
	default:
		<-d.command
		d.command <- command
	}
}

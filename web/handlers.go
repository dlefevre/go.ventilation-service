package web

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/dlefevre/go.ventilation-service/controller"
	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog/log"
)

// SimpleResponse is a simple response object, containing a result (ok, nok).
type SimpleResponse struct {
	Result string `json:"result"`
}

// ErrorResponse is a response object for errors, containing  a result (nok) and a message.
type ErrorResponse struct {
	SimpleResponse
	Message string `json:"message"`
}

// SpeedMessage is a message object for speed commands.
type SpeedMessage struct {
	Speed string `json:"speed"`
}

// TimerMessage is a message object for timer commands (should be 15, 30 or 30 minutes).
type TimerMessage struct {
	Duration int `json:"duration"`
}

// Health check handler.
func healthCheck(c echo.Context) error {
	return c.JSON(http.StatusOK, SimpleResponse{
		Result: "OK",
	})
}

func bodyParser(c echo.Context, dest interface{}) error {
	body, err := io.ReadAll(c.Request().Body)
	if err != nil {
		return fmt.Errorf("error reading request body: %v", err)
	}
	if err := json.Unmarshal(body, dest); err != nil {
		return fmt.Errorf("error parsing request body: %v", err)
	}
	return nil
}

// Handler for speed command
func speedHandler(c echo.Context) error {
	dc := controller.GetVentilationControllerService()

	var speed SpeedMessage
	if err := bodyParser(c, &speed); err != nil {
		log.Error().Msgf("Error parsing request: %v", err)
		return c.JSON(http.StatusBadRequest, ErrorResponse{
			SimpleResponse: SimpleResponse{Result: "nok"},
			Message:        fmt.Sprintf("Error parsing request: %v", err),
		})
	}

	switch speed.Speed {
	case "1", "low":
		dc.SendCommand(controller.CmdSpeed1)
	case "2", "medium":
		dc.SendCommand(controller.CmdSpeed2)
	case "3", "high":
		dc.SendCommand(controller.CmdSpeed3)
	default:
		log.Error().Msgf("Unknown speed: %s", speed.Speed)
		return c.JSON(http.StatusBadRequest, ErrorResponse{
			SimpleResponse: SimpleResponse{Result: "nok"},
			Message:        fmt.Sprintf("Invalid speed: %s", speed.Speed),
		})
	}

	return c.JSON(http.StatusOK, SimpleResponse{
		Result: "ok",
	})
}

// Handler for timer command
func timerHandler(c echo.Context) error {
	dc := controller.GetVentilationControllerService()

	var timer TimerMessage
	if err := bodyParser(c, &timer); err != nil {
		log.Error().Msgf("Error parsing request: %v", err)
		return c.JSON(http.StatusBadRequest, ErrorResponse{
			SimpleResponse: SimpleResponse{Result: "nok"},
			Message:        fmt.Sprintf("Error parsing request: %v", err),
		})
	}

	switch timer.Duration {
	case 15:
		dc.SendCommand(controller.CmdTimer15)
	case 30:
		dc.SendCommand(controller.CmdTimer30)
	case 60:
		dc.SendCommand(controller.CmdTimer60)
	default:
		log.Error().Msgf("Invalid duration: %d", timer.Duration)
		return c.JSON(http.StatusBadRequest, ErrorResponse{
			SimpleResponse: SimpleResponse{Result: "nok"},
			Message:        fmt.Sprintf("Invalid duration: %d", timer.Duration),
		})
	}

	return c.JSON(http.StatusOK, SimpleResponse{
		Result: "ok",
	})
}

// Handler for away command
func awayHandler(c echo.Context) error {
	dc := controller.GetVentilationControllerService()
	dc.SendCommand(controller.CmdAway)
	return c.JSON(http.StatusOK, SimpleResponse{
		Result: "ok",
	})
}

// Handler for auto command
func autoHandler(c echo.Context) error {
	dc := controller.GetVentilationControllerService()
	dc.SendCommand(controller.CmdAuto)
	return c.JSON(http.StatusOK, SimpleResponse{
		Result: "ok",
	})
}

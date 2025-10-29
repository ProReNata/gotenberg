package minibytes

import (
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/gotenberg/gotenberg/v8/pkg/modules/api"
	"github.com/labstack/echo/v4"
)

func minibytesRoute() api.Route {
	return api.Route{
		Method:      http.MethodPost,
		Path:        "/forms/minibytes",
		IsMultipart: true,
		Handler: func(c echo.Context) error {
			ctx := c.Get("context").(*api.Context)

			form := ctx.FormData()
			inputPaths := []string{}
			err := form.MandatoryPaths([]string{".rtf", ".pdf"}, &inputPaths).
				Validate()
			if err != nil {
				return fmt.Errorf("validate form data: %w", err)
			}

			// TODO: Will the framework outomatically handle multiple files?
			var outputPaths []string
			for _, inputPath := range inputPaths {
				// open uploaded file (inputPath is a path on disk)
				f, err := os.Open(inputPath)
				if err != nil {
					return err
				}

				buf := make([]byte, 100)
				n, err := f.Read(buf)
				// close ASAP
				_ = f.Close()
				if err != nil && err != io.EOF {
					return err
				}

				// create an output path and write the bytes
				outputPath := ctx.GeneratePath(".bin")
				if err := os.WriteFile(outputPath, buf[:n], 0644); err != nil {
					return err
				}

				outputPaths = append(outputPaths, outputPath)
			}

			// register outputs so Gotenberg will return them to the client
			return ctx.AddOutputPaths(outputPaths...)
		},
	}
}

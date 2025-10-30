package minibytes

import (
	"context"
	"fmt"
	"net/http"
	"os/exec"

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
				// generate a temporary output PNG path
				outputPath := ctx.GeneratePath(".png")

				// run Ghostscript to render first page as PNG
				cmd := exec.CommandContext(context.Background(),
					"gs",
					"-q",
					"-dNOPAUSE",
					"-dBATCH",
					"-dSAFER",
					"-sDEVICE=pngalpha",
					"-dFirstPage=1",
					"-dLastPage=1",
					"-r150", // resolution in DPI
					"-sOutputFile="+outputPath,
					inputPath,
				)

				if err := cmd.Run(); err != nil {
					return fmt.Errorf("failed to generate preview for %s: %w", inputPath, err)
				}

				outputPaths = append(outputPaths, outputPath)
			}

			// register outputs so Gotenberg will return them to the client
			return ctx.AddOutputPaths(outputPaths...)
		},
	}
}

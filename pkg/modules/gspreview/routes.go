package minibytes

import (
	"context"
	"fmt"
	"net/http"
	"os/exec"

	"github.com/gotenberg/gotenberg/v8/pkg/modules/api"
	// Where is this used???
	"github.com/labstack/echo/v4"
)

func gspreviewRoute() api.Route {
	return api.Route{
		Method:      http.MethodPost,
		Path:        "/forms/gspreview",
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
			xsize := 0
			err = form.Int("xsize", &xsize, -1).Validate()
			if err != nil {
				return fmt.Errorf("parse xsize: %w", err)
			}
			sizeArgument := "-r150"
			if xsize > 0 {
				sizeArgument = fmt.Sprintf("-r%d", xsize)
			}

			// TODO: Will the framework outomatically handle multiple files?
			var outputPaths []string
			for _, inputPath := range inputPaths {
				// generate a temporary output PNG path
				outputPath := ctx.GeneratePath(".png")

				// run Ghostscript to render first page as PNG
				// TODO: review options:
				//  Ensure non-transparent?
				//  Size in pixels? (not really possible in pure ghostscript)
				//  ???
				cmd := exec.CommandContext(context.Background(),
					"gs",
					"-q",
					"-dNOPAUSE",
					"-dBATCH",
					"-dSAFER",
					"-sDEVICE=pngalpha",
					"-dFirstPage=1",
					"-dLastPage=1",
					sizeArgument, // "-r150", // Prorenata: gs only support output size in as DPI... Maybe use magickimage?
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

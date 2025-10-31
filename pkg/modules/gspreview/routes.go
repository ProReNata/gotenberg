package minibytes

import (
	"context"
	"fmt"
	"net/http"
	"os/exec"

	"github.com/gotenberg/gotenberg/v8/pkg/modules/api"
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
			err = form.Int("xsize", &xsize, 0).Validate()
			if err != nil {
				return fmt.Errorf("parse xsize: %w", err)
			}
			sizeArgument := "1200x"
			if xsize > 0 {
				sizeArgument = fmt.Sprintf("%dx", xsize)
			}

			var outputPaths []string
			for _, inputPath := range inputPaths {
				outputPath := ctx.GeneratePath(".png") // Tmp output path
				var cmd *exec.Cmd
				if xsize > 0 {
					cmd = exec.CommandContext(context.Background(),
						"gm", "convert", "-adjoin",
						"-define", "pdf:use-cropbox=true",
						"-density", "150",
						"-resize", sizeArgument,
						"-quality", "100",
						fmt.Sprintf("%s[0]", inputPath), outputPath,
					)
				} else {
					cmd = exec.CommandContext(context.Background(),
						"gm", "convert", fmt.Sprintf("%s[0]", inputPath), outputPath,
					)
				}

				if err := cmd.Run(); err != nil {
					return fmt.Errorf("failed to generate preview for %s: %w", inputPath, err)
				}

				outputPaths = append(outputPaths, outputPath)
			}

			return ctx.AddOutputPaths(outputPaths...)
		},
	}
}

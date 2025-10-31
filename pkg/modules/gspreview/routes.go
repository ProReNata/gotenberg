package gspreview

import (
	"context"
	"fmt"
	"net/http"
	"os/exec"
	"strings"

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
			outputFormat := ""
			xsize := 0
			err := form.AnyMandatoryPaths(&inputPaths).Int("xsize", &xsize, 0).String("outputFormat", &outputFormat, "auto").Validate()
			if err != nil {
				return fmt.Errorf("validate form data: %w", err)
			}
			if outputFormat != "auto" && outputFormat != "png" && outputFormat != "pdf" {
				formatErrorMsg := "outputFormat must be one of [auto, png, pdf]"
				return api.WrapError(fmt.Errorf("%s", formatErrorMsg), api.NewSentinelHttpError(http.StatusBadRequest, formatErrorMsg))
			}
			sizeArgument := "1200x"
			if xsize > 0 {
				sizeArgument = fmt.Sprintf("%dx", xsize)
			}

			var outputPaths []string
			for _, inputPath := range inputPaths {
				toPng := (outputFormat == "auto" && strings.HasSuffix(strings.ToLower(inputPath), ".pdf")) || outputFormat == "png"
				var outputPath string
				var cmd *exec.Cmd
				if toPng {
					// "gm" parameters copied from Eketorp 3.80.0
					outputPath = ctx.GeneratePath(".png")
					cmd = exec.CommandContext(context.Background(),
						"gm", "convert", "-adjoin",
						"-define", "pdf:use-cropbox=true",
						"-density", "150",
						"-resize", sizeArgument,
						"-quality", "100",
						fmt.Sprintf("%s[0]", inputPath), outputPath,
					)
				} else {
					outputPath = ctx.GeneratePath(".pdf")
					cmd = exec.CommandContext(context.Background(),
						"gm", "convert", inputPath, outputPath,
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

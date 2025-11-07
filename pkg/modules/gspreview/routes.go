package gspreview

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/gotenberg/gotenberg/v8/pkg/gotenberg"
	"github.com/gotenberg/gotenberg/v8/pkg/modules/api"
	"github.com/labstack/echo/v4"
)

/*
If adding more modules, consider taking a look at this example project
that extends the PDF engine protocol:
https://github.com/Vrex123/gotenberg-ghostscript/tree/main

TODO: using a go library that links directly to graphicsmagic libraries
might give a moderate performance boost.
Maybe 40% based on testing under python.
*/
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
				return api.WrapError(
					fmt.Errorf("%s", formatErrorMsg),
					api.NewSentinelHttpError(http.StatusBadRequest, formatErrorMsg),
				)
			}
			sizeArgument := "1200x"
			if xsize > 0 {
				sizeArgument = fmt.Sprintf("%dx", xsize)
			}

			var outputPaths []string
			for _, inputPath := range inputPaths {
				toPng := (outputFormat == "auto" && strings.HasSuffix(strings.ToLower(inputPath), ".pdf")) || outputFormat == "png"
				var outputPath string
				var cmd *gotenberg.Cmd
				if toPng {
					// "gm" parameters copied from Eketorp 3.80.0
					outputPath = ctx.GeneratePath(".png")
					args := []string{
						"convert", "-adjoin",
						"-define", "pdf:use-cropbox=true",
						"-density", "150",
						"-resize", sizeArgument,
						"-quality", "100",
						fmt.Sprintf("%s[0]", inputPath), outputPath,
					}
					cmd, err = gotenberg.CommandContext(ctx, ctx.Log(), "gm", args...)
				} else {
					outputPath = ctx.GeneratePath(".pdf")
					args := []string{
						"convert",
						inputPath,
						outputPath,
					}
					cmd, err = gotenberg.CommandContext(ctx, ctx.Log(), "gm", args...)
				}
				if err != nil {
					return fmt.Errorf("create command: %w", err)
				}
				_, err = cmd.Exec()
				if err != nil {
					return fmt.Errorf("failed to convert for %s: %w", inputPath, err)
				}

				outputPaths = append(outputPaths, outputPath)
			}

			return ctx.AddOutputPaths(outputPaths...)
		},
	}
}

package gspreview

import (
	"bytes"
	"fmt"
	"image"
	"image/png"
	"net/http"
	"os"
	"strings"

	"github.com/gotenberg/gotenberg/v8/pkg/modules/api"
	"github.com/klippa-app/go-pdfium/requests"
	"github.com/labstack/echo/v4"

	_ "image/jpeg"
	_ "image/png"

	_ "golang.org/x/image/bmp"
	_ "golang.org/x/image/tiff"
	_ "golang.org/x/image/webp"

	"github.com/jung-kurt/gofpdf"
)

/*
PDF to PNG conversion via pdfium
*/
func RenderPdfFileToPng(m *Module, inputPath string, outputPath string, width int) error {
	// Mutext because c-bindings are not thread-safe
	m.pdfiumMutex.Lock()
	defer m.pdfiumMutex.Unlock()

	doc, err := m.pdfiumInstance.OpenDocument(&requests.OpenDocument{
		FilePath: &inputPath,
	})
	if err != nil {
		return fmt.Errorf("Error opening document: %w", err)
	}
	defer m.pdfiumInstance.FPDF_CloseDocument(&requests.FPDF_CloseDocument{
		Document: doc.Document,
	})

	render, err := m.pdfiumInstance.RenderPageInPixels(&requests.RenderPageInPixels{
		Page: requests.Page{
			ByIndex: &requests.PageByIndex{
				Document: doc.Document,
				Index:    0,
			},
		},
		Width:      width,
		Height:     0,
		RenderForm: true,
	})
	if err != nil {
		return fmt.Errorf("Render fail: %w", err)
	}

	outfile, err := os.Create(outputPath)
	if err != nil {
		return fmt.Errorf("Faile to create file: %w", err)
	}
	defer outfile.Close()
	err = png.Encode(outfile, render.Result.Image)
	if err != nil {
		return fmt.Errorf("Error in png encoding: %w", err)
	}

	fmt.Println("CONVERSION DONE!") // REMOVE: check that the code is running
	return nil
}

/*
Callback for generating a png
*/
func imageToPNG(img image.Image) *bytes.Buffer {
	buf := new(bytes.Buffer)
	_ = png.Encode(buf, img)
	return buf
}

/*
Convert image to single page pdf using fpdf
*/
func ConvertImageToPdf(inputPath string, outputPath string) error {
	file, err := os.Open(inputPath)
	if err != nil {
		return fmt.Errorf("Error reading file: %w", err)
	}
	defer file.Close()

	imgConfig, format, err := image.DecodeConfig(file)
	if err != nil {
		return fmt.Errorf("Error readin image header: %w", err)
	}

	// Create new empty pdf
	pdf := gofpdf.NewCustom(&gofpdf.InitType{
		UnitStr: "pt",
		Size: gofpdf.SizeType{
			Wd: float64(imgConfig.Width),
			Ht: float64(imgConfig.Height),
		},
	})
	pdf.AddPage()

	if format == "jpeg" || format == "png" {
		// use fpdf passthrough for supported formats
		pdf.Image(inputPath, 0, 0, float64(imgConfig.Width), float64(imgConfig.Height), false, "", 0, "")
	} else {
		// Re-encode to png for bmp/tiff/webp
		file.Seek(0, 0)
		img, _, err := image.Decode(file)
		if err != nil {
			return fmt.Errorf("Failed to decode %s: %w", format, err)
		}

		opts := gofpdf.ImageOptions{ImageType: "PNG", ReadDpi: true}
		pdf.RegisterImageOptionsReader("img", opts, imageToPNG(img))
		pdf.ImageOptions("img", 0, 0, float64(imgConfig.Width), float64(imgConfig.Height), false, opts, 0, "")
	}

	if pdf.Error() != nil {
		return fmt.Errorf("Error in PDF generation: %w", pdf.Error())
	}

	return pdf.OutputFileAndClose(outputPath)
}

/*
If adding more modules, consider taking a look at this example project
that extends the PDF engine protocol:
https://github.com/Vrex123/gotenberg-ghostscript/tree/main
*/
func gspreviewRoute(m *Module) api.Route {
	return api.Route{
		Method:      http.MethodPost,
		Path:        "/forms/gspreview",
		IsMultipart: true,
		Handler: func(c echo.Context) error {
			ctx := c.Get("context").(*api.Context)

			form := ctx.FormData()
			inputPaths := []string{}
			outputFormat := ""
			xsize := 1200
			err := form.AnyMandatoryPaths(&inputPaths).Int("xsize", &xsize, 1200).String("outputFormat", &outputFormat, "auto").Validate()
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

			var outputPaths []string
			for _, inputPath := range inputPaths {
				toPng := (outputFormat == "auto" && strings.HasSuffix(strings.ToLower(inputPath), ".pdf")) || outputFormat == "png"
				var outputPath string
				if toPng {
					outputPath = ctx.GeneratePath(".png")
					err = RenderPdfFileToPng(m, inputPath, outputPath, xsize)
					if err != nil {
						return fmt.Errorf("Error pdf->png: %s: %w", inputPath, err)
					}
				} else {
					outputPath = ctx.GeneratePath(".pdf")
					err = ConvertImageToPdf(inputPath, outputPath)
					if err != nil {
						return fmt.Errorf("Error image->pdf: %s: %w", inputPath, err)
					}
				}
				outputPaths = append(outputPaths, outputPath)
			}

			return ctx.AddOutputPaths(outputPaths...)
		},
	}
}

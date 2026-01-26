# Prorenata fork / extensions

This is a fork of [gotenberg](https://github.com/gotenberg/gotenberg) that changes the following:

- Added option to libreOffice route to convert document to plain text. Specifically made for rtf -> text conversions.
- Added route `gspreview` that add conversions by pdfium and gofpdf.
- Removed `chromium` and `pdftk` (including java) from Docker image.

## On branching

- `main`: Keep `main` branch in sync with the main branch of the source repository
- `dev`: Use as "main" branch for this fork

## Feature: rtf conversion to text 

Specify text as output in the request body:
- outputFormat=<pdf, text> (default=pdf)

Example
```
curl --request POST -F files=@test.rtf -F "outputFormat=text" http://localhost:3002/forms/libreoffice/convert -o output.txt
```

*Attention:* is not set up to work with multiple files in the request.

## Feature: PDF/Image Conversion with Ghostscript and GraphicsMagic

Unless specified a PDF is converted to an png-image, all other files are converted to pdf.

Force operation by setting output format in request body:
- outputFormat=<auto, png, pdf> (default=auto)

As with other routes in gotenberg it's possible to convert multiple files in 1 request which will return the result in a zip-archive.

### convert page 1 of pdf to png

Specify x dimension (pixels) in request body:
- xsize=<pixels> (default 1200 pixels)

Example:
```
curl --request POST -F "xsize=600" -F "outputFormat=auto" -F files=@test.pdf http://localhost:3002/forms/gspreview -o preview.png
```

### Feature: convert image to PDF

(takes no extra arguments)

Example
```
curl --request POST -F "outputFormat=auto" -F files=@test1.tiff http://localhost:3002/forms/gspreview -o output.pdf
```

## Testing

Testrunner is updated with tags so to not include removed parts

`make test-unit`

`make test-integration`

To only run our integration test

`make test-prorenata`

## Build/Push image

Set correct GOTENBERG_VERSION (defined in `.env`)

Run build script from gotenberg root path
`make build`

Creates image on `prorenata/gotenberg:GOTENBERG_VERSION`

(fold this part into a separate script)

Push images (Not possible on a normal docker user). 
Example:
`docker push prorenata/gotenberg:v8.24.0-prorenata-dev-amd64`
`docker push prorenata/gotenberg:v8.24.0-prorenata-dev-arm64`

Create multi architecture manifest. 
Example:
`docker manifest create  prorenata/gotenberg:v8.24.0-prorenata-dev --amend  prorenata/gotenberg:v8.24.0-prorenata-dev-amd64 --amend prorenata/gotenberg:v8.24.0-prorenata-dev-arm64`

`docker manifest push   prorenata/gotenberg:v8.24.0-prorenata-dev`


## ENV variables of note

- CHROMIUM_DISABLE_ROUTES: (Already disabled in docker images) 
- API_ENABLE_DEBUG_ROUTE: Enables some debug features. Example: `curl --request GET  http://localhost:3002/debug`
- GOTENBERG_ENABLE_PROMETHEUS: Enable prometheus

## Changelog

### [v8.24.0-prorenata-1.1.0] - 2026-01-26

#### Changed
- Now uses go-pdfium for PDF -> png (using prebuild pdfium library)
- Now uses gofpdf for image -> PDF (known support: tiff, png, bmp, jpg, webp). 

#### Removed
- Removed ghostscript from image
- Removed graphicsmagick from image


### [v8.24.0-prorenata-1.0.0] - 2025-11-06

#### Added

- Project fork documentation
- Document -> plain text conversion via libreOffice
- PDF <-> Image conversion via GraphicsMagick and Ghostscript

#### Changed
- Remove pdftk, java and chromium from Docker image
- Updated test suite to pass with removed modules
- Some changes to build scripts 

# Prorenata fork / extensions

- libreoffice route has an option to convert document to plain text. Specifically made for rtf -> text conversions.
- Added route `gspreview` that add conversions by graphicsmagick and ghopstscript
- Removed `chromium` and `pdftk` (including java) from docker image

## On branching

(work in progress)

- `main`: Keep `main` branch in sync with the main branch of the source reprository
- `dev`: Use as "main" branch for this fork

We can probably get away with 1 single branch as a common dev/stable/main/production branch.

## rtf conversion to text 

Specify text as output in the request body:
- outputFormat=<pdf, text> (default=pdf)

Example
```
curl --request POST -F files=@test.rtf -F "outputFormat=text" http://localhost:3002/forms/libreoffice/convert -o output.txt
```

*Attention:* is not setup to work with multiple files in the request.

## PDF/Image Conversion with ghostscript and graphicsmagic

Unless specified a pdf is converted to an png-image, all other files are converted to pdf.

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

### convert image to pdf

(takes no extra arguments)

Example
```
curl --request POST -F "outputFormat=auto" -F files=@test1.tiff http://localhost:3002/forms/gspreview -o output.pdf
```

## Build/Push image

Set correct GOTENBERG_VERSION (defined in `.env`)

Run build script from gotenberg root path
`make build`

Creates image on `prorenata/gotenberg:GOTENBERG_VERSION`

(fold this part into a separate script)

Push images (may require use of another Docker user such as `prorenataservice`). 
Example:
`docker push prorenata/gotenberg:v8.24.0-prorenata-dev-amd64`
`docker push prorenata/gotenberg:v8.24.0-prorenata-dev-arm64`

Create multi architecture manifest. 
Example:
`docker manifest create  prorenata/gotenberg:v8.24.0-prorenata-dev --amend  prorenata/gotenberg:v8.24.0-prorenata-dev-amd64 --amend prorenata/gotenberg:v8.24.0-prorenata-dev-arm64`

`docker manifest push   prorenata/gotenberg:v8.24.0-prorenata-dev`

## Docker Compose

Http server: port 3000

For local stress test: `docker compose up --scale gotenberg=16`

### ENV variables of note

- CHROMIUM_DISABLE_ROUTES: (Alreafy disabled in docker images) 
- API_ENABLE_DEBUG_ROUTE: Enables some debug features. Example: `curl --request GET  http://localhost:3002/debug`
- GOTENBERG_ENABLE_PROMETHEUS: Enable prometheus

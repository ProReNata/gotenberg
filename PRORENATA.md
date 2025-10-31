# Prorenata fork / extensions

- libreoffice route has an option to convert document to plain text. This is specifically made for rtf -> text conversions.
- Added route `gspreview` that add conversions by graphicsmagick and ghopstscript

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

## PDF/Image Conversion with ghostscript and graphicsmagic

Unless specified a pdf is converted to an png-image, all other files are converted to pdf.

Force mode by setting format in request body:
- outputFormat=<auto, png, pdf> (default=auto)

As with other routes in gotenberg it's possible to convert in batch which will return a zip-archive of all files.

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

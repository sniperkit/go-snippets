# scrubber
[![CircleCI](https://circleci.com/gh/redshiftzero/scrubber.svg?style=svg&circle-token=612ce2eb2ba545a51c7f0e73d4def1f49b431cdf)](https://circleci.com/gh/redshiftzero/scrubber)

Simple CLI to remove metadata from images

*Note:* This is in development and should not be used by high risk users.

```
$ scrubber -help
Usage of scrubber:
  -clean
    	Generate cleaned image (optional) (default true)
  -input string
    	Image file to scrub metadata from (required)
  -json
    	Print JSON metadata to stdout (optional)
  -output string
    	Output file of cleaned image (optional)
```

## Example

### Scrubbing image metadata and printing JSON metadata

```
$ scrubber -input=test/test.jpg -json
Clean version of test/test.jpg saved in the current directory as test_clean.jpg
{
  "Aperture": "5.66 EV (f/7.1)",
  "Color Space": "sRGB",
  "Components Configuration": "Y Cb Cr -",
  "Compression": "JPEG compression",
  "Custom Rendered": "Normal process",
  "Date and Time": "2010:05:28 13:35:33",
  "Date and Time (Digitized)": "2010:05:28 09:42:30",
  "Date and Time (Original)": "2010:05:28 09:42:30",
  "Exif Version": "Exif Version 2.21",
  "Exposure Bias": "0.00 EV",
  "Exposure Mode": "Auto exposure",
  "Exposure Program": "Not defined",
  "Exposure Time": "1/500 sec.",
  "F-Number": "f/7.1",
  "Filename": "test/test.jpg",
  "Flash": "Flash did not fire, compulsory flash mode",
  "FlashPixVersion": "FlashPix Version 1.0",
  "Focal Length": "300.0 mm",
  "Focal Plane Resolution Unit": "Inch",
  "Focal Plane X-Resolution": "3210.946",
  "Focal Plane Y-Resolution": "3230.241",
  "ISO Speed Ratings": "400",
  "Interoperability Index": "R98",
  "Interoperability Version": "0100",
  "Manufacturer": "Canon",
  "Metering Mode": "Pattern",
  "Model": "Canon EOS DIGITAL REBEL XTi",
  "Orientation": "Top-left",
  "Pixel X Dimension": "2816",
  "Pixel Y Dimension": "1880",
  "Resolution Unit": "Inch",
  "Scene Capture Type": "Standard",
  "Shutter Speed": "8.97 EV (1/500 sec.)",
  "Software": "Paint Shop Pro Photo 11.20",
  "User Comment": "",
  "White Balance": "Auto white balance",
  "X-Resolution": "72",
  "Y-Resolution": "72",
  "YCbCr Positioning": "Co-sited"
}
```

### Verifying metadata is gone in scrubbed image

```
$ scrubber -input=test_clean.jpg -json -clean=false
{
  "Filename": "test_clean.jpg"
}
```

## Install

Currently you need `libexif` installed:

```
apt-get install libexif-dev
```

Or on Mac OS:

```
brew install libexif
```

Then you can:

```
go get github.com/redshiftzero/scrubber
```

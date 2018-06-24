
export function fileSizeString(len = 0) {
  let rounded;
  if (len < 1000) {
    return `${len} bytes`;
  } else if (len < (1000 * 1000 * 1000)) {
    rounded = Math.round(len/1000);
    return `${rounded} kb`;
  } else if (len < (1000 * 1000 * 1000 * 1000)) {
    rounded = Math.round(len/1000/1000) / 100;
    return `${rounded} Mb`;
  } else {
    rounded = Math.round(len/1000/1000/1000) / 100;
    return `${rounded} Gb`;
  }
}

export function dateString(d) {
  d = new Date(d)
  return d.toISOString().slice(0,-14);
}

export function fileTypeString(mime) {
  switch (mime) {
    case "application/pdf":
      return "pdf";
    case "application/csv":
      return "csv";
    default:
      return mime;
  }
}
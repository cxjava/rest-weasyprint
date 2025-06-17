# rest-weasyprint

> REST API interface to leverage WeasyPrint — built with Go, wraps the WeasyPrint CLI and supports Docker deployment with advanced features.

This project provides a comprehensive, lightweight, and easy-to-deploy service for converting HTML content (inline, remote, or uploaded files) into PDF documents using [WeasyPrint](https://weasyprint.org/), exposed via a RESTful API. Built in Go, it's designed to be fast, scalable, and deployable via Docker with extensive customization options.

---

## 🚀 Features

### Core Functionality
- **Multiple Input Methods**: HTML files, inline HTML strings, remote URLs
- **Asset Support**: Upload and use custom CSS, fonts, images, and other assets
- **Advanced WeasyPrint Options**: Full support for WeasyPrint's extensive configuration
- **File Sharing Integration**: Built-in support for popular file sharing services
- **Streaming Output**: Efficient PDF generation with streaming response
- **Automatic Cleanup**: Auto-delete temporary files after processing
- **Health Monitoring**: Health check and version endpoints
- **Timeout Control**: Configurable request timeouts
- **Custom Filenames**: Set custom PDF filenames with UTF-8 support

### Supported WeasyPrint Options
- **Encoding & Media**: `encoding`, `media_type`, `base_url`
- **PDF Variants**: `pdf_identifier`, `pdf_variant`, `pdf_version`, `pdf_forms`
- **Output Control**: `uncompressed_pdf`, `custom_metadata`, `presentational_hints`
- **Image Processing**: `srgb`, `optimize_images`, `jpeg_quality`, `dpi`
- **Font Handling**: `full_fonts`, `hinting`
- **Performance**: `cache_folder`, `timeout`
- **Debugging**: `verbose`, `debug`, `quiet`

### File Sharing Services
- **file.io**: Temporary file hosting (14 days)
- **ki.tc**: Anonymous file sharing
- **c-v.sh**: Simple file upload service

---

## 🐳 Docker Deployment

### Quick Start with Docker

```bash
# Pull and run from GitHub Container Registry
docker run -d -p 8080:8080 --name rest-weasyprint ghcr.io/cxjava/rest-weasyprint

# Test if running
curl http://localhost:8080/
```

### Custom Configuration

```bash
docker run -d -p 8080:8080 \
  --env WEB_TIME_OUT_SECOND=120 \
  --name rest-weasyprint \
  ghcr.io/cxjava/rest-weasyprint
```

### Build from Source

```bash
git clone https://github.com/cxjava/rest-weasyprint.git
cd rest-weasyprint
docker build -t rest-weasyprint .
```

---

## 📡 API Endpoints

### Health Check
```
GET /
```
Returns service status and health information.

### Version Information
```
GET /api/v1/pdf/version
```
Returns WeasyPrint version information.

### HTML String Rendering
```
POST /api/v1/pdf/render/html
GET  /api/v1/pdf/render/html  (for testing)
```

### File Upload Rendering
```
POST /api/v1/pdf/render/file
```

---

## 🛠️ Usage Examples

### 1. Basic Health Check

```bash
curl http://localhost:8080/
```

### 2. Get WeasyPrint Version

```bash
curl http://localhost:8080/api/v1/pdf/version
```

### 3. Render Remote URL

```bash
curl -X POST http://localhost:8080/api/v1/pdf/render/html \
  -H "Content-Type: application/json" \
  -d '{"html": "https://weasyprint.org"}' \
  --output website.pdf
```

### 4. Render Inline HTML

```bash
curl -X POST http://localhost:8080/api/v1/pdf/render/html \
  -H "Content-Type: application/json" \
  -d '{
    "html": "<!DOCTYPE html><html><head><title>Test</title></head><body><h1>Hello World! 🌍</h1><p>This is a test document.</p></body></html>"
  }' \
  --output hello.pdf
```

### 5. Render with Custom Filename

```bash
curl -X POST "http://localhost:8080/api/v1/pdf/render/html?filename=custom-report.pdf" \
  -H "Content-Type: application/json" \
  -d '{"html": "<h1>Custom Report</h1><p>Generated on $(date)</p>"}' \
  --output custom-report.pdf
```

### 6. Render with WeasyPrint Options

```bash
curl -X POST http://localhost:8080/api/v1/pdf/render/html \
  -H "Content-Type: application/json" \
  -d '{
    "html": "<h1>High Quality PDF</h1>",
    "options": {
      "encoding": "UTF-8",
      "media_type": "print",
      "dpi": 300,
      "jpeg_quality": 95,
      "optimize_images": true,
      "pdf_version": "1.7"
    }
  }' \
  --output high-quality.pdf
```

### 7. Upload HTML File

```bash
curl -F "html=@document.html" \
  http://localhost:8080/api/v1/pdf/render/file \
  -o document.pdf
```

### 8. Upload with Multiple Assets

```bash
curl -F "html=@invoice.html" \
     -F "css.styles.css=@styles.css" \
     -F "css.print.css=@print.css" \
     -F "asset.logo.png=@logo.png" \
     -F "asset.font.ttf=@custom-font.ttf" \
     -F "options={\"base_url\":\"http://example.com/\",\"dpi\":300};type=application/json" \
     "http://localhost:8080/api/v1/pdf/render/file?filename=invoice-final.pdf" \
     -o invoice.pdf
```

### 9. File Sharing Integration

#### Upload to file.io
```bash
curl -X POST "http://localhost:8080/api/v1/pdf/render/html?share_service=file.io" \
  -H "Content-Type: application/json" \
  -d '{"html": "<h1>Shared Document</h1>"}' \
  | jq '.'
```

#### Upload to ki.tc
```bash
curl -X POST http://localhost:8080/api/v1/pdf/render/html \
  -H "Content-Type: application/json" \
  -d '{
    "html": "<h1>Anonymous Share</h1>",
    "share_service": "ki.tc"
  }' \
  | jq '.'
```

#### Upload to c-v.sh
```bash
curl -F "html=@document.html" \
     "http://localhost:8080/api/v1/pdf/render/file?share_service=c-v.sh&filename=document.pdf" \
     | jq '.'
```

---

## 📋 Request/Response Formats

### HTML Render Request (JSON)
```json
{
  "html": "<html>...</html>",  // HTML content or URL
  "options": {                 // Optional WeasyPrint options
    "encoding": "UTF-8",
    "media_type": "print",
    "base_url": "https://example.com/",
    "dpi": 300,
    "jpeg_quality": 95,
    "optimize_images": true,
    "pdf_version": "1.7",
    "pdf_variant": "pdf/a-1b",
    "uncompressed_pdf": false,
    "custom_metadata": true,
    "presentational_hints": true,
    "srgb": true,
    "full_fonts": false,
    "hinting": true,
    "verbose": false,
    "debug": false,
    "quiet": true,
    "timeout": 30
  },
  "share_service": "file.io"   // Optional: file.io, ki.tc, c-v.sh
}
```

### File Upload Fields
- `html`: Main HTML file (required)
- `css.<filename>`: CSS files (optional, multiple allowed)
- `asset.<filename>`: Asset files like fonts, images (optional, multiple allowed)
- `options`: JSON string with WeasyPrint options (optional)
- `filename`: Custom filename for the PDF (optional)
- `share_service`: Share service for the PDF (optional)

### File Sharing Response
```json
{
  "success": true,
  "link": "https://file.io/abc123",
  "service": "file.io",
  "filename": "document.pdf",
  "message": "Upload successful"
}
```

---

## ⚙️ Configuration Options

### Environment Variables
| Variable | Default | Description |
|----------|---------|-------------|
| `WEB_TIME_OUT_SECOND` | 30 | Request timeout in seconds |

### WeasyPrint Options Reference

#### Basic Options
| Option | Type | Default | Description |
|--------|------|---------|-------------|
| `encoding` | string | - | Character encoding |
| `media_type` | string | print | CSS media type |
| `base_url` | string | - | Base URL for relative links |

#### PDF Options
| Option | Type | Default | Description |
|--------|------|---------|-------------|
| `pdf_identifier` | string | - | PDF identifier |
| `pdf_variant` | string | - | PDF variant (pdf/a-1b, pdf/a-2b, etc.) |
| `pdf_version` | string | - | PDF version |
| `pdf_forms` | boolean | false | Enable PDF forms |

#### Quality Options
| Option | Type | Range | Default | Description |
|--------|------|-------|---------|-------------|
| `jpeg_quality` | integer | 0-95 | 80 | JPEG compression quality |
| `dpi` | integer | 50-600 | 96 | Output resolution |
| `timeout` | integer | 1-300 | 30 | Processing timeout |

#### Output Options
| Option | Type | Default | Description |
|--------|------|---------|-------------|
| `uncompressed_pdf` | boolean | false | Disable PDF compression |
| `custom_metadata` | boolean | false | Include custom metadata |
| `presentational_hints` | boolean | false | Use presentational hints |
| `srgb` | boolean | false | Use sRGB color space |
| `optimize_images` | boolean | false | Optimize embedded images |
| `full_fonts` | boolean | false | Embed full font files |
| `hinting` | boolean | false | Enable font hinting |

#### Debug Options
| Option | Type | Default | Description |
|--------|------|---------|-------------|
| `verbose` | boolean | false | Verbose logging |
| `debug` | boolean | false | Debug mode |
| `quiet` | boolean | false | Suppress output |

---

## 📦 Supported Input Formats

| Type | Format | Supported | Notes |
|------|--------|-----------|-------|
| HTML | Inline string | ✅ | Direct HTML content |
| HTML | Remote URL | ✅ | HTTP/HTTPS URLs |
| HTML | File upload | ✅ | .html files |
| CSS | File upload | ✅ | Multiple files supported |
| Assets | Images | ✅ | PNG, JPG, SVG, etc. |
| Assets | Fonts | ✅ | TTF, OTF, WOFF, etc. |
| Assets | Other | ✅ | Any file type |

---

## 🔧 Development

### Local Development
```bash
# With auto-reload (requires air)
air

# Direct run
cd cmd
go run main.go

# Build and run
go build -o rest-weasyprint cmd/main.go
./rest-weasyprint
```

### Requirements
- Go 1.20+
- WeasyPrint installed
- Docker (for containerized deployment)

---

## 🚨 Error Handling

The service provides detailed error messages for common issues:

- **400 Bad Request**: Invalid input, missing files, malformed JSON
- **500 Internal Server Error**: WeasyPrint execution errors, file system issues
- **Timeout**: Requests exceeding configured timeout limit

---

## 🤝 Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Add tests if applicable
5. Submit a pull request

---

## 📄 License

MIT License - see LICENSE file for details.

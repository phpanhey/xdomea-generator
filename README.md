
# Xdomea Generator

A drag-and-drop generator for Xdomea templates with a Go backend. Configure templates in the browser and export Xdomea.

## Features

- Drag-and-drop template configuration.
- Export Xdomea templates with one click.
- Go backend generates Xdomea from frontend.

## Project Structure
```

/ui
└─ index.html       # Frontend interface
/backend
└─ main.go          # Go server for Xdomea generation

````

## Getting Started

### 1. Start the Backend
```bash
cd backend
go build -o xdomea-server
./xdomea-server
````

* Runs on `http://localhost:8080`
* Endpoint: `/generatexdomea`

### 2. Open the Frontend

* Open `ui/index.html` in your browser.
* Configure your template and click **Export Xdomea**.

## License
MIT License

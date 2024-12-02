# Evolution Audio Converter

This project is a microservice in Go that processes audio files, converts them to **opus** or **mp3** format, and returns both the duration of the audio and the converted file (as base64 or S3 URL). The service accepts audio files sent as **form-data**, **base64**, or **URL**.

## Requirements

Before starting, you'll need to have the following installed:

- [Go](https://golang.org/doc/install) (version 1.21 or higher)
- [Docker](https://docs.docker.com/get-docker/) (to run the project in a container)
- [FFmpeg](https://ffmpeg.org/download.html) (for audio processing)

## Installation

### Clone the Repository

Clone this repository to your local machine:

```bash
git clone https://github.com/EvolutionAPI/evolution-audio-converter.git
cd evolution-audio-converter
```

### Install Dependencies

Install the project dependencies:

```bash
go mod tidy
```

### Install FFmpeg

The service depends on **FFmpeg** to convert the audio. Make sure FFmpeg is installed on your system.

- On Ubuntu:

  ```bash
  sudo apt update
  sudo apt install ffmpeg
  ```

- On macOS (via Homebrew):

  ```bash
  brew install ffmpeg
  ```

- On Windows, download FFmpeg [here](https://ffmpeg.org/download.html) and add it to your system `PATH`.

### Configuration

Create a `.env` file in the project's root directory. Here are the available configuration options:

#### Basic Configuration

```env
PORT=4040
API_KEY=your_secret_api_key_here
```

#### Transcription Configuration

```env
ENABLE_TRANSCRIPTION=true
TRANSCRIPTION_PROVIDER=openai  # or groq
OPENAI_API_KEY=your_openai_key_here
GROQ_API_KEY=your_groq_key_here
TRANSCRIPTION_LANGUAGE=en  # Default transcription language (optional)
```

#### Storage Configuration

```env
ENABLE_S3_STORAGE=true
S3_ENDPOINT=play.min.io
S3_ACCESS_KEY=your_access_key_here
S3_SECRET_KEY=your_secret_key_here
S3_BUCKET_NAME=audio-files
S3_REGION=us-east-1
S3_USE_SSL=true
S3_URL_EXPIRATION=24h
```

### Storage Options

The service supports two storage modes for the converted audio:

1. **Base64 (default)**: Returns the audio file encoded in base64 format
2. **S3 Compatible Storage**: Uploads to S3-compatible storage (AWS S3, MinIO, etc.) and returns a presigned URL

When S3 storage is enabled, the response will include a `url` instead of the `audio` field:

```json
{
  "duration": 120,
  "format": "ogg",
  "url": "https://your-s3-endpoint/bucket/file.ogg?signature...",
  "transcription": "Transcribed text here..." // if transcription was requested
}
```

If S3 upload fails, the service automatically falls back to base64 encoding.

## Running the Project

### Locally

To run the service locally:

```bash
go run main.go -dev
```

The server will be available at `http://localhost:4040`.

### Using Docker

1. **Build the Docker image**:

   ```bash
   docker build -t audio-service .
   ```

2. **Run the container**:

   ```bash
   docker run -p 4040:4040 --env-file=.env audio-service
   ```

## API Usage

### Authentication

All requests must include the `apikey` header with your API key.

### Endpoints

#### Process Audio

`POST /process-audio`

Accepts audio files in these formats:

- Form-data
- Base64
- URL

Optional parameters:

- `format`: Output format (`mp3` or `ogg`, default: `ogg`)
- `transcribe`: Enable transcription (`true` or `false`)
- `language`: Transcription language code (e.g., "en", "es", "pt")

#### Transcribe Only

`POST /transcribe`

Transcribes audio without format conversion.

Optional parameters:

- `language`: Transcription language code

### Example Requests

#### Form-data Upload

```bash
curl -X POST -F "file=@audio.mp3" \
  -F "format=ogg" \
  -F "transcribe=true" \
  -F "language=en" \
  http://localhost:4040/process-audio \
  -H "apikey: your_secret_api_key_here"
```

#### Base64 Upload

```bash
curl -X POST \
  -d "base64=$(base64 audio.mp3)" \
  -d "format=ogg" \
  http://localhost:4040/process-audio \
  -H "apikey: your_secret_api_key_here"
```

#### URL Upload

```bash
curl -X POST \
  -d "url=https://example.com/audio.mp3" \
  -d "format=ogg" \
  http://localhost:4040/process-audio \
  -H "apikey: your_secret_api_key_here"
```

### Response Format

With S3 storage disabled (default):

```json
{
  "duration": 120,
  "audio": "UklGR... (base64 of the file)",
  "format": "ogg",
  "transcription": "Transcribed text here..." // if requested
}
```

With S3 storage enabled:

```json
{
  "duration": 120,
  "url": "https://your-s3-endpoint/bucket/file.ogg?signature...",
  "format": "ogg",
  "transcription": "Transcribed text here..." // if requested
}
```

## License

This project is licensed under the [MIT](LICENSE) license.

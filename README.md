# Evolution Audio Converter

This project is a microservice in Go that processes audio files, converts them to **opus** or **mp3** format, and returns both the duration of the audio and the converted file in base64. The service accepts audio files sent as **form-data**, **base64**, or **URL**.

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

Create a `.env` file in the project's root directory with the following configuration:

```env
PORT=4040
API_KEY=your_secret_api_key_here
```

### Transcription Configuration

To enable audio transcription, configure the following variables in the `.env` file:

```env
ENABLE_TRANSCRIPTION=true
TRANSCRIPTION_PROVIDER=openai  # or groq
OPENAI_API_KEY=your_openai_key_here
GROQ_API_KEY=your_groq_key_here
TRANSCRIPTION_LANGUAGE=en  # Default transcription language (optional)
```

- `ENABLE_TRANSCRIPTION`: Enables or disables the transcription feature
- `TRANSCRIPTION_PROVIDER`: Chooses the AI provider for transcription (openai or groq)
- `OPENAI_API_KEY`: Your OpenAI API key (required if using openai)
- `GROQ_API_KEY`: Your Groq API key (required if using groq)
- `TRANSCRIPTION_LANGUAGE`: Sets the default transcription language (optional)

## Running the Project

### Locally

To run the service locally, use the following command:

```bash
go run main.go -dev
```

The server will be available at `http://localhost:4040`.

### Using Docker

If you prefer to run the service in a Docker container, follow the steps below:

1. **Build the Docker image**:

   ```bash
   docker build -t audio-service .
   ```

2. **Run the container**:

   ```bash
   docker run -p 4040:4040 --env-file=.env audio-service
   ```

   This will start the container on the port specified in the `.env` file.

## How to Use

You can send `POST` requests to the `/process-audio` endpoint with an audio file in the following formats:

- **Form-data** (to upload files)
- **Base64** (to send the audio encoded in base64)
- **URL** (to send the link to the audio file)

### Authentication

All requests must include the `apikey` header with the value of the `API_KEY` configured in the `.env` file.

### Optional Parameters

- **`format`**: You can specify the format for conversion by passing the `format` parameter in the request. Supported values:
  - `mp3`
  - `ogg` (default)

### Audio Transcription

You can get the audio transcription in two ways:

1. Along with audio processing by adding the `transcribe=true` parameter:

```bash
curl -X POST -F "file=@audio.mp3" \
  -F "transcribe=true" \
  -F "language=en" \
  http://localhost:4040/process-audio \
  -H "apikey: your_secret_api_key_here"
```

2. Using the specific transcription endpoint:

```bash
curl -X POST -F "file=@audio.mp3" \
  -F "language=en" \
  http://localhost:4040/transcribe \
  -H "apikey: your_secret_api_key_here"
```

Optional parameters:
- `language`: Audio language code (e.g., "en", "es", "pt"). If not specified, it will use the value defined in `TRANSCRIPTION_LANGUAGE` in `.env`. If neither is defined, the system will try to automatically detect the language.

The response will include the `transcription` field with the transcribed text:

```json
{
  "transcription": "Transcribed text here..."
}
```

When used with audio processing (`/process-audio`), the response will include both audio data and transcription:

```json
{
  "duration": 120,
  "audio": "UklGR... (base64 of the file)",
  "format": "ogg",
  "transcription": "Transcribed text here..."
}
```

### Example Requests Using cURL

#### Sending as Form-data

```bash
curl -X POST -F "file=@path/to/audio.mp3" http://localhost:4040/process-audio \
  -F "format=ogg" \
  -H "apikey: your_secret_api_key_here"
```

#### Sending as Base64

```bash
curl -X POST -d "base64=$(base64 path/to/audio.mp3)" http://localhost:4040/process-audio \
  -d "format=ogg" \
  -H "apikey: your_secret_api_key_here"
```

#### Sending as URL

```bash
curl -X POST -d "url=https://example.com/path/to/audio.mp3" http://localhost:4040/process-audio \
  -d "format=ogg" \
  -H "apikey: your_secret_api_key_here"
```

### Response

The response will be a JSON object containing the audio duration and the converted audio file in base64:

```json
{
  "duration": 120,
  "audio": "UklGR... (base64 of the file)",
  "format": "ogg"
}
```

- `duration`: The audio duration in seconds.
- `audio`: The converted audio file encoded in base64.
- `format`: The format of the converted file (`mp3` or `ogg`).

## License

This project is licensed under the [MIT](LICENSE) license.
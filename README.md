# Fampay-Assignment

## YouTube Video Fetcher API

A scalable API service that continuously fetches latest YouTube videos for a given search query and stores them in a database. The videos can be retrieved through a paginated API endpoint.

## Features

- Asynchronous background fetching of YouTube videos
- Paginated API endpoints to retrieve stored videos
- Smart API key rotation system with automatic fallback when quota exhausts
- PostgreSQL database for persistent storage
- Optimized database queries with proper indexing
- Reverse chronological sorting of videos by publish date

## Tech Stack

- Go (Golang)
- Gin Web Framework
- GORM (PostgreSQL)
- YouTube Data API v3

## Project Structure

```
├── cmd
│   └── server
│       └── main.go
├── config
├── go.mod
├── go.sum
├── internal
│   ├── config
│   │   └── config.go
│   ├── db
│   │   └── postgres_db
│   │       └── postgres.go
│   ├── helpers
│   │   └── int_to_base_64.go
│   ├── models
│   │   └── models.go
│   ├── router
│   │   └── router.go
│   └── videos
│       ├── handler.go
│       ├── routes.go
│       └── youtube_service.go
├── output.csv
└── README.md
```

## Data Models

### Video Model
```go
type Video struct {
    gorm.Model
    VideoID       string    `gorm:"not null;primaryKey" json:"videoId"`
    Title         string    `gorm:"not null" json:"title"`
    Description   *string   `gorm:"type:text" json:"description"`
    PublishedAt   time.Time `gorm:"type:timestamp;not null" json:"publishedAt"`
    ThumbnailURLs *string   `gorm:"type:text" json:"thumbnails"`
}
```

## Environment Variables

Create a `.env` file in the root directory with the following variables:

```env
YOUTUBE_API_KEY1=your_first_api_key
YOUTUBE_API_KEY2=your_second_api_key
YOUTUBE_API_KEY3=your_third_api_key
POSTGRES_SERVICE_URI=postgres://username:password@host:port/dbname
```

### Multiple API Key Support

The application implements a robust API key rotation system:

1. **Multiple Key Configuration**: Configure up to three YouTube API keys through environment variables
2. **Automatic Rotation**: When one key's quota is exhausted, the system automatically switches to the next available key
3. **Quota Management**: 
   - Monitors HTTP response codes (403, 429) to detect quota exhaustion
   - Implements a 5-second delay before trying the next key
   - Logs key usage and rotation events for monitoring
4. **Failsafe Mechanism**: If all API keys are exhausted, the system logs a fatal error to prevent unnecessary API calls

## Database Schema

The service uses a `videos` table with the following structure:

```sql
CREATE TABLE videos (
    id BIGSERIAL PRIMARY KEY,
    created_at TIMESTAMP WITH TIME ZONE,
    updated_at TIMESTAMP WITH TIME ZONE,
    deleted_at TIMESTAMP WITH TIME ZONE,
    video_id VARCHAR NOT NULL PRIMARY KEY,
    title VARCHAR NOT NULL,
    description TEXT,
    published_at TIMESTAMP NOT NULL,
    thumbnail_urls TEXT
);

CREATE INDEX idx_videos_published_at ON videos(published_at DESC);
CREATE INDEX idx_videos_video_id ON videos(video_id);
```

## API Endpoints

### 1. Get Videos (Paginated)
```
GET /fam/videos?limit={limit}&pageToken={pageToken}
```

**Parameters:**
- `limit` (required): Number of videos to return per page
- `pageToken` (optional): Token for fetching next page of results

**Sample Response:**
```json
{
    "pageToken": "BQ==",
    "videos": [
        {
            "ID": 5,
            "CreatedAt": "2024-12-31T08:23:08.955754+05:30",
            "UpdatedAt": "2024-12-31T08:23:08.955754+05:30",
            "DeletedAt": null,
            "videoId": "-FMbrt4Pvy0",
            "title": "Sample Video Title",
            "description": "Video description...",
            "publishedAt": "2024-12-30T11:06:07Z",
            "thumbnails": "[\"url1\",\"url2\",\"url3\"]"
        }
    ]
}
```

### 2. Get Video by ID
```
GET /fam/videos/{videoId}
```

**Sample Response:**
```json
{
    "ID": 1,
    "CreatedAt": "2024-12-31T08:23:08.955754+05:30",
    "UpdatedAt": "2024-12-31T08:23:08.955754+05:30",
    "DeletedAt": null,
    "videoId": "o4Xkt62NQfQ",
    "title": "Sample Video Title",
    "description": "Video description...",
    "publishedAt": "2024-12-30T07:49:11Z",
    "thumbnails": "[\"url1\",\"url2\",\"url3\"]"
}
```

## Installation & Setup

1. Clone the repository:
```bash
git clone https://github.com/harsh082ip/Fampay-Assignment
cd Fampay-Assignment
```

2. Install dependencies:
```bash
go mod download
```

3. Set up your PostgreSQL database and configure the environment variables in `.env` file.

4. Run the application:
```bash
go run cmd/server/main.go
```

## Implementation Details

### Background Video Fetching
- Continuous fetching every 10 seconds
- Batch processing for database insertions (100 videos per batch)
- Automatic API key rotation on quota exhaustion
- Error handling and logging

### Database Optimization
- Indexed `published_at` for efficient sorting
- Indexed `video_id` for quick lookups
- Soft delete support via `gorm.Model`
- Cursor-based pagination for consistent results

### Error Handling
- Proper HTTP status codes
- Detailed error messages
- Validation for required parameters
- Graceful API key rotation

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.
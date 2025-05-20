# Community API Endpoints

This document outlines the API endpoints for the community feature of the algo trading platform.

## Authentication

All authenticated endpoints require authentication via Bearer token.

```
Authorization: Bearer <token>
```

## Public Endpoints (No Authentication Required)

### 1. List Topics

Returns a list of all community topics with optional filtering.

```bash
curl -X GET \
  'https://api.example.com/community/topics?page=1&perPage=10&category=1&tag=2&sort=recent&search=keyword' \
  -H 'Content-Type: application/json'
```

**Query Parameters:**
- `page`: Page number (default: 1)
- `perPage`: Number of items per page (default: 10)
- `category`: Filter by category ID
- `tag`: Filter by tag ID
- `sort`: Sort field (options: recent, popular, views)
- `search`: Search topics by keyword

**Response:**
```json
{
  "success": true,
  "data": {
    "items": [
      {
        "id": 123,
        "title": "Trading Strategy Discussion",
        "excerpt": "Let's discuss the new moving average strategy...",
        "author": {
          "id": 456,
          "name": "John Doe"
        },
        "categoryId": 2,
        "categoryName": "Technical Analysis",
        "tags": [
          {
            "id": 3,
            "name": "Moving Average"
          },
          {
            "id": 5,
            "name": "Trend Following"
          }
        ],
        "postCount": 15,
        "viewCount": 230,
        "likeCount": 45,
        "createdAt": "2025-04-15T08:23:45Z",
        "updatedAt": "2025-04-18T14:12:30Z",
        "isPinned": false
      }
    ],
    "total": 120,
    "page": 1,
    "perPage": 10,
    "totalPages": 12
  }
}
```

### 2. Get Topic by ID

Returns detailed information about a specific topic.

```bash
curl -X GET \
  'https://api.example.com/community/topics/123' \
  -H 'Content-Type: application/json'
```

**Response:**
```json
{
  "success": true,
  "data": {
    "id": 123,
    "title": "Trading Strategy Discussion",
    "content": "Let's discuss the new moving average strategy...",
    "author": {
      "id": 456,
      "name": "John Doe",
      "profileImage": "https://example.com/profiles/johndoe.jpg",
      "reputation": 1250
    },
    "categoryId": 2,
    "categoryName": "Technical Analysis",
    "tags": [
      {
        "id": 3,
        "name": "Moving Average"
      },
      {
        "id": 5,
        "name": "Trend Following"
      }
    ],
    "postCount": 15,
    "viewCount": 231,
    "likeCount": 45,
    "createdAt": "2025-04-15T08:23:45Z",
    "updatedAt": "2025-04-18T14:12:30Z",
    "isPinned": false
  }
}
```

### 3. List Categories

Returns a list of all community categories.

```bash
curl -X GET \
  'https://api.example.com/community/categories' \
  -H 'Content-Type: application/json'
```

**Response:**
```json
{
  "success": true,
  "data": [
    {
      "id": 1,
      "name": "General Discussion",
      "description": "General trading discussions",
      "topicCount": 45
    },
    {
      "id": 2,
      "name": "Technical Analysis",
      "description": "Discussions about technical indicators and chart patterns",
      "topicCount": 32
    },
    {
      "id": 3,
      "name": "Algorithmic Trading",
      "description": "Discussions about algo trading strategies",
      "topicCount": 28
    },
    {
      "id": 4,
      "name": "Market News",
      "description": "Latest news affecting the markets",
      "topicCount": 15
    }
  ]
}
```

### 4. Get Posts for a Topic

Returns all posts in a specific topic.

```bash
curl -X GET \
  'https://api.example.com/community/topics/123/posts?page=1&perPage=20' \
  -H 'Content-Type: application/json'
```

**Query Parameters:**
- `page`: Page number (default: 1)
- `perPage`: Number of items per page (default: 20)

**Response:**
```json
{
  "success": true,
  "data": {
    "items": [
      {
        "id": 456,
        "content": "I think this strategy works well in trending markets...",
        "author": {
          "id": 789,
          "name": "Jane Smith",
          "profileImage": "https://example.com/profiles/janesmith.jpg",
          "reputation": 780
        },
        "likeCount": 12,
        "parentPostId": null,
        "createdAt": "2025-04-15T09:45:22Z",
        "updatedAt": "2025-04-15T09:45:22Z"
      },
      {
        "id": 457,
        "content": "Yes, but it fails in sideways markets.",
        "author": {
          "id": 456,
          "name": "John Doe",
          "profileImage": "https://example.com/profiles/johndoe.jpg",
          "reputation": 1250
        },
        "likeCount": 8,
        "parentPostId": 456,
        "createdAt": "2025-04-15T10:12:30Z",
        "updatedAt": "2025-04-15T10:12:30Z"
      }
    ],
    "total": 15,
    "page": 1,
    "perPage": 20,
    "totalPages": 1
  }
}
```

### 5. Get User Profile

Returns the community profile for a specific user.

```bash
curl -X GET \
  'https://api.example.com/community/users/456/profile' \
  -H 'Content-Type: application/json'
```

**Response:**
```json
{
  "success": true,
  "data": {
    "id": 456,
    "name": "John Doe",
    "profileImage": "https://example.com/profiles/johndoe.jpg",
    "bio": "Professional trader with 10 years of experience",
    "joinedDate": "2023-05-12T09:30:15Z",
    "reputation": 1250,
    "badgeCount": 8,
    "topicCount": 24,
    "postCount": 153
  }
}
```

### 6. Get User Topics

Returns all topics created by a specific user.

```bash
curl -X GET \
  'https://api.example.com/community/users/456/topics?page=1&perPage=10' \
  -H 'Content-Type: application/json'
```

**Query Parameters:**
- `page`: Page number (default: 1)
- `perPage`: Number of items per page (default: 10)

**Response:**
```json
{
  "success": true,
  "data": {
    "items": [
      {
        "id": 123,
        "title": "Trading Strategy Discussion",
        "excerpt": "Let's discuss the new moving average strategy...",
        "categoryId": 2,
        "categoryName": "Technical Analysis",
        "postCount": 15,
        "viewCount": 230,
        "likeCount": 45,
        "createdAt": "2025-04-15T08:23:45Z",
        "updatedAt": "2025-04-18T14:12:30Z"
      }
    ],
    "total": 24,
    "page": 1,
    "perPage": 10,
    "totalPages": 3
  }
}
```

### 7. Get User Badges

Returns all badges earned by a specific user.

```bash
curl -X GET \
  'https://api.example.com/community/users/456/badges' \
  -H 'Content-Type: application/json'
```

**Response:**
```json
{
  "success": true,
  "data": [
    {
      "id": 1,
      "name": "Helpful",
      "description": "Provided 10 answers marked as helpful",
      "iconUrl": "https://example.com/badges/helpful.png",
      "earnedAt": "2024-06-18T14:12:30Z"
    },
    {
      "id": 3,
      "name": "Strategy Expert",
      "description": "Created 5 highly-rated strategy discussions",
      "iconUrl": "https://example.com/badges/strategy-expert.png",
      "earnedAt": "2024-08-22T09:45:12Z"
    }
  ]
}
```

## Authenticated Endpoints (Require Authentication)

### 8. Create a Topic

Creates a new topic in the community.

```bash
curl -X POST \
  'https://api.example.com/community/topics' \
  -H 'Content-Type: application/json' \
  -H 'Authorization: Bearer YOUR_ACCESS_TOKEN' \
  -d '{
    "title": "New Trading Strategy Discussion",
    "content": "Let's discuss the new moving average strategy...",
    "categoryId": 2,
    "tagIds": [3, 5, 8],
    "isPinned": false
  }'
```

**Request Body:**
- `title`: Topic title (required)
- `content`: Topic content (required)
- `categoryId`: Category ID (required)
- `tagIds`: Array of tag IDs (optional)
- `isPinned`: Whether the topic should be pinned (optional, admin only)

**Response:**
```json
{
  "success": true,
  "message": "Topic created successfully",
  "data": {
    "id": 234,
    "title": "New Trading Strategy Discussion",
    "excerpt": "Let's discuss the new moving average strategy...",
    "author": {
      "id": 456,
      "name": "John Doe"
    },
    "categoryId": 2,
    "categoryName": "Technical Analysis",
    "tags": [
      {
        "id": 3,
        "name": "Moving Average"
      },
      {
        "id": 5,
        "name": "Trend Following"
      },
      {
        "id": 8,
        "name": "Backtesting"
      }
    ],
    "postCount": 0,
    "viewCount": 0,
    "likeCount": 0,
    "createdAt": "2025-05-20T10:12:30Z",
    "updatedAt": "2025-05-20T10:12:30Z",
    "isPinned": false
  }
}
```

### 9. Update a Topic

Updates an existing topic.

```bash
curl -X PUT \
  'https://api.example.com/community/topics/123' \
  -H 'Content-Type: application/json' \
  -H 'Authorization: Bearer YOUR_ACCESS_TOKEN' \
  -d '{
    "title": "Updated Trading Strategy Discussion",
    "content": "Let's continue discussing the moving average strategy...",
    "categoryId": 2,
    "tagIds": [3, 5, 8, 10],
    "isPinned": true
  }'
```

**Request Body:**
- `title`: Topic title (optional)
- `content`: Topic content (optional)
- `categoryId`: Category ID (optional)
- `tagIds`: Array of tag IDs (optional)
- `isPinned`: Whether the topic should be pinned (optional, admin only)

**Response:**
```json
{
  "success": true,
  "message": "Topic updated successfully",
  "data": {
    "id": 123,
    "title": "Updated Trading Strategy Discussion",
    "excerpt": "Let's continue discussing the moving average strategy...",
    "author": {
      "id": 456,
      "name": "John Doe"
    },
    "categoryId": 2,
    "categoryName": "Technical Analysis",
    "tags": [
      {
        "id": 3,
        "name": "Moving Average"
      },
      {
        "id": 5,
        "name": "Trend Following"
      },
      {
        "id": 8,
        "name": "Backtesting"
      },
      {
        "id": 10,
        "name": "Risk Management"
      }
    ],
    "postCount": 15,
    "viewCount": 230,
    "likeCount": 45,
    "createdAt": "2025-04-15T08:23:45Z",
    "updatedAt": "2025-05-20T10:30:15Z",
    "isPinned": true
  }
}
```

### 10. Delete a Topic

Deletes a topic.

```bash
curl -X DELETE \
  'https://api.example.com/community/topics/123' \
  -H 'Authorization: Bearer YOUR_ACCESS_TOKEN'
```

**Response:**
```json
{
  "success": true,
  "message": "Topic deleted successfully"
}
```

### 11. Create a Post in a Topic

Creates a new post in a topic.

```bash
curl -X POST \
  'https://api.example.com/community/topics/123/posts' \
  -H 'Content-Type: application/json' \
  -H 'Authorization: Bearer YOUR_ACCESS_TOKEN' \
  -d '{
    "content": "I think this strategy works well in trending markets...",
    "parentPostId": null
  }'
```

**Request Body:**
- `content`: Post content (required)
- `parentPostId`: ID of the parent post if this is a reply (optional)

**Response:**
```json
{
  "success": true,
  "message": "Post created successfully",
  "data": {
    "id": 789,
    "content": "I think this strategy works well in trending markets...",
    "author": {
      "id": 456,
      "name": "John Doe",
      "profileImage": "https://example.com/profiles/johndoe.jpg",
      "reputation": 1250
    },
    "likeCount": 0,
    "parentPostId": null,
    "createdAt": "2025-05-20T10:45:22Z",
    "updatedAt": "2025-05-20T10:45:22Z"
  }
}
```

### 12. Update a Post

Updates an existing post.

```bash
curl -X PUT \
  'https://api.example.com/community/posts/456' \
  -H 'Content-Type: application/json' \
  -H 'Authorization: Bearer YOUR_ACCESS_TOKEN' \
  -d '{
    "content": "Updated: I think this strategy works well only in trending markets..."
  }'
```

**Request Body:**
- `content`: Updated post content (required)

**Response:**
```json
{
  "success": true,
  "message": "Post updated successfully",
  "data": {
    "id": 456,
    "content": "Updated: I think this strategy works well only in trending markets...",
    "author": {
      "id": 789,
      "name": "Jane Smith",
      "profileImage": "https://example.com/profiles/janesmith.jpg",
      "reputation": 780
    },
    "likeCount": 12,
    "parentPostId": null,
    "createdAt": "2025-04-15T09:45:22Z",
    "updatedAt": "2025-05-20T11:12:30Z"
  }
}
```

### 13. Delete a Post

Deletes a post.

```bash
curl -X DELETE \
  'https://api.example.com/community/posts/456' \
  -H 'Authorization: Bearer YOUR_ACCESS_TOKEN'
```

**Response:**
```json
{
  "success": true,
  "message": "Post deleted successfully"
}
```

### 14. Like Content (Topic or Post)

Likes a topic or post.

```bash
curl -X POST \
  'https://api.example.com/community/like' \
  -H 'Content-Type: application/json' \
  -H 'Authorization: Bearer YOUR_ACCESS_TOKEN' \
  -d '{
    "contentType": "topic",
    "contentId": 123
  }'
```

**Request Body:**
- `contentType`: Type of content to like (required, values: "topic" or "post")
- `contentId`: ID of the content to like (required)

**Response:**
```json
{
  "success": true,
  "message": "Content liked successfully",
  "data": {
    "contentType": "topic",
    "contentId": 123,
    "likeCount": 46
  }
}
```

### 15. Unlike Content (Topic or Post)

Unlikes a previously liked topic or post.

```bash
curl -X DELETE \
  'https://api.example.com/community/unlike' \
  -H 'Content-Type: application/json' \
  -H 'Authorization: Bearer YOUR_ACCESS_TOKEN' \
  -d '{
    "contentType": "post",
    "contentId": 456
  }'
```

**Request Body:**
- `contentType`: Type of content to unlike (required, values: "topic" or "post")
- `contentId`: ID of the content to unlike (required)

**Response:**
```json
{
  "success": true,
  "message": "Content unliked successfully",
  "data": {
    "contentType": "post",
    "contentId": 456,
    "likeCount": 11
  }
}
```

### 16. Bookmark a Topic

Bookmarks a topic for the authenticated user.

```bash
curl -X POST \
  'https://api.example.com/community/topics/123/bookmark' \
  -H 'Authorization: Bearer YOUR_ACCESS_TOKEN'
```

**Response:**
```json
{
  "success": true,
  "message": "Topic bookmarked successfully"
}
```

### 17. Unbookmark a Topic

Removes a bookmark from a topic.

```bash
curl -X DELETE \
  'https://api.example.com/community/topics/123/bookmark' \
  -H 'Authorization: Bearer YOUR_ACCESS_TOKEN'
```

**Response:**
```json
{
  "success": true,
  "message": "Topic unbookmarked successfully"
}
```

### 18. Get User Bookmarks

Returns all bookmarked topics for the authenticated user.

```bash
curl -X GET \
  'https://api.example.com/community/bookmarks?page=1&perPage=20' \
  -H 'Authorization: Bearer YOUR_ACCESS_TOKEN'
```

**Query Parameters:**
- `page`: Page number (default: 1)
- `perPage`: Number of items per page (default: 20)

**Response:**
```json
{
  "success": true,
  "data": {
    "items": [
      {
        "id": 123,
        "title": "Trading Strategy Discussion",
        "excerpt": "Let's discuss the new moving average strategy...",
        "author": {
          "id": 456,
          "name": "John Doe"
        },
        "categoryId": 2,
        "categoryName": "Technical Analysis",
        "tags": [
          {
            "id": 3,
            "name": "Moving Average"
          },
          {
            "id": 5,
            "name": "Trend Following"
          }
        ],
        "postCount": 15,
        "viewCount": 230,
        "likeCount": 46,
        "createdAt": "2025-04-15T08:23:45Z",
        "updatedAt": "2025-04-18T14:12:30Z",
        "bookmarkedAt": "2025-05-18T09:45:30Z"
      }
    ],
    "total": 8,
    "page": 1,
    "perPage": 20,
    "totalPages": 1
  }
}
```
// Package docs Code generated by swaggo/swag. DO NOT EDIT
package docs

import "github.com/swaggo/swag"

const docTemplate = `{
    "schemes": {{ marshal .Schemes }},
    "swagger": "2.0",
    "info": {
        "description": "{{escape .Description}}",
        "title": "{{.Title}}",
        "contact": {},
        "version": "{{.Version}}"
    },
    "host": "{{.Host}}",
    "basePath": "{{.BasePath}}",
    "paths": {
        "/jwt_login": {
            "post": {
                "description": "Login using JWT authentication",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Auth"
                ],
                "summary": "JWT Login",
                "parameters": [
                    {
                        "description": "JWT login request",
                        "name": "request",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/jwtLogin.Request"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/jwtLogin.Response"
                        }
                    }
                }
            }
        },
        "/refresh": {
            "post": {
                "description": "Refresh JWT token",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Auth"
                ],
                "summary": "Refresh JWT",
                "parameters": [
                    {
                        "description": "JWT refresh request",
                        "name": "request",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/refresh.Request"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/refresh.Response"
                        }
                    }
                }
            }
        },
        "/session_login": {
            "post": {
                "description": "Login using session-based authentication",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Auth"
                ],
                "summary": "Session Login",
                "parameters": [
                    {
                        "description": "Session login request",
                        "name": "request",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/sessionLogin.Request"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/sessionLogin.Response"
                        }
                    }
                }
            }
        },
        "/users": {
            "get": {
                "description": "Get a list of all users",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Users"
                ],
                "summary": "Get All Users",
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "array",
                            "items": {
                                "$ref": "#/definitions/getAllUsers.Response"
                            }
                        }
                    }
                }
            },
            "post": {
                "description": "Create a new user",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Users"
                ],
                "summary": "Create User",
                "parameters": [
                    {
                        "description": "Create user request",
                        "name": "request",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/createUser.Request"
                        }
                    }
                ],
                "responses": {
                    "201": {
                        "description": "Created",
                        "schema": {
                            "$ref": "#/definitions/createUser.Response"
                        }
                    }
                }
            }
        },
        "/users/{userID}": {
            "get": {
                "description": "Get user by ID",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Users"
                ],
                "summary": "Get User",
                "parameters": [
                    {
                        "type": "string",
                        "description": "User ID",
                        "name": "userID",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/getUser.Response"
                        }
                    }
                }
            },
            "put": {
                "description": "Update user information by ID",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Users"
                ],
                "summary": "Update User",
                "parameters": [
                    {
                        "type": "string",
                        "description": "User ID",
                        "name": "userID",
                        "in": "path",
                        "required": true
                    },
                    {
                        "description": "Update user request",
                        "name": "request",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/updateUser.Request"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/updateUser.Response"
                        }
                    }
                }
            },
            "delete": {
                "description": "Delete user by ID",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Users"
                ],
                "summary": "Delete User",
                "parameters": [
                    {
                        "type": "string",
                        "description": "User ID",
                        "name": "userID",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "204": {
                        "description": "No Content"
                    }
                }
            }
        },
        "/v1/logout": {
            "get": {
                "description": "Logout and invalidate JWT token",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Auth"
                ],
                "summary": "JWT Logout",
                "responses": {
                    "200": {
                        "description": "OK"
                    }
                }
            }
        },
        "/v1/posts": {
            "get": {
                "description": "Get a list of all posts",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Posts"
                ],
                "summary": "Get All Posts",
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "array",
                            "items": {
                                "$ref": "#/definitions/getAllPosts.Response"
                            }
                        }
                    }
                }
            },
            "post": {
                "description": "Create a new post",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Posts"
                ],
                "summary": "Create Post",
                "parameters": [
                    {
                        "description": "Create post request",
                        "name": "request",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/createPost.Request"
                        }
                    }
                ],
                "responses": {
                    "201": {
                        "description": "Created",
                        "schema": {
                            "$ref": "#/definitions/createPost.Response"
                        }
                    }
                }
            }
        },
        "/v1/posts/{postID}": {
            "get": {
                "description": "Get post by ID",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Posts"
                ],
                "summary": "Get Post",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Post ID",
                        "name": "postID",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/getPost.Response"
                        }
                    }
                }
            },
            "put": {
                "description": "Update post information by ID",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Posts"
                ],
                "summary": "Update Post",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Post ID",
                        "name": "postID",
                        "in": "path",
                        "required": true
                    },
                    {
                        "description": "Update post request",
                        "name": "request",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/updatePost.Request"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/updatePost.Response"
                        }
                    }
                }
            },
            "delete": {
                "description": "Delete post by ID",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Posts"
                ],
                "summary": "Delete Post",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Post ID",
                        "name": "postID",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "204": {
                        "description": "No Content"
                    }
                }
            }
        },
        "/v2/logout": {
            "get": {
                "description": "Logs out the user by clearing the session stored in the cookie \"session_id\".",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "session_auth"
                ],
                "summary": "Logout from session-based authentication",
                "responses": {
                    "200": {
                        "description": "Successfully logged out",
                        "schema": {
                            "type": "string"
                        }
                    }
                }
            }
        },
        "/v2/posts": {
            "get": {
                "description": "Retrieve all posts in the system with session-based authentication (requires \"session_id\" cookie).",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "posts"
                ],
                "summary": "Get all posts",
                "responses": {
                    "200": {
                        "description": "List of posts",
                        "schema": {
                            "type": "array",
                            "items": {
                                "$ref": "#/definitions/getAllPosts.Response"
                            }
                        }
                    }
                }
            },
            "post": {
                "description": "Create a post in the system with session-based authentication (requires \"session_id\" cookie).",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "posts"
                ],
                "summary": "Create a new post",
                "parameters": [
                    {
                        "description": "Post details",
                        "name": "post",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/createPost.Request"
                        }
                    }
                ],
                "responses": {
                    "201": {
                        "description": "Post created successfully",
                        "schema": {
                            "$ref": "#/definitions/createPost.Response"
                        }
                    }
                }
            }
        },
        "/v2/posts/{postID}": {
            "get": {
                "description": "Retrieve a specific post by its ID with session-based authentication (requires \"session_id\" cookie).",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "posts"
                ],
                "summary": "Get post by ID",
                "parameters": [
                    {
                        "type": "string",
                        "description": "ID of the post",
                        "name": "postID",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "Post details",
                        "schema": {
                            "$ref": "#/definitions/getPost.Response"
                        }
                    }
                }
            },
            "put": {
                "description": "Update a specific post by its ID with session-based authentication (requires \"session_id\" cookie).",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "posts"
                ],
                "summary": "Update a post by ID",
                "parameters": [
                    {
                        "type": "string",
                        "description": "ID of the post",
                        "name": "postID",
                        "in": "path",
                        "required": true
                    },
                    {
                        "description": "Updated post details",
                        "name": "post",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/updatePost.Request"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "Post updated successfully",
                        "schema": {
                            "$ref": "#/definitions/updatePost.Response"
                        }
                    }
                }
            },
            "delete": {
                "description": "Delete a specific post by its ID with session-based authentication (requires \"session_id\" cookie).",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "posts"
                ],
                "summary": "Delete a post by ID",
                "parameters": [
                    {
                        "type": "string",
                        "description": "ID of the post",
                        "name": "postID",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "Post deleted successfully",
                        "schema": {
                            "type": "string"
                        }
                    }
                }
            }
        }
    },
    "definitions": {
        "createPost.Request": {
            "type": "object",
            "required": [
                "title"
            ],
            "properties": {
                "content": {
                    "type": "string"
                },
                "tags": {
                    "type": "array",
                    "items": {
                        "type": "string"
                    }
                },
                "title": {
                    "type": "string"
                }
            }
        },
        "createPost.Response": {
            "type": "object",
            "properties": {
                "error": {
                    "type": "string"
                },
                "post": {
                    "$ref": "#/definitions/database.PostDTO"
                },
                "status": {
                    "type": "string"
                }
            }
        },
        "createUser.Request": {
            "type": "object",
            "required": [
                "password",
                "username"
            ],
            "properties": {
                "description": {
                    "type": "string"
                },
                "password": {
                    "type": "string"
                },
                "username": {
                    "type": "string"
                }
            }
        },
        "createUser.Response": {
            "type": "object",
            "properties": {
                "error": {
                    "type": "string"
                },
                "status": {
                    "type": "string"
                },
                "user": {
                    "$ref": "#/definitions/database.UserDTO"
                }
            }
        },
        "database.PostDTO": {
            "type": "object",
            "properties": {
                "content": {
                    "type": "string"
                },
                "created_at": {
                    "$ref": "#/definitions/pgtype.Timestamp"
                },
                "id": {
                    "type": "integer"
                },
                "tags": {
                    "type": "array",
                    "items": {
                        "type": "string"
                    }
                },
                "title": {
                    "type": "string"
                },
                "user_id": {
                    "type": "integer"
                }
            }
        },
        "database.UserDTO": {
            "type": "object",
            "properties": {
                "dateJoined": {
                    "$ref": "#/definitions/pgtype.Date"
                },
                "description": {
                    "type": "string"
                },
                "id": {
                    "type": "integer"
                },
                "password": {
                    "type": "string"
                },
                "username": {
                    "type": "string"
                }
            }
        },
        "getAllPosts.Response": {
            "type": "object",
            "properties": {
                "error": {
                    "type": "string"
                },
                "posts": {
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/database.PostDTO"
                    }
                },
                "status": {
                    "type": "string"
                }
            }
        },
        "getAllUsers.Response": {
            "type": "object",
            "properties": {
                "error": {
                    "type": "string"
                },
                "status": {
                    "type": "string"
                },
                "usernames": {
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/database.UserDTO"
                    }
                }
            }
        },
        "getPost.Response": {
            "type": "object",
            "properties": {
                "error": {
                    "type": "string"
                },
                "post": {
                    "$ref": "#/definitions/database.PostDTO"
                },
                "status": {
                    "type": "string"
                }
            }
        },
        "getUser.Response": {
            "type": "object",
            "properties": {
                "error": {
                    "type": "string"
                },
                "status": {
                    "type": "string"
                },
                "user": {
                    "$ref": "#/definitions/database.UserDTO"
                }
            }
        },
        "jwtLogin.Request": {
            "type": "object",
            "required": [
                "password",
                "username"
            ],
            "properties": {
                "password": {
                    "type": "string"
                },
                "username": {
                    "type": "string"
                }
            }
        },
        "jwtLogin.Response": {
            "type": "object",
            "properties": {
                "access_token": {
                    "type": "string"
                },
                "error": {
                    "type": "string"
                },
                "refresh_token": {
                    "type": "string"
                },
                "status": {
                    "type": "string"
                }
            }
        },
        "pgtype.Date": {
            "type": "object",
            "properties": {
                "infinityModifier": {
                    "$ref": "#/definitions/pgtype.InfinityModifier"
                },
                "time": {
                    "type": "string"
                },
                "valid": {
                    "type": "boolean"
                }
            }
        },
        "pgtype.InfinityModifier": {
            "type": "integer",
            "enum": [
                1,
                0,
                -1
            ],
            "x-enum-varnames": [
                "Infinity",
                "Finite",
                "NegativeInfinity"
            ]
        },
        "pgtype.Timestamp": {
            "type": "object",
            "properties": {
                "infinityModifier": {
                    "$ref": "#/definitions/pgtype.InfinityModifier"
                },
                "time": {
                    "description": "Time zone will be ignored when encoding to PostgreSQL.",
                    "type": "string"
                },
                "valid": {
                    "type": "boolean"
                }
            }
        },
        "refresh.Request": {
            "type": "object",
            "required": [
                "refresh_token"
            ],
            "properties": {
                "refresh_token": {
                    "type": "string"
                }
            }
        },
        "refresh.Response": {
            "type": "object",
            "properties": {
                "access_token": {
                    "type": "string"
                },
                "error": {
                    "type": "string"
                },
                "refresh_token": {
                    "type": "string"
                },
                "status": {
                    "type": "string"
                }
            }
        },
        "sessionLogin.Request": {
            "type": "object",
            "required": [
                "password",
                "username"
            ],
            "properties": {
                "password": {
                    "type": "string"
                },
                "username": {
                    "type": "string"
                }
            }
        },
        "sessionLogin.Response": {
            "type": "object",
            "properties": {
                "error": {
                    "type": "string"
                },
                "session_id": {
                    "type": "string"
                },
                "status": {
                    "type": "string"
                }
            }
        },
        "updatePost.Request": {
            "type": "object",
            "properties": {
                "content": {
                    "type": "string"
                },
                "tags": {
                    "type": "array",
                    "items": {
                        "type": "string"
                    }
                },
                "title": {
                    "type": "string"
                }
            }
        },
        "updatePost.Response": {
            "type": "object",
            "properties": {
                "error": {
                    "type": "string"
                },
                "post": {
                    "$ref": "#/definitions/database.PostDTO"
                },
                "status": {
                    "type": "string"
                }
            }
        },
        "updateUser.Request": {
            "type": "object",
            "properties": {
                "description": {
                    "type": "string"
                },
                "password": {
                    "type": "string"
                },
                "username": {
                    "type": "string"
                }
            }
        },
        "updateUser.Response": {
            "type": "object",
            "properties": {
                "error": {
                    "type": "string"
                },
                "status": {
                    "type": "string"
                },
                "user": {
                    "$ref": "#/definitions/database.UserDTO"
                }
            }
        }
    }
}`

// SwaggerInfo holds exported Swagger Info so clients can modify it
var SwaggerInfo = &swag.Spec{
	Version:          "1.0",
	Host:             "localhost:8000",
	BasePath:         "/",
	Schemes:          []string{},
	Title:            "API DOCUMENTATION",
	Description:      "This is a sample server.",
	InfoInstanceName: "swagger",
	SwaggerTemplate:  docTemplate,
	LeftDelim:        "{{",
	RightDelim:       "}}",
}

func init() {
	swag.Register(SwaggerInfo.InstanceName(), SwaggerInfo)
}

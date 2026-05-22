CREATE TYPE user_role AS ENUM ('user', 'designer', 'admin');

CREATE TABLE users (
    id UUID DEFAULT gen_random_uuid() PRIMARY KEY,
    nickname VARCHAR(20) NOT NULL UNIQUE,
    role user_role NOT NULL,
    password_hash TEXT NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT now()

);

CREATE TABLE designer_profile (
    id UUID DEFAULT gen_random_uuid() PRIMARY KEY,
    avatar_key TEXT,
    nickname VARCHAR(60),
    user_id UUID NOT NULL UNIQUE,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);

CREATE TABLE refresh_token (
    id UUID DEFAULT gen_random_uuid() PRIMARY KEY,
    token_hash TEXT UNIQUE NOT NULL
)
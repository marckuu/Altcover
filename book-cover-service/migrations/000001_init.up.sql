CREATE TABLE book (
    id UUID DEFAULT gen_random_uuid() PRIMARY KEY,
    title VARCHAR(300) NOT NULL,
    description VARCHAR(300)
);

CREATE TABLE cover (
    id UUID DEFAULT gen_random_uuid() PRIMARY KEY,
    title VARCHAR(100) NOT NULL,
    description VARCHAR(200),
    images_keys TEXT,
    status SMALLINT NOT NULL CHECK (status in (0, 1, 2)),
    book_id UUID,
    FOREIGN KEY (book_id) REFERENCES book(id),
    user_id UUID NOT NULL,
    designer_nickname UUID NOT NULL,
    designer_avatar_key TEXT NOT NULL
);

CREATE TABLE comment (
    id UUID DEFAULT gen_random_uuid() PRIMARY KEY,
    text TEXT,
    cover_id UUID,
    FOREIGN KEY (cover_id) REFERENCES cover(id),
    user_id UUID NOT NULL
);


CREATE TABLE favorites (
    user_id UUID NOT NULL,
    cover_id UUID NOT NULL,
    PRIMARY KEY (user_id, cover_id),
    FOREIGN KEY (cover_id) REFERENCES cover(id) ON DELETE CASCADE
);

CREATE TABLE cover_like (
    user_id UUID NOT NULL,
    cover_id UUID NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT now()
);

CREATE TABLE designer_profile_snapshot (
    id UUID PRIMARY KEY,
    avatar_key TEXT,
    nickname VARCHAR(60)
);
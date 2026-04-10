CREATE TYPE status AS ENUM ('uploaded', 'onModeration', 'hidden');

CREATE TABLE cover (
    id UUID DEFAULT gen_random_uuid() PRIMARY KEY,
    title VARCHAR(100) NOT NULL,
    description VARCHAR(200),
    likes INT NOT NULL,
    images_keys TEXT,
    status status NOT NULL,
    book_id UUID,
    FOREIGN KEY (book_id) REFERENCES book(id),
    designer_id UUID NOT NULL,
    designer_nickname UUID NOT NULL,
    designer_avatar_key TEXT NOT NULL,
);

CREATE TABLE book (
    id UUID DEFAULT gen_random_uuid() PRIMARY KEY,
    title VARCHAR(300) NOT NULL,
    description VARCHAR(300)
);

CREATE TABLE comment (
    id UUID DEFAULT gen_random_uuid() PRIMARY KEY,
    text TEXT,
    cover_id UUID,
    FOREIGN KEY (cover_id) REFERENCES cover(id),
    user_id UUID NOT NULL
);



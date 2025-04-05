CREATE TABLE users (
    id UUID DEFAULT gen_random_uuid() PRIMARY KEY,
    full_name TEXT NOT NULL,
    username TEXT NOT NULL,
    email TEXT NOT NULL,
    password TEXT NOT NULL
);

CREATE TABLE messages (
    id UUID DEFAULT gen_random_uuid() PRIMARY KEY,
    room_id TEXT NOT NULL,
    message TEXT NOT NULL,
    sender_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    created_at TIMESTAMP DEFAULT NOW(),
    read_at TIMESTAMP NULL
);

CREATE TABLE rooms (
    id UUID DEFAULT gen_random_uuid() PRIMARY KEY,
    property_id TEXT NOT NULL,
    property_owner_id TEXT NOT NULL,
    customer_id TEXT NOT NULL
);

CREATE TABLE listings (
    id UUID DEFAULT gen_random_uuid() PRIMARY KEY,
    title TEXT NOT NULL,
    description TEXT NOT NULL,
    type TEXT NOT NULL,
    price INTEGER NOT NULL CHECK (price > 0) DEFAULT 0,
    location TEXT NOT NULL,
    bathrooms INTEGER NOT NULL CHECK (bathrooms >= 0) DEFAULT 0,
    bedrooms INTEGER NOT NULL CHECK (bedrooms >= 0) DEFAULT 0,
    created_at TIMESTAMP DEFAULT NOW(),
    image_urls TEXT[] NOT NULL,
    is_air_conditioned BOOLEAN NOT NULL DEFAULT FALSE,
    is_balcony_available BOOLEAN NOT NULL DEFAULT FALSE,
    is_dryer_available BOOLEAN NOT NULL DEFAULT FALSE,
    is_heated BOOLEAN NOT NULL DEFAULT FALSE,
    is_parking_available BOOLEAN NOT NULL DEFAULT FALSE,
    is_pool_available BOOLEAN NOT NULL DEFAULT FALSE,
    is_washer_available BOOLEAN NOT NULL DEFAULT FALSE,
    is_wifi_available BOOLEAN NOT NULL DEFAULT FALSE,
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE
);